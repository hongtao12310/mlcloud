package fs

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "github.com/deepinsight/mlcloud/src/utils/log"
    fsutil "github.com/deepinsight/mlcloud/src/cmd/mlcloud/utils"

    "strings"
    "bytes"
    "mime/multipart"
    "io"
)


type uploadChunkResponse struct {
    Err string `json:"err"`
}

func getChunkReader(path string, offset int64) (*os.File, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }

    _, err = f.Seek(offset, 0)
    if err != nil {
        defer Close(f)

        return nil, err
    }

    return f, nil
}

func encodeChunkName(src *Chunk, dst string) string {
    chunk := Chunk{
        Path:   dst,
        Offset: src.Offset,
        Size:   src.Size,
    }

    return chunk.ToURLParam().Encode()
}

func postChunk(chunk *Chunk, dst string) ([]byte, error) {
    log.Debugf("post chunk: %#v", chunk)
    f, err := os.Open(chunk.Path)
    if err != nil {
        log.Errorf("failed to open file: %s", chunk.Path)
        return nil, err
    }

    _, err = f.Seek(chunk.Offset, 0)
    if err != nil {
        log.Errorf("Failed to seek file: %s offset: %d", chunk.Path, chunk.Offset)
        return nil, err
    }

    // close the source file
    defer Close(f)

    // target URI
    uri := fmt.Sprintf("%s/api/v1/fs/storage/chunks", Config.ActiveConfig.Endpoint)
    log.Debugf("Post Chunk URI: %s", uri)

    bodyBuf := &bytes.Buffer{}
    chunkWriter := multipart.NewWriter(bodyBuf)
    if err := chunkWriter.SetBoundary(DefaultMultiPartBoundary); err != nil {
        return nil, err
    }

    chunkName := encodeChunkName(chunk, dst)
    log.Debugf("post chunk - chunk name: %s", chunkName)
    part, err := chunkWriter.CreateFormFile("chunk", chunkName)
    if err != nil {
        return nil, err
    }

    _, err = io.CopyN(part, f, chunk.Size)
    if err != nil {
        return []byte{}, err
    }

    contentType := chunkWriter.FormDataContentType()
    if err = chunkWriter.Close(); err != nil {
        return []byte{}, err
    }

    req, err := fsutil.MakeRequestToken(uri, "POST", bodyBuf, contentType, nil)
    if err != nil {
        return []byte{}, err
    }

    return fsutil.GetResponse(req)

}


func uploadChunks(src string, dst string, diffMeta []ChunkMeta) error {
    if len(diffMeta) == 0 {
        log.Debugf("srcfile: %s and destfile: %s are same\n", src, dst)
        return nil
    }

    for _, meta := range diffMeta {
        log.Debugf("diffMeta:%v\n", meta)

        chunk := Chunk{
            Path:   src,
            Offset: meta.Offset,
            Size:   meta.Len,
        }

        body, err := postChunk(&chunk, dst)
        if err != nil {
            return err
        }

        resp := uploadChunkResponse{}
        if err := json.Unmarshal(body, &resp); err != nil {
            return err
        }

        if len(resp.Err) == 0 {
            continue
        }

        return errors.New(resp.Err)
    }

    return nil
}

func uploadFile(src, dst string, srcFileSize int64) error {

    log.Debugf("touch %s size:%d\n", dst, srcFileSize)

    dstMeta, err := remoteChunkMeta(dst, defaultChunkSize)
    if err != nil && !strings.Contains(err.Error(), StatusText(StatusFileNotFound)) {
        return err
    }

    log.Debugf("dst %s chunkMeta:%#v\n", dst, dstMeta)

    srcMeta, err := GetChunkMeta(src, defaultChunkSize)
    if err != nil {
        return err
    }
    log.Debugf("src %s chunkMeta:%#v\n", src, srcMeta)

    diffMeta, err := GetDiffChunkMeta(srcMeta, dstMeta)
    if err != nil {
        return err
    }
    log.Debugf("diff chunkMeta:%#v\n", diffMeta)

    return uploadChunks(src, dst, diffMeta)
}

