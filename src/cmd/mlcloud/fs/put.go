package fs

import (
    "fmt"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/spf13/cobra"
    "os"
)

const (
    PutCmdName        = "put"
)

// PutCmdResult means the copy-command's result.
type PutCmdResult struct {
    Src string `json:"Path"`
    Dst string `json:"Dst"`
}

// PutCmd means copy-command.
type PutCmd struct {
    Method string
    V      bool
    Src    string
    Dst    string
}

// PartToString prints command's info.
func (p *PutCmd) PartToString(src, dst string) string {
    if p.V {
        return fmt.Sprintf("put -v %s %s\n", src, dst)
    }
    return fmt.Sprintf("put %s %s\n", src, dst)
}


func (p *PutCmd) CheckInput() error {
    // return err if source file not exist
    if exist, err := PathExists(p.Src); !exist {
        return err
    }

    return nil
}

// RunPut runs PutCmd.
func RunPut(cmd *cobra.Command, args []string) {
    verbose, err := cmd.Flags().GetBool("verbose")
    if err != nil {
        log.Fatal(err)
    }

    putCmd := &PutCmd{
        Method: PutCmdName,
        V: verbose,
        Src: args[0],
        Dst: args[1],
    }

    if err := putCmd.CheckInput(); err != nil {
        fmt.Fprint(os.Stderr, err.Error())
        return
    }

    if err := Upload(putCmd.Src, putCmd.Dst); err != nil {
        fmt.Fprint(os.Stderr, err.Error())
        return
    }
}

func NewPutCommand() *cobra.Command {
    putCmd := &cobra.Command{
        Use:   "put [src] [dst]",
        Short: "upload file",
        Long: `upload a file or directory from local to remote directory`,
        Args: cobra.MinimumNArgs(2),
        Example: "mlcloud fs put ~/test.txt /",
        Run: RunPut,
    }

    putCmd.Flags().BoolP("verbose", "v", false, "output the details of the put process")

    // add children commands here

    return putCmd
}