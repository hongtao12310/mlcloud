package fs

import (
    "github.com/deepinsight/mlcloud/src/cmd/mlcloud/utils"
    "strings"
    "os"
)

// Config is global config object for pfs commandline
var Config, _ = utils.ParseDefaultConfig()

const (
    defaultChunkSize = 2 * 1024 * 1024
    defaultMaxChunkSize = 4 * 1024 * 1024
    defaultMinChunkSize = 1
    //defaultMinChunkSize = 4 * 1024
)

func removeRoot(root, path string) string {
    if root == "/" {
        return path
    }

    if strings.HasPrefix(path, root) {
        return strings.TrimPrefix(path, root)
    }

    return path
}

// check if path exists
func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func IsDir(path string) (bool, error) {
    info, err := os.Stat(path)
    if err != nil {
        return false, err
    }

    return info.IsDir(), nil
}
