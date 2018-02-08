package fs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
    "github.com/deepinsight/mlcloud/src/cmd/mlcloud/utils"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "path/filepath"
	"errors"
)

const (
	defaultMaxCreateFileSize = int64(10 * 1024 * 1024 * 1024)
)

const (
	// TouchCmdName is the name of touch command.
	TouchCmdName = "touch"
)

// TouchResult represents touch-command's result.
type TouchResult struct {
	Path string `json:"path"`
}

// TouchCmd is holds touch command's variables.
type TouchCmd struct {
	Method   string `json:"method"`
	FileSize int64  `json:"size"`
	Path     string `json:"path"`
	Root     string `json:"root"`
}

func (p *TouchCmd) checkFileSize() error {
	if p.FileSize < 0 || p.FileSize > defaultMaxCreateFileSize {
		return errors.New(StatusText(StatusBadFileSize) + ":" + fmt.Sprint(p.FileSize))
	}
	return nil
}

// ValidateLocalArgs check the conditions when running local.
func (p *TouchCmd) ValidateLocalArgs() error {
	return p.checkFileSize()
}

// ValidateCloudArgs checks the conditions when running on cloud.
func (p *TouchCmd) ValidateCloudArgs() error {
	if err := ValidatePfsPath(p.Path); err != nil {
		return err
	}

	return p.checkFileSize()
}

// ToURLParam encodes a TouchCmd to a URL encoding string.
func (p *TouchCmd) ToURLParam() url.Values {
	parameters := url.Values{}
	parameters.Add("method", p.Method)
	parameters.Add("path", p.Path)

	str := fmt.Sprint(p.FileSize)
	parameters.Add("size", str)

	return parameters
}

// ToJSON encodes a TouchCmd to a JSON string.
func (p *TouchCmd) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

func remoteTouch(cmd *TouchCmd) error  {
	j, err := cmd.ToJSON()
	if err != nil {
		return err
	}

	t := fmt.Sprintf("%s/api/v1/fs/files", Config.ActiveConfig.Endpoint)
	body, _ := utils.PostCall(t, j)

	log.Debugf("Touch Call Body: %s", string(body[:]))

	type touchResponse struct {
		Code    int         `json:"code"`
		Err     string      `json:"err"`
		Results TouchResult `json:"results"`
	}

	resp := touchResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}

	if len(resp.Err) == 0 {
		return nil
	}

	return errors.New(resp.Err)
}

// format user touch command
func FormatUserTouchCmd(cmd *TouchCmd) *TouchCmd {
    newCmd := TouchCmd{}
    newCmd.Method = cmd.Method
    newCmd.Root = cmd.Root
    newCmd.Path = filepath.Join(cmd.Root, cmd.Path)

    return &newCmd
}

// CreateSizedFile creates a file with specified size.
func CreateSizedFile(path string, size int64) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer Close(fd)

	if size <= 0 {
		return nil
	}

	_, err = fd.Seek(size-1, 0)
	if err != nil {
		return err
	}

	_, err = fd.Write([]byte{0})
	return err
}

// Run is a function runs TouchCmd.
func (p *TouchCmd) Run() (interface{}, error) {
	if p.FileSize < 0 || p.FileSize > defaultMaxCreateFileSize {
		return nil, errors.New(StatusText(StatusBadFileSize))
	}

	fi, err := os.Stat(p.Path)
	if os.IsExist(err) && fi.IsDir() {
		return nil, errors.New(StatusText(StatusDirectoryAlreadyExist))
	}

	if os.IsNotExist(err) || fi.Size() != p.FileSize {
		if err := CreateSizedFile(p.Path, p.FileSize); err != nil {
			return nil, errors.New(StatusText(CreateFileErr))
		}
	}

	return &TouchResult{
		Path: p.Path,
	}, nil
}
