package job

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"flag"
)

func NewJobCommand() *cobra.Command {
	subCommands := []string {
		"list",
		"get",
		"submit",
		"delete",
		"logs",
	}

	jobCmd :=  &cobra.Command {
		Use:   "job [list|get|submit|delete|logs]",
		Short: "training job operations",
		Long: `training job operations`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, command := range subCommands {
				if args[0] == command {
					return
				}
			}

			fmt.Fprint(os.Stderr, "no matched subcommand\n\n")
			cmd.Usage()
		},
	}

	// add children commands here
	jobCmd.AddCommand(NewGetCommand())
	jobCmd.AddCommand(NewListCommand())
	jobCmd.AddCommand(NewDeleteCommand())
	jobCmd.AddCommand(NewSubmitCommand())
	jobCmd.AddCommand(NewGetLogsCommand())

	jobCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return jobCmd
}