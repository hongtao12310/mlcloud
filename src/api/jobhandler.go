package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/deepinsight/mlcloud/src/pkg/dao"
	"github.com/deepinsight/mlcloud/src/pkg/models"

	"fmt"
	"strings"

	"github.com/deepinsight/mlcloud/src/pkg/crm"
	httputil "github.com/deepinsight/mlcloud/src/utils/http"
	"github.com/deepinsight/tf-operator/pkg/util"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

type JobHandler struct {
	*BaseHandler
}

// NewJobHandler created JobHandler instance.
func NewJobHandler(baseHandler *BaseHandler) *JobHandler {
	return &JobHandler{
		BaseHandler: baseHandler,
	}
}

func (self *JobHandler) Register(router *mux.Router) {
	bs := router.PathPrefix("/api/v1").Subrouter()

	bs.HandleFunc("/users/{username}/jobs", self.ListUserJobs).
		Methods("GET")

	bs.HandleFunc("/users/{username}/jobs", self.SubmitJob).
		Methods("POST")

	bs.HandleFunc("/users/{username}/jobs/{jobname}", self.GetUserJob).
		Methods("GET")

	bs.HandleFunc("/users/{username}/jobs/{jobname}", self.DeleteUserJob).
		Methods("DELETE")

	bs.HandleFunc("/users/{username}/jobs/{jobname}/{task}/{index}/logs", self.GetLogs).
		Queries("follow", "{follow}").
		Queries("timestamp", "{timestamp}").
		Queries("limit-bytes", "{limit}").
		Queries("tail", "{tail}").
		Queries("since", "{since}").
		Methods("GET")
}

