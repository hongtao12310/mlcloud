package fs

import (
    "encoding/json"
    "fmt"
    "net/url"
    "os"
    "github.com/deepinsight/mlcloud/src/cmd/mlcloud/utils"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "path/filepath"
    "github.com/spf13/cobra"
    "errors"
)

const (
    mkdirCmdName = "mkdir"
)

// MkdirResult means Mkdir command's result.
type MkdirResult struct {
    Path string `json:"path"`
}

// MkdirCmd means Mkdir command.
// swagger:parameters mkdirCmd
type MkdirCmd struct {
    // in: body
    Root   string   `json:"root"`
    // method - mkdir
    // in: body
    // required: true
    Method string   `json:"method"`

    // file path
    // in: body
    // required: true
    Path   string   `json:"path"`
}

// ValidateLocalArgs checks the conditions when running on local.
func (p *MkdirCmd) ValidateLocalArgs() error {
    return nil
}

// ValidateCloudArgs checks the conditions when running on cloud.
func (p *MkdirCmd) ValidateCloudArgs() error {
    return ValidatePfsPath(p.Path)
}

// ToURLParam need not to be implemented.
func (p *MkdirCmd) ToURLParam() url.Values {
    panic("not implemented")
}

// ToJSON encodes MkdirCmd to JSON string.
func (p *MkdirCmd) ToJSON() ([]byte, error) {
    return json.Marshal(p)
}

// NewMkdirCmd returns a new MkdirCmd.
func NewMkdirCmd(path string) *MkdirCmd {
    return &MkdirCmd{
        Method: mkdirCmdName,
        Path:  path,
    }
}

// format user mkdir command
func FormatUserMkdirCmd(cmd *MkdirCmd, user *models.User) *MkdirCmd {

    newCmd := MkdirCmd{}
    newCmd.Method = cmd.Method
    newCmd.Root = cmd.Root
    newCmd.Path = filepath.Join(cmd.Root, cmd.Path)

    return &newCmd
}


// Run runs MkdirCmd.
func (p *MkdirCmd) Run() (interface{}, error) {
    var results []MkdirResult
    fi, err := os.Stat(p.Path)

    if os.IsExist(err) && !fi.IsDir() {
        return results, errors.New(
            StatusText(StatusAlreadyExist))
    }

    if err := os.MkdirAll(p.Path, 0755); err != nil {
        return results, err
    }

    results = append(results, MkdirResult{Path: p.Path})

    return results, nil
}

// Name returns name of MkdirComand.
//func (*MkdirCmd) Name() string {
//    return "mkdir"
//}
//
//// Synopsis returns synopsis of MkdirCmd.
//func (*MkdirCmd) Synopsis() string {
//    return "mkdir directoies on Machine Learning Cloud"
//}
//
//// Usage returns usage of MkdirCmd.
//func (*MkdirCmd) Usage() string {
//    return `mkdir <pfspath>:
//	mkdir directories on Machine Learning Cloud
//	Options:
//`
//}
//
//// SetFlags sets MkdirCmd's parameters.
//func (p *MkdirCmd) SetFlags(f *flag.FlagSet) {
//}

func formatMkdirPrint(results []MkdirResult, err error) {
    if err != nil {
        fmt.Println("\t" + err.Error())
    }
}

// RemoteMkdir creat a directory on cloud.
func RemoteMkdir(cmd *MkdirCmd) ([]MkdirResult, error) {
    j, err := cmd.ToJSON()
    if err != nil {
        return nil, err
    }

    uri := fmt.Sprintf("%s/api/v1/fs/files", Config.ActiveConfig.Endpoint)
    log.Debugf("mkdir uri: %s", uri)

    body, _ := utils.PostCall(uri, j)

    log.Debugf("Mkdir Call Body: %s", string(body[:]))

    type mkdirResponse struct {
        Code    int           `json:"code"`
        Err     string        `json:"err"`
        Results []MkdirResult `json:"results"`
    }

    resp := mkdirResponse{}
    if err := json.Unmarshal(body, &resp); err != nil {
        return resp.Results, err
    }

    log.Debugf("%#v\n", resp)

    if len(resp.Err) == 0 {
        return resp.Results, nil
    }

    return resp.Results, errors.New(resp.Err)
}

//func remoteMkdir(cmd *MkdirCmd) error {
//    subcmd := NewMkdirCmd(cmd.Path)
//
//    fmt.Printf("mkdir %s\n", cmd.Path)
//    results, err := RemoteMkdir(subcmd)
//    if err != nil {
//        return err
//    }
//
//    formatMkdirPrint(results, err)
//    return nil
//
//}

// Execute runs a MkdirCmd.
//func (p *MkdirCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
//    if f.NArg() < 1 {
//        f.Usage()
//        return subcommands.ExitFailure
//    }
//
//    cmd, err := newMkdirCmdFromFlag(f)
//    if err != nil {
//        return subcommands.ExitFailure
//    }
//    log.Debugf("%#v\n", cmd)
//
//    if err := remoteMkdir(cmd); err != nil {
//        return subcommands.ExitFailure
//    }
//
//    return subcommands.ExitSuccess
//}

func RunMkDir(cmd *cobra.Command, args []string) {
    mkDirCmd := &MkdirCmd{
        Method: mkdirCmdName,
        Path: args[0],
    }

    _, err := RemoteMkdir(mkDirCmd)
    if err != nil {
        log.Fatal(err)
    }

    //formatMkdirPrint(results, err)

}

func NewMkDirCommand() *cobra.Command {
    mkDirCmd :=  &cobra.Command{
        Use:   "mkdir [path]",
        Short: "make directory",
        Long: `create a directory in remote cloud platform`,
        Args: cobra.MinimumNArgs(1),
        Run: RunMkDir,
    }

    // add children commands here
    return mkDirCmd
}