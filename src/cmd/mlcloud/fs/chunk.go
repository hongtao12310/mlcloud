package fs

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"github.com/deepinsight/mlcloud/src/utils/log"
	"errors"
)

// Chunk respresents a chunk info.
type Chunk struct {
	Path   string
	Offset int64
	Size   int64
}

// ToURLParam encodes variables to url encoding parameters.
func (p *Chunk) ToURLParam() url.Values {
	parameters := url.Values{}
	parameters.Add("path", p.Path)

	str := fmt.Sprint(p.Offset)
	parameters.Add("offset", str)

	str = fmt.Sprint(p.Size)
	parameters.Add("size", str)

	return parameters
}

func (p *Chunk) checkChunkSize() error {
    if p.Size < defaultMinChunkSize ||
        p.Size > defaultMaxChunkSize {
        return errors.New(StatusText(StatusBadChunkSize))
    }

    return nil
}

// ValidateCloudArgs checks the conditions when running on cloud.
func (p *Chunk) ValidateCloudArgs() error {
    if err := ValidatePfsPath(p.Path); err != nil {
        return err
    }

    return p.checkChunkSize()
}

// ParseChunk get a Chunk struct from path.
// path example:
// 	  path=/test/1.txt&offset=4096&size=4096
func ParseChunk(rawQuery string) (*Chunk, error) {
	cmd := Chunk{}

	values, err := url.ParseQuery(rawQuery)
    if err != nil {
        return nil, err
    }

    offset := values.Get("offset")
    path := values.Get("path")
    if len(offset) == 0 ||
        len(path) == 0 {
        return nil, errors.New(StatusText(StatusInvalidArgs))
    }

    // check chunk size
    chunkStr := values.Get("size")
    cmd.Size = int64(defaultChunkSize)
	if len(chunkStr) > 0 {
        size, err := strconv.ParseInt(chunkStr, 10, 64)
		if err != nil {
			return nil, errors.New(StatusText(StatusBadChunkSize))
		}

        cmd.Size = size
	}

    cmd.Path = path
	cmd.Offset, err = strconv.ParseInt(offset, 10, 64)
	if err != nil {
		return nil, errors.New(StatusText(ParseStrToIntErr))
	}

    log.Debugf("parse chunk: %#v", cmd)
    if err := cmd.ValidateCloudArgs(); err != nil {
        return nil, err
    }

	return &cmd, nil
}

// LoadChunkData loads a specified chunk to io.Writer.
// basically, it will load cloud file chunk to http response body
func (p *Chunk) LoadChunkData(w io.Writer) error {
    // throw exception if file is directory
    fileInfo, err := os.Stat(p.Path)
    if err != nil && os.IsNotExist(err) {
        return errors.New(StatusText(StatusFileNotFound))
    }

    if fileInfo.IsDir() {
        return errors.New( StatusText(StatusDestShouldBeDirectory))
    }

	f, err := os.Open(p.Path)
	if err != nil {
		return err
	}
	defer Close(f)

	_, err = f.Seek(p.Offset, 0)
	if err != nil {
		return err
	}

	loaded, err := io.CopyN(w, f, p.Size)
	log.Debugf("loaded:%d\n", loaded)
    if err != nil && err != io.EOF {
        return err
    }

    return nil
}

// SaveChunkData save data from io.Reader.
func (p *Chunk) SaveChunkData(r io.Reader) error {
    // create file if not exist
	if _, err := os.Stat(p.Path); os.IsNotExist(err) {
        if _, err := os.Create(p.Path); err != nil {
            log.Errorf("failed create file: %s", p.Path)
            return err
        }
    }

	f, err := os.OpenFile(p.Path, os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer Close(f)

	if _, err := f.Seek(p.Offset, 0); err != nil {
        return err
    }

    writen, err := io.CopyN(f, r, p.Size)
    log.Debugf("chunksize: %d writen: %d\n", p.Size, writen)

    if err != nil && err != io.EOF {
        log.Errorf("save chunk data error: %#v", err)
		return err
    }

    return nil
}
