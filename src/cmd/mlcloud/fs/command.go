package fs

import (
    "io"
    "net/url"
    "strings"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "errors"
)

const (
    // DefaultMultiPartBoundary is the default multipart form boudary.
    DefaultMultiPartBoundary = "8d7b0e5709d756e21e971ff4d9ac3b20"

    // MaxJSONRequestSize is the max body size when server receives a request.
    MaxJSONRequestSize = 2048
)

// Command is a interface of all commands.
type Command interface {
    // ToURLParam generates url.Values of the command struct.
    ToURLParam() url.Values
    // ToJSON generates JSON string of the command struct.
    ToJSON() ([]byte, error)
    // Run runs a command.
    Run() (interface{}, error)
    // ValidateLocalArgs validates arguments when running locally.
    ValidateLocalArgs() error
    // ValidateCloudArgs validates arguments when running on cloud.
    ValidateCloudArgs() error
}

type BaseCommand struct {
    Method string   `json:"method"`
}


// ValidatePfsPath returns whether a path is a pfspath.
func ValidatePfsPath(path string) error {
    if !strings.HasPrefix(path, "/") {
        return errors.New(StatusText(StatusShouldBePfsPath))
    }
    return nil
}

// IsCloudPath returns whether a path is a pfspath.
func IsCloudPath(path string) bool {
    return strings.HasPrefix(path, "/fs/")
}

// Close closes c and log it.
func Close(c io.Closer) {
    err := c.Close()
    if err != nil {
        log.Error(err)
    }
}
