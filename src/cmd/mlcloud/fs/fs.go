package fs

import (
    "github.com/spf13/cobra"
    "fmt"
    "os"
)

func NewFSCommand() *cobra.Command {
    subCommands := []string {
        "ls",
        "rm",
        "mkdir",
        "put",
        "get",
    }

    fsCmd :=  &cobra.Command {
        Use:   "fs [ls|rm|mkdir|put|get]",
        Short: "file system operations",
        Long: `list/remove/mkdir/upload/download file operation`,
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

    // add children commands here
    fsCmd.AddCommand(NewLSCommand())
    fsCmd.AddCommand(NewMkDirCommand())
    fsCmd.AddCommand(NewRMCommand())
    fsCmd.AddCommand(NewPutCommand())
    fsCmd.AddCommand(NewGetCommand())

    return fsCmd
}