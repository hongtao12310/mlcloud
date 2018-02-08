package models

import "github.com/astaxie/beego/orm"

type TensorflowJob struct {
	Id              int       `json:"ID" orm:"pk;auto"`
	NumPs           int       `json:"NumPs"`
	NumWorkers      int       `json:"NumWorkers"`
	Image           string    `json:"Image"`
	DataDir         string    `json:"dataDir"`
	LogDir          string    `json:"logDir"`
	Command         string    `json:"Command,omitempty"` // "<dir>/learn.py"
	Arguments       string    `json:"-"`
	ArgsSlice       *[]string `json:"Arguments" orm:"-"`
	NumGPU          int       `json:"NumGpu" orm:"column(num_gpu)"`
	Tensorboard     bool      `json:"Tensorboard"`
	TensorboardHost string    `json:"TensorboardHost"`
	HasMaster       bool      `json:"hasMaster"`
}

func init() {
	orm.RegisterModel(new(TensorflowJob))
}
