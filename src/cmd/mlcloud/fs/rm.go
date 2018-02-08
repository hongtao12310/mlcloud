package fs

import (
    "encoding/json"
    "fmt"
    "net/url"
    "os"
    "path/filepath"
    "github.com/deepinsight/mlcloud/src/cmd/mlcloud/utils"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/spf13/cobra"
    "errors"
)

const (
    rmCmdName = "rm"
)

// RmResult means Rm-command's result.
type RmResult struct {
    Path string `json:"path"`
}

// RmCmd means Rm command.
// swagger:parameters rmCmd
type RmCmd struct {
    // list file method - rm
    // in: body
    // required: true
    Method string       `json:"method"`
    // recursively - true or false
    // in: body
    R      bool         `json:"r"`
    // path root
    // in: body
    Root   string       `json:"root"`
    // file full path
    // in: body
    // required: true
    Path   string    `json:"path"`
}

// ValidateLocalArgs checks the conditions when running local.
func (p *RmCmd) ValidateLocalArgs() error {
    if len(p.Path) == 0 {
        return errors.New(StatusText(StatusInvalidArgs))
    }
    return nil
}

// ValidateCloudArgs checks the conditions when running on cloud.
func (p *RmCmd) ValidateCloudArgs() error {
    return ValidatePfsPath(p.Path)
}

// ToURLParam needs not to be implemented.
func (p *RmCmd) ToURLParam() url.Values {
    panic("not implemented")
}

// ToJSON encodes RmCmd to JSON string.
func (p *RmCmd) ToJSON() ([]byte, error) {
    return json.Marshal(p)
}

// NewRmCmd returns a new RmCmd.
func NewRmCmd(r bool, path string) *RmCmd {
    return &RmCmd{
        Method: rmCmdName,
        R:      r,
        Path:   path,
    }
}

func FormatUserRmCmd(cmd *RmCmd, user *models.User) *RmCmd {
    newCmd := RmCmd{}
    newCmd.Method = cmd.Method
    newCmd.R = cmd.R
    newCmd.Root = cmd.Root
    newCmd.Path = filepath.Join(cmd.Root, cmd.Path)
    return &newCmd
}


// Run runs RmCmd.
func (p *RmCmd) Run() (interface{}, error) {
    var result []RmResult

    list, err := filepath.Glob(p.Path)
    if err != nil {
        return nil, err
    }

    for _, arg := range list {
        fi, err := os.Stat(arg)
        if err != nil && os.IsNotExist(err) {
            log.Errorf("stat command file: %s not exist", p.Path)
            return nil, errors.New(StatusText(StatusFileNotFound))
        }

        if fi.IsDir() && !p.R {
            return result, errors.New(StatusText(StatusCannotDelDirectory) + ":" + arg)
        }

        if err := os.RemoveAll(arg); err != nil {
            return result, nil
        }

        result = append(result, RmResult{Path: arg})
    }

    return result, nil
}


func formatRmPrint(results []RmResult, err error) {
    for _, result := range results {
        fmt.Printf("rm %s\n", result.Path)
    }

    if err != nil {
        fmt.Println("\t" + err.Error())
    }

    return
}

// RemoteRm gets RmCmd Result from cloud.
func RemoteRm(cmd *RmCmd) ([]RmResult, error) {
    j, err := cmd.ToJSON()
    if err != nil {
        return nil, err
    }

    uri := fmt.Sprintf("%s/api/v1/fs/files", Config.ActiveConfig.Endpoint)
    body, _ := utils.DeleteCall(uri, j)

    log.Debugf("DELETE Call Body: %s", string(body[:]))

    type rmResponse struct {
        Code    int           `json:"code"`
        Err     string        `json:"err"`
        Results []RmResult    `json:"path"`
    }

    resp := rmResponse{}
    if err := json.Unmarshal(body, &resp); err != nil {
        return resp.Results, err
    }

    log.Debugf("Delete Call Response %#v\n", resp)

    if len(resp.Err) == 0 {
        return resp.Results, nil
    }

    return resp.Results, errors.New(resp.Err)

}


func RunRM(cmd *cobra.Command, args []string) {
    recursive, err := cmd.Flags().GetBool("recursive")
    if err != nil {
        log.Fatal(err)
    }

    rmCmd := &RmCmd{
        Method: rmCmdName,
        R: recursive,
        Path: args[0],
    }

    if _, err := RemoteRm(rmCmd); err != nil {
        log.Fatal(err)
    }
}

func NewRMCommand() *cobra.Command {
    rmCmd :=  &cobra.Command{
        Use:   "rm [flags] [src]",
        Short: "remove file or directory",
        Long: `remove a file or directory from remote cloud`,
        Args: cobra.MinimumNArgs(1),
        Run: RunRM,
    }

    // init flag
    rmCmd.Flags().BoolP("recursive", "r", false, "remove a directory")

    // add children commands here
    return rmCmd
}