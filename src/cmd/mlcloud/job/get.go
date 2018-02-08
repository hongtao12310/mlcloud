package job

import (
	"github.com/spf13/cobra"
	"github.com/deepinsight/mlcloud/src/pkg/job"
	"github.com/golang/glog"
)

func runGet(cmd *cobra.Command, args []string) {

	glog.V(4).Info("run mlcloud job get ", args[0])

	err := job.GetJob(args[0])
	if err != nil {
		glog.Fatal(err)
	}
}

func NewGetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get [jobname]",
		Short: "get a training job",
		//Long: `download a file or directory from remote to local`,
		Args: cobra.MinimumNArgs(1),
		Run: runGet,
	}

	// add children commands here

	return getCmd
}
