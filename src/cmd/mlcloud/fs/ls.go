package fs

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strconv"
    "time"

    "github.com/deepinsight/mlcloud/src/utils/log"
    fsutil "github.com/deepinsight/mlcloud/src/cmd/mlcloud/utils"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/spf13/cobra"
    "errors"
)

const (
    lsCmdName = "ls"
)

// LsResult represents a LsCmd's result.
type LsResult struct {
    Path    string `json:"Path"`
    ModTime int64  `json:"ModTime"`
    Size    int64  `json:"Size"`
    IsDir   bool   `json:"IsDir"`
}

// LsCmd means LsCmd structure.
// swagger:parameters lsCmd
type LsCmd struct {
    // list file method - ls
    // required: true
    Method string     `json:"method"`

    // recursively
    R      bool       `json:"r"`

    // the root Path
    Root   string     `json:"root"`

    // the real Path
    // required: true
    Path   string     `json:"path"`
}

// ToURLParam encoding LsCmd to URL Encoding string.
func (p *LsCmd) ToURLParam() url.Values {
    parameters := url.Values{}
    parameters.Add("method", p.Method)
    parameters.Add("r", strconv.FormatBool(p.R))
    parameters.Add("path", p.Path)

    return parameters
}

// ToJSON does't need to be implemented.
func (p *LsCmd) ToJSON() ([]byte, error) {
    panic("not implemented")
}

//func newLsCmdFromFlag(f *flag.FlagSet) (*LsCmd, error) {
//    cmd := LsCmd{
//        Method: lsCmdName,
//    }
//
//    var err error
//    f.Visit(func(flag *flag.Flag) {
//        if flag.Name == "r" {
//            cmd.R, err = strconv.ParseBool(flag.Value.String())
//            if err != nil {
//                log.Error("meets error when parsing argument r")
//                return
//            }
//        }
//    })
//
//    if err != nil {
//        return nil, err
//    }
//
//    cmd.Path = f.Arg(0)
//
//    return &cmd, nil
//}

// NewLsCmdFromURLParam returns a new LsCmd according path variable.
func NewLsCmdFromURLParam(root, rawPath string, user *models.User) (*LsCmd, error) {
    cmd := LsCmd{
        Root: root,
    }

    values, err := url.ParseQuery(rawPath)
    if err != nil {
        return nil, err
    }

    log.Debugf("parameters: %+v", values)

    method := values.Get("method")
    r := values.Get("r")
    path := values.Get("path")


    if method != lsCmdName {
        return nil, errors.New(http.StatusText(http.StatusMethodNotAllowed) + ":" + cmd.Method)
    }

    cmd.Method = method
    if len(r) == 0 {
        r = "false"
    }

    if len(path) == 0 {
        path = "/"
    }

    cmd.R, err = strconv.ParseBool(r)
    if err != nil {
        return nil, err
    }

    // format command file path
    cmd.Path = filepath.Join(root, path)
    return &cmd, nil
}

// NewLsCmd return a new LsCmd according r and path variable.
func NewLsCmd(r bool, path string) *LsCmd {
    return &LsCmd{
        Method: lsCmdName,
        R:      r,
        Path:   path,
    }
}

func lsPath(path string, r bool) ([]LsResult, error ) {
    var ret []LsResult

    err := filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {

        if err != nil {
            return err
        }

        log.Debugf("path: %s, subpath: %s", path, subpath)
        m := LsResult{}
        m.Path = subpath
        m.Size = info.Size()
        m.ModTime = info.ModTime().UnixNano()
        m.IsDir = info.IsDir()

        if subpath == path {
            if info.IsDir() {
            } else {
                ret = append(ret, m)
            }
        } else {
            ret = append(ret, m)
        }

        if info.IsDir() && !r && subpath != path {
            return filepath.SkipDir
        }

        return nil
    })

    return ret, err
}

// ValidateCloudArgs checks the conditions when running on cloud.
func (p *LsCmd) ValidateCloudArgs() error {
    return ValidatePfsPath(p.Path)
}

// ValidateLocalArgs checks the conditions when running local.
func (p *LsCmd) ValidateLocalArgs() error {
    if len(p.Path) == 0 {
        return errors.New(StatusText(StatusNotEnoughArgs))
    }
    return nil
}


