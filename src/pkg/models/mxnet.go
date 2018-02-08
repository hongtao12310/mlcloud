package models

import "github.com/astaxie/beego/orm"

type MxnetJob struct {
	Id         int       `json:"ID" orm:"pk;auto"`
	Mode       string    `json:"Mode" orm:"column(mode)"`
	NumPs      int       `json:"NumPs"`
	NumWorkers int       `json:"NumWorkers"`
	Image      string    `json:"Image"`
	DataDir    string    `json:"dataDir"`
	LogDir     string    `json:"logDir"`
	Command    string    `json:"Command,omitempty"` // "<dir>/learn.py"
	Arguments  string    `json:"-"`
	ArgsSlice  *[]string `json:"Arguments" orm:"-"`
	NumGPU     int       `json:"NumGpu" orm:"column(num_gpu)"`
}

func init() {
	orm.RegisterModel(new(MxnetJob))
}