func (self *JobHandler) ListUserJobs(writer http.ResponseWriter, request *http.Request) {
	resp := httputil.Response{}
	user, err := self.getUserInfo(writer, request)
	if err != nil {
		resp.Code = http.StatusUnauthorized
		resp.Err = "Authentication failed"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	username := mux.Vars(request)["username"]
	if !user.Sysadmin && user.Username != username { // user (not admin) wants to access other's namespace
		glog.V(4).Infof("%s wants to access %s's space", user.Username, username)
		resp.Code = http.StatusUnauthorized
		resp.Err = "Unauthorized"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	// Get user by username
	user, err = dao.GetUser(&models.User{Username: username})
	if err != nil {
		resp.Code = http.StatusBadRequest
		glog.V(4).Infof("%v", err)
		resp.Err = fmt.Sprintf("User (%s) doesn't exist", username)
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	jobs, err := dao.ListJobsByUser(user)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Err = "Cannot get jobs in database"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	resp.Code = http.StatusOK
	resp.Results = jobs
	httputil.WriteJSONResponse(writer, resp)
}

func (self *JobHandler) SubmitJob(writer http.ResponseWriter, request *http.Request) {
	glog.V(4).Infof("Post %s: %v", request.RequestURI, request.Body)

	resp := httputil.Response{}
	user, err := self.getUserInfo(writer, request)
	if err != nil {
		glog.Errorf("submit tfjob: %v", err)
		resp.Code = http.StatusUnauthorized
		resp.Err = "Authentication failed"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	username := mux.Vars(request)["username"]
	if !user.Sysadmin && user.Username != username { // user (not admin) wants to access other's namespace
		glog.Errorf("%s wants to access %s's space", user.Username, username)
		resp.Code = http.StatusUnauthorized
		resp.Err = "Unauthorized"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	// Get user by username
	user, err = dao.GetUser(&models.User{Username: username})
	if err != nil || user == nil {
		glog.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Err = fmt.Sprintf("User (%s) doesn't exist", username)
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	job := &models.Job{}
	if err := httputil.ReadEntity(request, job); err != nil {
		glog.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Err = err.Error()
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	job.ArgSlice2Str()

	if err = job.Validate(); err != nil {
		glog.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Err = "invalid job"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	job.User = user

	if _, err := dao.CreateJob(job); err != nil {
		if e, ok := err.(dao.JobAlreadyExistError); ok {
			// Job already exists
			resp.Code = http.StatusNotAcceptable
			resp.Err = e.Error()
			httputil.WriteJSONResponse(writer, resp)
			return
		}

		// Failed to create job
		glog.V(4).Infof("failed to create job: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Err = err.Error()
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	switch job.Type {
	case models.JobTypeTensorflow:
		tfjob, err := crm.GenerateTFJob(job)
		glog.V(4).Info(util.Pformat(tfjob))
		if err != nil {
			glog.Error(err)
			resp.Code = http.StatusInternalServerError
			resp.Err = "cannot generate tfjob"
			httputil.WriteJSONResponse(writer, resp)
			return
		}
		self.clientset.TFJob().Create(job.User.Username, tfjob)
	case models.JobTypeMxnet:
		mxjob, err := crm.GenerateMXJob(job)
		glog.V(4).Info(util.Pformat(mxjob))
		if err != nil {
			glog.Error(err)
			resp.Code = http.StatusInternalServerError
			resp.Err = "cannot generate mxjob"
			httputil.WriteJSONResponse(writer, resp)
			return
		}
		self.clientset.MXJob().Create(job.User.Username, mxjob)
		fmt.Printf("sucess to create mx job \n", job.Name)
	}

	resp.Results = job
	resp.Code = http.StatusOK
	httputil.WriteJSONResponse(writer, resp)
}

func (self *JobHandler) GetUserJob(writer http.ResponseWriter, request *http.Request) {
	resp := httputil.Response{}
	user, err := self.getUserInfo(writer, request)
	if err != nil {
		glog.Error(err)
		resp.Code = http.StatusUnauthorized
		resp.Err = "Authentication failed"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	username := mux.Vars(request)["username"]
	if !user.Sysadmin && user.Username != username { // user (not admin) wants to access other's namespace
		glog.Errorf("%s wants to access %s's space", user.Username, username)
		resp.Code = http.StatusUnauthorized
		resp.Err = "Unauthorized"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	// Get user by username
	user, err = dao.GetUser(&models.User{Username: username})
	if err != nil {
		glog.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Err = fmt.Sprintf("User (%s) doesn't exist", username)
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	jobname := mux.Vars(request)["jobname"]

	job, err := dao.GetJob(&models.Job{Name: jobname, User: user})
	if err != nil {
		glog.Errorf("job not found: %v", err)
		resp.Code = http.StatusBadRequest
		resp.Err = fmt.Sprintf("Job (%s) doesn't exist", jobname)
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	job.ArgStr2Slice()

	var js *models.JobStatus = &models.JobStatus{}
	js.Jb = job
	switch job.Type {
	case models.JobTypeTensorflow:
		tfjob, err := self.clientset.TFJob().Get(user.Username, jobname)
		js.TfStatus = &tfjob.Status
		if err == nil {
			resp.Results = js
		} else {
			glog.Errorf("get tfjob crd err: %v", err)
			resp.Code = http.StatusBadRequest
			resp.Err = fmt.Sprintf("get tfjob crd err: %v", err)
			httputil.WriteJSONResponse(writer, resp)
			return
		}
	case models.JobTypeMxnet:
		mxjob, err := self.clientset.MXJob().Get(user.Username, jobname)
		js.MxStatus = &mxjob.Status
		if err == nil {
			resp.Results = js
		} else {
			glog.Errorf("get mxjob crd err: %v", err)
			resp.Code = http.StatusBadRequest
			resp.Err = fmt.Sprintf("get mxjob crd err: %v", err)
			httputil.WriteJSONResponse(writer, resp)
			return
		}
	}

	resp.Code = http.StatusOK
	httputil.WriteJSONResponse(writer, resp)
}

func (self *JobHandler) DeleteUserJob(writer http.ResponseWriter, request *http.Request) {
	resp := httputil.Response{}
	user, err := self.getUserInfo(writer, request)
	if err != nil {
		glog.Error(err)
		resp.Code = http.StatusUnauthorized
		resp.Err = "Authentication failed"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	username := mux.Vars(request)["username"]
	if !user.Sysadmin && user.Username != username { // user (not admin) wants to access other's namespace
		glog.Errorf("%s wants to access %s's space", user.Username, username)
		resp.Code = http.StatusUnauthorized
		resp.Err = "Unauthorized"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	// Get user by username
	user, err = dao.GetUser(&models.User{Username: username})
	if err != nil {
		glog.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Err = fmt.Sprintf("User (%s) doesn't exist", username)
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	jobname := mux.Vars(request)["jobname"]

	job, err := dao.GetJob(&models.Job{Name: jobname, User: user})
	if err != nil {
		glog.Errorf("delete job not found: %v", err)
		resp.Code = http.StatusBadRequest
		resp.Err = fmt.Sprintf("Job (%s) doesn't exist", jobname)
		httputil.WriteJSONResponse(writer, resp)
		return
	}
	job.ArgStr2Slice()

	glog.V(4).Infof(util.Pformat(job))

	// Delete TjJob in k8s cluster
	switch job.Type {
	case models.JobTypeTensorflow:
		if err := self.clientset.TFJob().Delete(job.User.Username, job.Name); err != nil {
			glog.Errorf("fail to delete tfjob: %v", err)
			// MARK: Sometimes tfjob doesn't exist in the cluster (Debugging, switch of cluster, etc),
			// so temporally ignore this error
			//resp.Code = http.StatusInternalServerError
			//resp.Err = fmt.Sprintf( "Failed to delete tfjob %s", job.Name)
			//httputil.WriteJSONResponse(writer, resp)
			//return
		}
	case models.JobTypeMxnet:
		if _, err := self.clientset.MXJob().Delete(job.User.Username, job.Name); err != nil {
			glog.Errorf("fail to delete mxjob: %v", err)
			// MARK: Sometimes mxjob doesn't exist in the cluster (Debugging, switch of cluster, etc),
			// so temporally ignore this error
			//resp.Code = http.StatusInternalServerError
			//resp.Err = fmt.Sprintf( "Failed to delete tfjob %s", job.Name)
			//httputil.WriteJSONResponse(writer, resp)
			//return
		}
		fmt.Printf("success to delete mxjob %v\n", job.Name)
	}

	err = dao.DeleteJob(job)
	if err != nil {
		glog.Error(err)
		resp.Code = http.StatusInternalServerError
		resp.Err = "cannot delete job in database"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	resp.Code = http.StatusOK
	resp.Results = job
	httputil.WriteJSONResponse(writer, resp)
}

func (self *JobHandler) GetLogs(writer http.ResponseWriter, request *http.Request) {
	resp := httputil.Response{}
	user, err := self.getUserInfo(writer, request)
	if err != nil {
		glog.Error(err)
		resp.Code = http.StatusUnauthorized
		resp.Err = "Authentication failed"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	username := mux.Vars(request)["username"]
	if !user.Sysadmin && user.Username != username { // user (not admin) wants to access other's namespace
		glog.Errorf("%s wants to access %s's space", user.Username, username)
		resp.Code = http.StatusUnauthorized
		resp.Err = "Unauthorized"
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	// Get user by username
	user, err = dao.GetUser(&models.User{Username: username})
	if err != nil {
		glog.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Err = fmt.Sprintf("User (%s) doesn't exist", username)
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	// Get Job
	jobname := mux.Vars(request)["jobname"]
	job, err := dao.GetJob(&models.Job{Name: jobname, User: user})
	if err != nil {
		glog.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Err = fmt.Sprintf("Job (%s) doesn't exist", jobname)
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	// Validate Task (ps/worker) & Index
	task := strings.ToUpper(mux.Vars(request)["task"])
	taskIndex := mux.Vars(request)["index"]

	if err = job.ValidateTaskAndIndex(task, taskIndex); err != nil {
		glog.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Err = fmt.Sprintf(err.Error())
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	// Get Pod
	// TODO(bug): here labels don't contain jobname, so mlcloud could retrieve pods of other jobs
	var labels []string
	if job.Type == "tensorflow" {
		labels = []string{
			"mlkube.io=",
			"job_type=" + task,
			"task_index=" + taskIndex,
		}
	} else if job.Type == "mxnet" {
		labels = []string{
			"mxnet.mlkube.io=",
			"job_type=" + task,
			"task_index=" + taskIndex,
		}
	}

	options := metav1.ListOptions{
		LabelSelector: strings.Join(labels, ","),
	}
	pods, err := self.clientset.CoreV1().Pods(user.Username).List(options)
	if err != nil {
		glog.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Err = fmt.Sprintf("Cannot list pods")
		httputil.WriteJSONResponse(writer, resp)
		return
	}

	if len(pods.Items) > 0 {
		pod := pods.Items[0]

		follow, _ := strconv.ParseBool(mux.Vars(request)["follow"])
		timestamp, _ := strconv.ParseBool(mux.Vars(request)["timestamp"])
		limit, _ := strconv.ParseInt(mux.Vars(request)["limit"], 10, 64)
		tail, _ := strconv.ParseInt(mux.Vars(request)["tail"], 10, 64)
		since, _ := strconv.ParseInt(mux.Vars(request)["since"], 10, 64)

		opt := &v1.PodLogOptions{
			Follow:     follow,
			Timestamps: timestamp,
		}
		if limit > 0 {
			opt.LimitBytes = &limit
		}
		if tail != -1 {
			opt.TailLines = &tail
		}
		if since != 0 {
			opt.SinceSeconds = &since
		}

		rc, err := self.clientset.CoreV1().Pods(user.Username).GetLogs(pod.Name, opt).Stream()
		if err != nil {
			glog.V(4).Infof(err.Error())
			resp.Code = http.StatusBadRequest
			resp.Err = fmt.Sprintf("%s", err.Error())
			httputil.WriteJSONResponse(writer, resp)
			return
		}
		defer rc.Close()

		flushWriter := httputil.Wrap(writer)
		_, err = io.Copy(flushWriter, rc)
		if err != nil {
			glog.Error(err)
		}
		return
	} else {
		glog.V(4).Infof("Please wait, k8s haven't created the training container yet")
		resp.Code = http.StatusBadRequest
		resp.Err = fmt.Sprintf("Please wait, k8s haven't created the training container yet")
		httputil.WriteJSONResponse(writer, resp)
		return
	}
}
