package job

import (
	"flag"
	"math"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"fmt"

	job "github.com/deepinsight/mlcloud/src/pkg/job"
)

func runGetLogs(cmd *cobra.Command, args []string) {

	flag.CommandLine.Parse([]string{})

	glog.V(4).Info("run mlcloud job logs ", args[0])

	task, err := cmd.Flags().GetString("task")
	if err != nil {
		fmt.Println("please specify task (ps/worker)")
		return
	}

	index, err := cmd.Flags().GetInt("index")
	if err != nil {
		fmt.Println("please specify task index")
		return
	}

	opt := &job.LogOpt{}

	opt.Follow, err = cmd.Flags().GetBool("follow")
	if err != nil {
		fmt.Println("input err!")
		fmt.Println(err)
		return
	}

	opt.TimeStamps, err = cmd.Flags().GetBool("timestamps")
	if err != nil {
		fmt.Println("input err!")
		fmt.Println(err)
		return
	}

	opt.LimitBytes, err = cmd.Flags().GetInt64("limit-bytes")
	if err != nil {
		fmt.Println("input err!")
		fmt.Println(err)
		return
	}

	opt.TailLine, err = cmd.Flags().GetInt64("tail")
	if err != nil {
		fmt.Println("input err!")
		fmt.Println(err)
		return
	}

	sinceSeconds, err := cmd.Flags().GetDuration("since")
	if err != nil {
		fmt.Println("input err!")
		fmt.Println(err)
		return
	}
	opt.SinceSeconds = int64(math.Ceil(float64(sinceSeconds) / float64(time.Second)))

	err = job.GetJobLogs(args[0], task, index, opt)
	if err != nil {
		glog.Fatal(err)
	}
}

func NewGetLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [jobname]",
		Short: "get training job's log",
		//Long: `download a file or directory from remote to local`,
		Args: cobra.MinimumNArgs(1),
		Run:  runGetLogs,
	}

	// add children commands here
	cmd.Flags().String("task", "worker", "job task typer(tensorflow: master/ps/worker;mxnet: scheduler/server/worker)")
	cmd.Flags().Int("index", 0, "job task index")
	cmd.Flags().BoolP("follow", "f", false, "Specify if the logs should be streamed.")
	cmd.Flags().Bool("timestamps", false, "Include timestamps on each line in the log output")
	cmd.Flags().Int64("limit-bytes", 0, "Maximum bytes of logs to return. Defaults to no limit.")
	cmd.Flags().Int64("tail", -1, "Lines of recent log file to display. Defaults to -1 with no selector, showing all log lines otherwise 10, if a selector is provided.")
	cmd.Flags().Duration("since", 0, "Only return logs newer than a relative duration like 5s, 2m, or 3h. Defaults to all logs. Only one of since-time / since may be used.")

	return cmd
}
