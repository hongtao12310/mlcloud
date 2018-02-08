package fs

import (
    "net/url"
    "os"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "path/filepath"
    "fmt"
    "github.com/deepinsight/mlcloud/src/cmd/mlcloud/utils"
    "encoding/json"
    "errors"
)

const (
    // StatCmdName means stat command name.
    StatCmdName = "stat"
)

// StatCmd means stat command.
type StatCmd struct {
    Method string   `json:"method"`
    Path   string   `json:"path"`
    // Base indicate the base dir for the user in
    // the distribute file system
    Root   string   `json:"root"`
}

// ToURLParam encodes StatCmd to URL Encoding string.
func (p *StatCmd) ToURLParam() url.Values {
    parameters := url.Values{}
    parameters.Add("method", p.Method)
    parameters.Add("path", p.Path)

    return parameters

}

// ToJSON here need not tobe implemented.
func (p *StatCmd) ToJSON() ([]byte, error) {
    return nil, nil
}

// NewStatCmdFromURLParam return a new StatCmd.
func NewStatCmdFromURLParam(root, rawQuery string, user *models.User) (*StatCmd, error) {
    cmd := StatCmd{
        Root: root,
    }

    values, err := url.ParseQuery(rawQuery)

    if err != nil {
        return nil, errors.New(StatusText(BadRawQueryURL))
    }

    log.Debugf("parameters: %+v", values)

    method := values.Get("method")
    path := values.Get("path")

    if method != StatCmdName {
        return nil, errors.New("illegal stat command")
    }

    cmd.Method = method
    if len(path) == 0 {
        path = "/"
    }

    cmd.Path = filepath.Join(root, path)
    return &cmd, nil
}

// ValidateLocalArgs checks the condition when running local.
func (p *StatCmd) ValidateLocalArgs() error {
    panic("not implement")
}

// ValidateCloudArgs checks the conditions when running on cloud.
func (p *StatCmd) ValidateCloudArgs() error {
    return ValidatePfsPath(p.Path)
}

func RemoteStat(cmd *StatCmd) (LsResult, error) {
    t := fmt.Sprintf("%s/api/v1/fs/files", Config.ActiveConfig.Endpoint)
    body, _ := utils.GetCall(t, cmd.ToURLParam())
    //if err != nil {
    //    return LsResult{}, err
    //}

    log.Debugf("Stat Call Body: %s", string(body[:]))

    type statResponse struct {
        Err     string            `json:"err"`
        Code    int               `json:"code"`
        Results LsResult          `json:"results"`
    }

    resp := statResponse{}
    if err := json.Unmarshal(body, &resp); err != nil {
        return resp.Results, errors.New( StatusText(StatusJSONErr))
    }

    if len(resp.Err) == 0  {
        return resp.Results, nil
    }

    //"stat", cmd.Path, resp.Err
    return resp.Results, errors.New(resp.Err)
}

// Run runs the StatCmd.
func (p *StatCmd) Run() (interface{}, error) {
    fi, err := os.Stat(p.Path)
    if err != nil && os.IsNotExist(err) {
        log.Errorf("stat command file: %s not exist", p.Path)
        return nil, errors.New(StatusText(StatusFileNotFound))
    }

    return &LsResult{
        Path:    removeRoot(p.Root, p.Path),
        ModTime: fi.ModTime().UnixNano(),
        IsDir:   fi.IsDir(),
        Size:    fi.Size(),
    }, nil
}

// NewLsCmd return a new LsCmd according r and path variable.
func NewStatCmd(path string) *StatCmd {
    return &StatCmd{
        Method: StatCmdName,
        Path:   path,
    }
}