package main

import (
    "github.com/spf13/cobra"
    "github.com/deepinsight/mlcloud/src/cmd/mlcloud/fs"
    "github.com/deepinsight/mlcloud/src/cmd/mlcloud/job"
    "fmt"
    "os"
)

var subCommands = []string {
    "fs",
    "job",
}

var rootCmd =  &cobra.Command {
    Use:   "mlcloud [fs|job]",
    Short: "machine learning cloud CLI tool",
    Long: `submit machine learning jobs to kubernetes based cloud platform`,
    Args: cobra.MinimumNArgs(2),
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

func init() {
    rootCmd.AddCommand(job.NewJobCommand())
    rootCmd.AddCommand(fs.NewFSCommand())

}

func main()  {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, )
    }
}



