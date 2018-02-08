package fs

import (
    "fmt"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/spf13/cobra"
    "os"
)

const (
    GetCmdName = "get"
)

// GetCmdResult means the copy-command's result.
type GetCmdResult struct {
    Src string `json:"Path"`
    Dst string `json:"Dst"`
}

// GetCmd means copy-command.
type GetCmd struct {
    Method string
    V      bool
    Src    string
    Dst    string
}

// PartToString prints command's info.
func (p *GetCmd) PartToString(src, dst string) string {
    if p.V {
        return fmt.Sprintf("get -v %s %s\n", src, dst)
    }
    return fmt.Sprintf("get %s %s\n", src, dst)
}


func (p *GetCmd) CheckInput() error {
    // return err if source file not exist
    if exist, err := PathExists(p.Src); !exist {
        return err
    }

    return nil
}

func runGet(cmd *cobra.Command, args []string) {
    verbose, err := cmd.Flags().GetBool("verbose")
    if err != nil {
        log.Fatal(err)
    }

    getCmd := &GetCmd{
        Method: GetCmdName,
        V: verbose,
        Src: args[0],
        Dst: args[1],
    }

    if err := getCmd.CheckInput(); err != nil {
        fmt.Fprint(os.Stderr, err.Error())
        return
    }

    if err := Download(getCmd.Src, getCmd.Dst); err != nil {
        fmt.Fprint(os.Stderr, err.Error())
        return
    }
}

func NewGetCommand() *cobra.Command {
    getCmd := &cobra.Command{
        Use:   "get [src] [dst]",
        Short: "download file",
        Long: `download a file or directory from remote to local`,
        Args: cobra.MinimumNArgs(2),
        Run: runGet,
    }

    getCmd.Flags().BoolP("verbose", "v", false, "output the details of the get process")

    // add children commands here

    return getCmd
}

