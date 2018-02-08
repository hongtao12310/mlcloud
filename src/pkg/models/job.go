package models

import (
	"errors"
	"strings"
	"time"

	"fmt"
	"strconv"

	"github.com/astaxie/beego/orm"
	mxspec "github.com/deepinsight/mxnet-operator/pkg/spec"
	tfspec "github.com/deepinsight/tf-operator/pkg/spec"
)

const (
	JobTypeTensorflow = "tensorflow"
	JobTypeMxnet      = "mxnet"
)

func init() {
	orm.RegisterModel(new(Job))
}

type Job struct {
	Id           int            `orm:"pk;auto" json:"ID"`
	Name         string         `json:"Name"`
	Type         string         `json:"Type"`
	User         *User          `orm:"rel(fk)" json:"-"`
	Tensorflow   *TensorflowJob `orm:"-" json:"Tensorflow,omitempty"` // orm: ignore
	Mxnet        *MxnetJob      `orm:"-" json:"Mxnet,omitempty"`
	CreationTime time.Time      `json:"creationTime"`
	UpdateTime   time.Time      `json:"updateTime"`
}

type JobStatus struct {
	Jb       *Job                `json:"JobDescription"`
	MxStatus *mxspec.MxJobStatus `json:"mxStatus,omitempty"`
	TfStatus *tfspec.TfJobStatus `json:"tfStatus,omitempty"`
}

func (job *Job) ArgSlice2Str() {
	switch job.Type {
	case JobTypeTensorflow:
		if job.Tensorflow.ArgsSlice != nil {
			job.Tensorflow.Arguments = strings.Join(*job.Tensorflow.ArgsSlice, " ")
		}
	case JobTypeMxnet:
		if job.Mxnet.ArgsSlice != nil {
			job.Mxnet.Arguments = strings.Join(*job.Mxnet.ArgsSlice, " ")
		}
	}
}
func (job *Job) ArgStr2Slice() {
	switch job.Type {
	case JobTypeTensorflow:
		if len(job.Tensorflow.Arguments) > 0 {
			argslice := strings.Split(job.Tensorflow.Arguments, " ")
			job.Tensorflow.ArgsSlice = &argslice
		}
	case JobTypeMxnet:
		if len(job.Mxnet.Arguments) > 0 {
			argslice := strings.Split(job.Mxnet.Arguments, " ")
			job.Mxnet.ArgsSlice = &argslice
		}
	}
}

func (job *Job) Validate() error {

	if job == nil {
		return fmt.Errorf("ValidateTaskTypeAndIndex: job is nil")
	}

	switch job.Type {
	case JobTypeTensorflow:
	case JobTypeMxnet:
	default:
		return errors.New("invalid job type")
	}

	return nil
}

func (job *Job) ValidateTaskAndIndex(task, taskIndex string) error {

	if job == nil {
		return fmt.Errorf("ValidateTaskTypeAndIndex: job is nil")
	}

	index, err := strconv.Atoi(taskIndex)
	if err != nil {
		return fmt.Errorf("Task index (%s) should be int", taskIndex)
	}

	switch job.Type {
	case JobTypeTensorflow:
		switch tfspec.TfReplicaType(task) {
		case tfspec.MASTER:
		case tfspec.PS:
			if index >= job.Tensorflow.NumPs || index < 0 {
				return fmt.Errorf("invalid task (%s) and index (%s)", task, taskIndex)
			}
		case tfspec.WORKER:
			if index >= job.Tensorflow.NumWorkers || index < 0 {
				return fmt.Errorf("invalid task (%s) and index (%s)", task, taskIndex)
			}
		default:
			return fmt.Errorf("Invalid task: %s", task)
		}

	case JobTypeMxnet:
		switch mxspec.MxReplicaType(task) {
		case mxspec.SCHEDULER:
			if index > 1 || index < 0 {
				return fmt.Errorf("invalid task index (%s)", taskIndex)
			}
		case mxspec.SERVER:
			if index > job.Mxnet.NumPs || index < 0 {
				return fmt.Errorf("invalid task index (%s)", taskIndex)
			}
		case mxspec.WORKER:
			if index > job.Mxnet.NumWorkers || index < 0 {
				return fmt.Errorf("invalid task index (%s)", taskIndex)
			}
		default:
			return fmt.Errorf("Invalid task: %s", task)
		}
	default:
		return errors.New("invalid job type")
	}

	return nil
}
