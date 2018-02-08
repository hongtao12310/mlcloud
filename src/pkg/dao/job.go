package dao

import (
	"fmt"
	"time"

	"github.com/deepinsight/mlcloud/src/pkg/models"
	"github.com/deepinsight/mlcloud/src/utils/log"
	"github.com/deepinsight/tf-operator/pkg/util"
	"github.com/golang/glog"
)

// Get job by id or (jobName, userId)
func GetJob(job *models.Job) (*models.Job, error) {

	if job == nil {
		return nil, fmt.Errorf("GetJob: job is nil")
	}

	o := GetOrmer()

	if job.Id != 0 {
		// If job has id, use id to retrieve job
		err := o.Read(job)
		if err != nil {
			return nil, fmt.Errorf("cannot get job: %v", err)
		}

		err = o.Read(job.User)
		if err != nil {
			return nil, fmt.Errorf("cannot get user: %v", err)
		}
	} else {
		// if job has no id, use jobname & userId to retrieve job
		if job.Name != "" {
			err := o.Read(job, "Name", "User")
			if err != nil {
				return nil, fmt.Errorf("cannot get job: %v", err)
			}

			err = o.Read(job.User)
			if err != nil {
				return nil, fmt.Errorf("cannot get user: %v", err)
			}
		} else {
			glog.V(4).Infof("failed to get job: %v", util.Pformat(job))
			return nil, fmt.Errorf("failed to get job: job info is not complete")
		}
	}

	switch job.Type {
	case models.JobTypeTensorflow:
		job.Tensorflow = &models.TensorflowJob{Id: job.Id}
		err := o.Read(job.Tensorflow)
		if err != nil {
			return job, fmt.Errorf("Cannot get tensorflow job")
		}
	case models.JobTypeMxnet:
		job.Mxnet = &models.MxnetJob{Id: job.Id}
		err := o.Read(job.Mxnet)
		if err != nil {
			return job, fmt.Errorf("Cannot get mxnet job")
		}
	}

	return job, nil
}

func CreateJob(job *models.Job) (int64, error) {
	if job == nil {
		return -1, fmt.Errorf("job is nil")
	}

	if j, _ := GetJob(job); j != nil {
		return -1, JobAlreadyExistError{job.Name, job.User.Username}
	}

	o := GetOrmer()
	job.CreationTime = time.Now()
	job.UpdateTime = time.Now()

	id, err := o.Insert(job)
	if err != nil {
		return -1, fmt.Errorf("failed to create job: %v", err)
	}

	switch job.Type {
	case models.JobTypeTensorflow:
		if job.Tensorflow == nil {
			return -1, fmt.Errorf("tensorflow job undefined")
		}

		job.Tensorflow.Id = int(id)

		_, err := o.Insert(job.Tensorflow)
		if err != nil {
			return id, fmt.Errorf("failed to insert tensorflow job into db: %v", err)
		}
	case models.JobTypeMxnet:
		job.Mxnet.Id = int(id)

		_, err := o.Insert(job.Mxnet)
		if err != nil {
			return id, fmt.Errorf("failed to insert mxnet job into db: %v", err)
		}
	}

	return id, nil
}

// Delete job
func DeleteJob(job *models.Job) error {

	if job == nil {
		return fmt.Errorf("Delete Job: job cannot be nil")
	}

	o := GetOrmer()

	switch job.Type {
	case models.JobTypeTensorflow:
		if job.Tensorflow != nil {
			_, err := o.Delete(job.Tensorflow)
			if err != nil {
				return fmt.Errorf("Delete Job: failed to delete tensorflow job from db")
			}
		}

	case models.JobTypeMxnet:
		if job.Mxnet != nil {
			_, err := o.Delete(job.Mxnet)
			if err != nil {
				return fmt.Errorf("Delete Job: failed to delete mxnet job from db")
			}
		}
	}

	_, err := o.Delete(job)

	return err
}

func ListJobsByUser(user *models.User) (*[]models.Job, error) {
	o := GetOrmer()

	var jobs []models.Job

	if user == nil {
		return &jobs, nil
	}

	_, err := o.QueryTable("job").Filter("user_id", user.Id).All(&jobs)
	if err != nil {
		return nil, err
	}

	glog.V(4).Infof(util.Pformat(jobs))

	for index, job := range jobs {
		switch job.Type {
		case models.JobTypeTensorflow:
			jobs[index].Tensorflow = &models.TensorflowJob{Id: job.Id}
			err = o.Read(jobs[index].Tensorflow)
			if err != nil {
				log.Error(err)
			}
		case models.JobTypeMxnet:
			jobs[index].Mxnet = &models.MxnetJob{Id: job.Id}
			err = o.Read(jobs[index].Mxnet)
			jobs[index].ArgStr2Slice()
			if err != nil {
				log.Error(err)
			}

		}
	}

	return &jobs, err
}
