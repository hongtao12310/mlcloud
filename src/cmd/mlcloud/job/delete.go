package job

import (
	"github.com/spf13/cobra"
	"github.com/deepinsight/mlcloud/src/utils/log"
	"github.com/deepinsight/mlcloud/src/pkg/job"
	"github.com/golang/glog"
)

func runDelete(cmd *cobra.Command, args []string) {

	glog.V(4).Info("run mlcloud job delete ", args[0])

	err := job.DeleteJob(args[0])
	if err != nil {
		log.Fatal(err)
	}
}

func NewDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [jobname]",
		Short: "delete job",
		//Long: `download a file or directory from remote to local`,
		Args: cobra.MinimumNArgs(1),
		Run: runDelete,
	}

	// add children commands here

	return cmd
}