// send the stat command to remote server to check if the destination path
// is directory
func isRemoteDir(dstPath string) bool {
    // if dst file not exist. return the error
    // reminder user "create the remote folder first"
    dstStat, err := RemoteStat(NewStatCmd(dstPath))
    if err != nil {
        return false
    }

    if dstStat.IsDir {
        return true
    }

    return false
}

// src: local path
// dst: the path on cloud
func Upload(src, dst string) error {
    // get the absolute path
    var srcAbs string
    var err error
    if !filepath.IsAbs(src) {
        srcAbs, err = filepath.Abs(src)
        if err != nil {
            return err
        }
    } else {
        srcAbs = src
    }

    srcAbs = strings.TrimRight(srcAbs, "/")

    log.Debugf("source path absolution path: %s", srcAbs)

    srcIsDir, err := IsDir(srcAbs)
    if err != nil {
        return err
    }

    if !isRemoteDir(dst) {
        return errors.New(fmt.Sprintf("the destination path: %s should be directory," +
            " you must be create it first\n\n" +
            "\tmlcloud fs mkdir %s\n", dst, dst))
    }
    if !srcIsDir {

        // case 1: source is a file and dst must be an existing directory
        //if !dstExist || !dstStat.IsDir {
        //    return errors.New(fmt.Sprintf("destination directory: %s doesn't exist." +
        //        " you must create it fist\n", dst))
        //}

        // list local file
        lsCmd := NewLsCmd(false, srcAbs)
        srcRet, err := lsCmd.Run()
        if err != nil {
            return err
        }

        srcMeta := srcRet.([]LsResult)[0]

        // case 2: source is a file and dst directory exists
        _, file := filepath.Split(srcAbs)
        realDst := filepath.Join(dst, file)
        log.Debugf("put src_file: %s into dst_path: %s\n", srcAbs, realDst)
        if err := uploadFile(src, realDst, srcMeta.Size); err != nil {
            return err
        }
    } else {
        // extract source path parent
        // for example: we get /mysql from /home/mysql
        srcPrefix, srcFolder := filepath.Split(srcAbs)
        realDst := filepath.Join(dst, srcFolder)

        log.Debugf("srcPrefix: %s, realDst: %s", srcPrefix, realDst)

        // case 5: src is a folder and dst is folder
        // put the src folder files recursively
        // get source file metadata
        lsCmd := NewLsCmd(true, srcAbs)
        srcRet, err := lsCmd.Run()
        if err != nil {
            return err
        }

        var realSrc string
        srcMetas := srcRet.([]LsResult)
        for _, srcMeta := range srcMetas {
            log.Debugf("source meta: %#v", srcMeta)
            if strings.HasPrefix(srcMeta.Path, srcAbs) {
                realSrc = strings.TrimPrefix(srcMeta.Path, srcAbs)
            }

            log.Debugf("real source: %s, real dst: %s", realSrc, realDst)
            // if directory we create it if not exist
            if srcMeta.IsDir {
                dirPath := filepath.Join(realDst, realSrc)
                log.Debugf("mkdir dst path: %s", dirPath)
                cmd := &MkdirCmd {
                    Method: mkdirCmdName,
                    Path: dirPath,
                }
                if _, err := RemoteMkdir(cmd); err != nil {
                    return err
                }
                continue
            }

            // if src is a file
            // create the parent folder if needed
            pDir, _ := filepath.Split(realSrc)
            pDir = filepath.Join(realDst, pDir)
            cmd := &MkdirCmd{
                Method: mkdirCmdName,
                Path: pDir,
            }
            if _, err := RemoteMkdir(cmd); err != nil {
                return err
            }

            // upload file
            dstPath := filepath.Join(realDst, realSrc)

            log.Debugf("upload src_path: %s dst_path: %s\n", srcAbs, dstPath)
            if err := uploadFile(srcMeta.Path, dstPath, srcMeta.Size); err != nil {
                return err
            }
        }
    }

    return nil
}
