package fs

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"strconv"
	"github.com/deepinsight/mlcloud/src/pkg/models"
	"path/filepath"
	"errors"
	"github.com/deepinsight/mlcloud/src/utils/log"
)

const (
	// ChunkMetaCmdName is the name of GetChunkMeta command.
	ChunkMetaCmdName = "GetChunkMeta"
)

// ChunkMeta holds the chunk meta's info.
type ChunkMeta struct {
	Offset   int64  `json:"offset"`
	Checksum string `json:"checksum"`
	Len      int64  `json:"len"`
}

// ChunkMetaCmd is a command.
type ChunkMetaCmd struct {
	Method    string `json:"method"`
	Root      string `json:"root"`
	FilePath  string `json:"path"`
	ChunkSize int64  `json:"chunksize"`
}

// ToURLParam encodes ChunkMetaCmd to URL encoding string.
func (p *ChunkMetaCmd) ToURLParam() url.Values {
	parameters := url.Values{}
	parameters.Add("method", p.Method)
	parameters.Add("path", p.FilePath)

	str := fmt.Sprint(p.ChunkSize)
	parameters.Add("chunksize", str)

	return parameters
}

// ToJSON encodes ChunkMetaCmd to JSON string.
func (p *ChunkMetaCmd) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// Run is a functions which run ChunkMetaCmd.
func (p *ChunkMetaCmd) Run() (interface{}, error) {
	return GetChunkMeta(p.FilePath, p.ChunkSize)
}

func (p *ChunkMetaCmd) checkChunkSize() error {
	if p.ChunkSize < defaultMinChunkSize ||
		p.ChunkSize > defaultMaxChunkSize {
		return errors.New(StatusText(StatusBadChunkSize))
	}

	return nil
}

// ValidateCloudArgs checks the conditions when running on cloud.
func (p *ChunkMetaCmd) ValidateCloudArgs() error {
	if err := ValidatePfsPath(p.FilePath); err != nil {
		return err
	}

	return p.checkChunkSize()
}

// ValidateLocalArgs checks the conditions when running locally.
func (p *ChunkMetaCmd) ValidateLocalArgs() error {
	return p.checkChunkSize()
}

// NewChunkMetaCmdFromURLParam get a new ChunkMetaCmd.
func NewChunkMetaCmdFromURLParam(root, rawQuery string, user *models.User) (*ChunkMetaCmd, error ) {
    values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil, err
	}

	method := values.Get("method")
	path := values.Get("path")
	chunkStr := values.Get("chunksize")

	if len(method) == 0 ||
        method != ChunkMetaCmdName ||
		len(path) == 0 {
		return nil, errors.New(StatusText(StatusNotEnoughArgs))
	}

    chunkSize := int64(defaultChunkSize)
    if len(chunkStr) > 0 {
        chunkSize, err = strconv.ParseInt(chunkStr, 10, 64)
        if err != nil {
            return nil, errors.New( StatusText(StatusBadChunkSize))
        }
    }

	cmd :=  &ChunkMetaCmd{
		Method:    method,
		FilePath:  path,
		ChunkSize: chunkSize,
	}

    if err := cmd.ValidateCloudArgs(); err != nil {
        return nil, err
    }

	// the filepath should be rootPath + path
    cmd.FilePath = filepath.Join(root, path)

    return cmd, nil

}

type metaSlice []ChunkMeta

func (a metaSlice) Len() int           { return len(a) }
func (a metaSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a metaSlice) Less(i, j int) bool { return a[i].Offset < a[j].Offset }

// GetDiffChunkMeta gets difference between srcMeta and dstMeta.
func GetDiffChunkMeta(srcMeta []ChunkMeta, dstMeta []ChunkMeta) ([]ChunkMeta, error) {
	if len(dstMeta) == 0 || len(srcMeta) == 0 {
		return srcMeta, nil
	}

	if !sort.IsSorted(metaSlice(srcMeta)) {
		sort.Sort(metaSlice(srcMeta))
	}

	if !sort.IsSorted(metaSlice(dstMeta)) {
		sort.Sort(metaSlice(dstMeta))
	}

	dstIdx := 0
	srcIdx := 0
	diff := make([]ChunkMeta, 0, len(srcMeta))

	for {
		if srcMeta[srcIdx].Offset < dstMeta[dstIdx].Offset {
			diff = append(diff, srcMeta[srcIdx])
			srcIdx++
		} else if srcMeta[srcIdx].Offset > dstMeta[dstIdx].Offset {
			dstIdx++
		} else {
			if srcMeta[srcIdx].Checksum != dstMeta[dstIdx].Checksum {
				diff = append(diff, srcMeta[srcIdx])
			}

			dstIdx++
			srcIdx++
		}

		if dstIdx >= len(dstMeta) {
			break
		}

		if srcIdx >= len(srcMeta) {
			break
		}
	}

	if srcIdx < len(srcMeta) {
		diff = append(diff, srcMeta[srcIdx:]...)
	}

	return diff, nil
}

// GetChunkMeta gets chunk metas from path of file.
func GetChunkMeta(path string, len int64) ([]ChunkMeta, error) {
    // throw exception if file is directory
    fileInfo, err := os.Stat(path)

    log.Debugf("chunk meta data path: %s", path)
    log.Debugf("fileInfo: %v", fileInfo)
    log.Debugf("err: %v", err)

	// TODO: error can be missed here and crash
	if err != nil && os.IsNotExist(err) {
		log.Warning(err.Error())
		return []ChunkMeta{}, errors.New( StatusText(StatusFileNotFound))
	}

    if fileInfo.IsDir() {
        return nil, errors.New(StatusText(StatusDirectoryNotAFile))
    }

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer Close(f)

	if len > defaultMaxChunkSize || len < defaultMinChunkSize {
		return nil, errors.New(StatusText(StatusBadChunkSize))
	}

	var metas []ChunkMeta

	data := make([]byte, len)
	offset := int64(0)

	for {
		n, err := f.Read(data)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err == io.EOF {
			break
		}

		m := ChunkMeta{}
		m.Offset = offset
		sum := md5.Sum(data[:n])
		m.Checksum = hex.EncodeToString(sum[:])
		m.Len = int64(n)

		metas = append(metas, m)

		offset += int64(n)
	}

	return metas, nil
}
