package job

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/deepinsight/mlcloud/src/pkg/crm"
	jobs "github.com/deepinsight/mlcloud/src/pkg/job"
	"github.com/deepinsight/mlcloud/src/pkg/models"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func runSubmit(cmd *cobra.Command, args []string) {

	flag.CommandLine.Parse([]string{})

	glog.V(4).Info("run mlcloud job submit")

	job, err := parseArguments(cmd)
	if err != nil {
		glog.Error(err)
	}

	glog.V(4).Info(job)

	err = jobs.SubmitJob(job)
	if err != nil {
		//glog.Fatal(err)
	}
}

func parseArguments(cmd *cobra.Command) (*models.Job, error) {
	hasMaster, err := cmd.Flags().GetBool("master")
	if err != nil {
		return nil, err
	}

	nbWorkers, err := cmd.Flags().GetInt("num-worker")
	if err != nil {
		return nil, err
	}

	nbPS, err := cmd.Flags().GetInt("num-ps")
	if err != nil {
		return nil, err
	}

	tensorboard, err := cmd.Flags().GetBool("tensorboard")
	if err != nil {
		return nil, err
	}

	numGPU, err := cmd.Flags().GetInt("num-gpu")
	if err != nil {
		return nil, err
	}

	jobName, err := cmd.Flags().GetString("name")
	if err != nil {
		return nil, err
	}
	if jobName == "" {
		fmt.Println("must specify the job's name")
		os.Exit(1)
	}

	image, err := cmd.Flags().GetString("image")
	if err != nil {
		return nil, err
	}

	jobType, err := cmd.Flags().GetString("type")
	if err != nil {
		return nil, err
	}

	if image == "" {
		switch jobType {
		case models.JobTypeTensorflow:
			if numGPU > 0 {
				image = crm.DefaultTfGPUImage
			} else {
				image = crm.DefaultTfImage
			}
		case models.JobTypeMxnet:
			if numGPU > 0 {
				image = crm.DefaultMxGPUImage
			} else {
				image = crm.DefaultMxImage
			}
		}
	}

	command, err := cmd.Flags().GetString("command")
	if err != nil {
		return nil, err
	}

	args, err := cmd.Flags().GetString("args")
	if err != nil {
		return nil, err
	}
	argslice := strings.Fields(args)

	logDir, err := cmd.Flags().GetString("log-dir")
	if err != nil {
		return nil, err
	}

	mxnetMode, err := cmd.Flags().GetString("mxnet-mode")
	if err != nil {
		return nil, err
	}

	job := &models.Job{
		Name: jobName,
		Type: jobType,
	}

	switch job.Type {
	case models.JobTypeTensorflow:
		job.Tensorflow = &models.TensorflowJob{
			HasMaster:       hasMaster,
			NumWorkers:      nbWorkers,
			Image:           image,
			Command:         command,
			ArgsSlice:       &argslice,
			Tensorboard:     tensorboard,
			TensorboardHost: jobName + "." + jobs.Username + ".test-cloud.bigdata.test.cn",
			LogDir:          logDir,
			NumGPU:          numGPU,
		}
	case models.JobTypeMxnet:
		if mxnetMode != "local" && mxnetMode != "dist" {
			return nil, errors.New("mxnet mode must be local or dist")
		}
		if mxnetMode == "local" {
			nbWorkers = 1
		}
		job.Mxnet = &models.MxnetJob{
			Mode:       mxnetMode,
			NumWorkers: nbWorkers,
			NumPs:      nbPS,
			Image:      image,
			Command:    command,
			ArgsSlice:  &argslice,
			LogDir:     logDir,
			NumGPU:     numGPU,
		}
	default:
		//todo: raise error
		return nil, errors.New("unsupport machine learning framework; current support list: tensorflow, mxnet.")
	}

	return job, nil
}

func NewSubmitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "submit job",
		//Long: `download a file or directory from remote to local`,
		//Args: cobra.MinimumNArgs(0),
		Run: runSubmit,
	}

	cmd.Flags().BoolP("master", "", false, "distributed training contains master")
	cmd.Flags().Int("num-worker", 0, "number of workers")
	cmd.Flags().Int("num-ps", 0, "number of parameter servers")
	cmd.Flags().BoolP("tensorboard", "", false, "use tensorboard")
	cmd.Flags().Int("num-gpu", 0, "number of gpus")
	cmd.Flags().String("name", "", "job name")
	cmd.Flags().String("image", "", "docker image")
	cmd.Flags().String("type", "", "machine learning framework")
	cmd.Flags().String("command", "", "training stript to execute")
	cmd.Flags().String("args", "", "arguments to training commnad")
	cmd.Flags().String("log-dir", "", "logging directory")

	cmd.Flags().String("mxnet-mode", "", "runing mode for mxnet: local, dist")

	// add children commands here

	return cmd
}
