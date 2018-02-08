package job

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/deepinsight/mlcloud/src/cmd/mlcloud/utils"
	"github.com/deepinsight/mlcloud/src/pkg/models"
	httputil "github.com/deepinsight/mlcloud/src/utils/http"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
)

var endpoint, Username string

func init() {
	config, err := utils.ParseDefaultConfig()
	if err != nil {
		fmt.Printf("parse configuration error: %v", err)
		os.Exit(1)
	}

	endpoint = config.ActiveConfig.Endpoint
	Username = config.ActiveConfig.Username
}

func ListJobs() error {
	targetURL := endpoint + "/api/v1/users/" + Username + "/jobs"

	req, err := utils.MakeRequestToken(targetURL, "GET", nil, "", nil)
	if err != nil {
		return fmt.Errorf("ListJobs: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jobs := make([]models.Job, 0)
	jrsp := httputil.Response{Results: &jobs}

	if err = json.Unmarshal(body, &jrsp); err != nil {
		return err
	}

	if jrsp.Code != 200 {
		fmt.Println("GetJob error! ", jrsp.Err)
		return nil
	}

	if len(jobs) == 0 {
		return nil
	}

	y, err := yaml.Marshal(jobs)
	if err != nil {
		return err
	}

	fmt.Println(string(y))
	return nil
}

func GetJob(jobname string) error {
	targetURL := endpoint + "/api/v1/users/" + Username + "/jobs/" + jobname

	req, err := utils.MakeRequestToken(targetURL, "GET", nil, "", nil)
	if err != nil {
		return fmt.Errorf("GetJob: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	js := &models.JobStatus{}
	jrsp := httputil.Response{Results: js}

	err = json.Unmarshal(body, &jrsp)
	if err != nil {
		return fmt.Errorf("GetJob: %v", err)
	}

	if jrsp.Code != 200 {
		fmt.Println("GetJob error: ", jrsp.Err)
		return nil
	}

	if js == nil {
		fmt.Println("js is nil")
		return nil
	}

	y, err := yaml.Marshal(*js)
	if err != nil {
		return fmt.Errorf("GetJob: %v", err)
	}
	fmt.Println(string(y))
	return nil
}

func DeleteJob(jobname string) error {
	targetURL := endpoint + "/api/v1/users/" + Username + "/jobs/" + jobname

	req, err := utils.MakeRequestToken(targetURL, "DELETE", nil, "", nil)
	if err != nil {
		return fmt.Errorf("DeleteJob: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jrsp := httputil.Response{}
	if err = json.Unmarshal(body, &jrsp); err != nil {
		return err
	}

	if jrsp.Code != 200 {
		fmt.Println("Delete Job error: ", jrsp.Err)
		return nil
	}

	fmt.Printf("Delete job %s success!\n", jobname)
	return nil
}

type LogOpt struct {
	Follow       bool
	TimeStamps   bool
	LimitBytes   int64
	TailLine     int64
	SinceSeconds int64
}

func GetJobLogs(jobname, task string, index int, opt *LogOpt) error {

	queryString := "?follow=" + strconv.FormatBool(opt.Follow)
	queryString += "&timestamp=" + strconv.FormatBool(opt.TimeStamps)
	queryString += "&limit-bytes=" + strconv.FormatInt(opt.LimitBytes, 10)
	queryString += "&tail=" + strconv.FormatInt(opt.TailLine, 10)
	queryString += "&since=" + strconv.FormatInt(opt.SinceSeconds, 10)

	targetURL := endpoint + "/api/v1/users/" + Username + "/jobs/" + jobname + "/" +
		task + "/" + strconv.Itoa(index) + "/logs" + queryString

	req, err := utils.MakeRequestToken(targetURL, "GET", nil, "", nil)
	if err != nil {
		glog.V(4).Infof("GetJobLogs: %v", err)
		return fmt.Errorf("GetJobLogs: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

	return nil
}

func SubmitJob(job *models.Job) error {

	if job == nil {
		glog.V(4).Info("SubmitJob: job is nil")
		return fmt.Errorf("SubmitJob: job is nil")
	}

	url := endpoint + "/api/v1/users/" + Username + "/jobs"

	j, err := json.Marshal(job)
	if err != nil {
		glog.V(4).Infof("SubmitJob: unable to marshal job: %v", err)
		return err
	}

	body, err := utils.PostCall(url, j)
	if err != nil {
		return fmt.Errorf("SubmitJob: PostCall error: %v", err)
	}

	jrsp := httputil.Response{}
	if err = json.Unmarshal(body, &jrsp); err != nil {
		return err
	}

	if jrsp.Code != 200 {
		fmt.Println("Submit Job error: ", jrsp.Err)
		return nil
	}

	y, err := yaml.Marshal(job)
	if err != nil {
		return err
	}
	fmt.Printf("Submit Job %s success!\n", job.Name)
	fmt.Println(string(y))
	return nil
}