// Run functions runs LsCmd and return LsResult and error if any happened.
func (p *LsCmd) Run() (interface{}, error ) {
    var results []LsResult

    log.Infof("ls %s\n", p.Path)

    list, err := filepath.Glob(p.Path)
    if err != nil {
        return nil, err
    }

    if len(list) == 0 {
        return results, errors.New(StatusText(StatusFileNotFound))
    }

    for _, path := range list {
        rets, err := lsPath(path, p.R)
        if err != nil {
            return nil, err
        }

        for _, ret := range rets {
            ret.Path = removeRoot(p.Root, ret.Path)
            results = append(results, ret)
        }

    }

    return results, nil
}

//// Name returns LsCmd's name.
//func (*LsCmd) Name() string {
//    return "ls"
//}
//
//// Synopsis returns Synopsis of LsCmd.
//func (*LsCmd) Synopsis() string {
//    return "List files on Machine Learning Cloud"
//}
//
//// Usage returns usage of LsCmd.
//func (*LsCmd) Usage() string {
//    return `ls [-r] <pfspath>:
//	List files on Machine Learning Cloud
//	Options:
//`
//}
//
//// SetFlags sets LsCmd's parameters.
//func (p *LsCmd) SetFlags(f *flag.FlagSet) {
//    f.BoolVar(&p.R, "r", false, "list files recursively")
//}

// getFormatPrint gets max width of filesize and return format string to print.
func getFormatString(result []LsResult) string {
    max := 0
    for _, t := range result {
        str := fmt.Sprintf("%d", t.Size)

        if len(str) > max {
            max = len(str)
        }
    }

    return fmt.Sprintf("%%s %%s %%%dd %%s\n", max)
}

func formatPrint(result []LsResult) {
    formatStr := getFormatString(result)

    for _, t := range result {
        timeStr := time.Unix(0, t.ModTime).Format("2006-01-02 15:04:05")

        if t.IsDir {
            fmt.Printf(formatStr, timeStr, "d", t.Size, t.Path)
        } else {
            fmt.Printf(formatStr, timeStr, "f", t.Size, t.Path)
        }
    }

    fmt.Printf("\n")
}

// RemoteLs gets LsCmd result from cloud.
func RemoteLs(cmd *LsCmd) ([]LsResult, error) {
    t := fmt.Sprintf("%s/api/v1/fs/files", Config.ActiveConfig.Endpoint)
    body, _ := fsutil.GetCall(t, cmd.ToURLParam())

    type lsResponse struct {
        Err     string     `json:"err"`
        Code    int        `json:"code"`
        Results []LsResult `json:"results"`
    }

    resp := lsResponse{}
    if err := json.Unmarshal(body, &resp); err != nil {
        return resp.Results, err
    }

    if len(resp.Err) == 0  {
        return resp.Results, nil
    }

    return resp.Results, errors.New(resp.Err)

}

//func remoteLs(cmd *LsCmd) error {
//    subCmd := NewLsCmd(
//        cmd.R,
//        cmd.Path,
//    )
//    result, err := RemoteLs(subCmd)
//
//    if err != nil {
//        fmt.Printf("  error:%s\n\n", err.Error())
//        return err
//    }
//
//    formatPrint(result)
//    return nil
//}

// Execute runs a LsCmd.
//func (p *LsCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
//    for _, arg := range f.Args() {
//        log.Debugf("ls arg: %s", arg)
//    }
//
//    return subcommands.ExitSuccess
//    if f.NArg() < 1 {
//        f.Usage()
//        return subcommands.ExitFailure
//    }
//
//    cmd, err := newLsCmdFromFlag(f)
//    if err != nil {
//        return subcommands.ExitFailure
//    }
//    log.Debugf("%#v\n", cmd)
//
//    if err := remoteLs(cmd); err != nil {
//        return subcommands.ExitFailure
//    }
//    return subcommands.ExitSuccess
//}

func RunLS(cmd *cobra.Command, args []string) {
    recusive, err := cmd.Flags().GetBool("recursive")
    if err != nil {
        log.Fatal(err)
    }

    lsCmd := &LsCmd{
        Method: lsCmdName,
        R: recusive,
        Path: args[0],
    }

    if result, err := RemoteLs(lsCmd); err != nil {
        log.Fatal(err)
    } else {
        formatPrint(result)
    }


}

func NewLSCommand() *cobra.Command {
    fsCmd :=  &cobra.Command{
        Use:   "ls [src]",
        Short: "list file",
        Long: `list a file or directory metadata`,
        Args: cobra.MinimumNArgs(1),
        Run: RunLS,
    }

    // init flag
    fsCmd.Flags().BoolP("recursive", "r", false, "list a directory recursively")

    // add children commands here
    return fsCmd
}