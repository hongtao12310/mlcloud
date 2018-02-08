package job

import (
	"github.com/spf13/cobra"
	"github.com/deepinsight/mlcloud/src/utils/log"
	"github.com/deepinsight/mlcloud/src/pkg/job"
)

func runList(cmd *cobra.Command, args []string) {

	log.Debug("run mlcloud job list")

	err := job.ListJobs()
	if err != nil {
		log.Fatal(err)
	}
}

func NewListCommand() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list [src] [dst]",
		Short: "list jobs",
		Long: `list jobs`,
		//Args: cobra.MinimumNArgs(2),
		Run: runList,
	}

	// add children commands here

	return listCmd
}