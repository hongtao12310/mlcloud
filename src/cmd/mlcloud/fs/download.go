package fs

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "github.com/deepinsight/mlcloud/src/cmd/mlcloud/utils"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "strings"
)

func remoteChunkMeta(path string, chunkSize int64) ([]ChunkMeta, error) {
    cmd := ChunkMetaCmd {
        Method:    ChunkMetaCmdName,
        FilePath:  path,
        ChunkSize: chunkSize,
    }

    t := fmt.Sprintf("%s/api/v1/fs/chunks", Config.ActiveConfig.Endpoint)
    ret, err := utils.GetCall(t, cmd.ToURLParam())
    if err != nil {
        return nil, err
    }

    type chunkMetaResponse struct {
        Err     string      `json:"err"`
        Results []ChunkMeta `json:"results"`
    }

    resp := chunkMetaResponse{}
    if err := json.Unmarshal(ret, &resp); err != nil {
        return nil, err
    }

    if len(resp.Err) == 0 {
        return resp.Results, nil
    }

    return []ChunkMeta{}, errors.New(resp.Err)
}

func getChunkData(target string, chunk Chunk, dst string) error {
    log.Debugf("getChunkData: %+v", chunk)

    resp, err := utils.GetChunk(target, chunk.ToURLParam())
    if err != nil {
        return err
    }
    defer Close(resp.Body)

    if resp.Status != utils.HTTPOK {
        return errors.New("http server returned non-200 status: " + resp.Status)
    }

    partReader := multipart.NewReader(resp.Body, DefaultMultiPartBoundary)
    for {
        part, error := partReader.NextPart()
        if error == io.EOF {
            break
        }

        if part.FormName() == "chunk" {
            recvCmd, err := ParseChunk(part.FileName())
            if err != nil {
                return err
            }

            recvCmd.Path = dst

            if err := recvCmd.SaveChunkData(part); err != nil {
                return err
            }
        }
    }

    return nil
}

func downloadChunks(src, dst string, diffMeta []ChunkMeta) error {
    if len(diffMeta) == 0 {
        log.Debugf("srcfile:%s and dstfile:%s are already same\n", src, dst)
        return nil
    }

    uri := fmt.Sprintf("%s/api/v1/fs/storage/chunks", Config.ActiveConfig.Endpoint)
    for _, meta := range diffMeta {
        chunk := Chunk{
            Path:   src,
            Offset: meta.Offset,
            Size:   meta.Len,
        }

        err := getChunkData(uri, chunk, dst)
        if err != nil {
            return err
        }
    }

    return nil
}

// src: source file path
// srcFileSize
// dst: destination path
func downloadFile(src string, srcFileSize int64, dst string) error {
    srcMeta, err := remoteChunkMeta(src, defaultChunkSize)
    if err != nil {
        return err
    }
    log.Debugf("download file srcMeta:%#v\n\n", srcMeta)

    dstMeta, err := GetChunkMeta(dst, defaultChunkSize)
    if err != nil && !strings.Contains(err.Error(),
        StatusText(StatusFileNotFound)) {
        return err
    }

    log.Debugf("download file dstMeta:%#v\n", dstMeta)

    diffMeta, err := GetDiffChunkMeta(srcMeta, dstMeta)
    if err != nil {
        return err
    }

    log.Debugf("chunk meta diff: %#v", diffMeta)

    err = downloadChunks(src, dst, diffMeta)
    return err
}

// check destination path is file or directory
// when download multiple files from a folder the destination must be
// a directory
func checkBeforeDownLoad(src []LsResult, dst string) (bool, error) {
    var bDir bool
    fi, err := os.Stat(dst)
    if err == nil {
        bDir = fi.IsDir()
        if !fi.IsDir() && len(src) > 1 {
            return bDir, errors.New(StatusText(StatusDestShouldBeDirectory))
        }
    } else if os.IsNotExist(err) {
        return false, nil
    }

    return bDir, err
}


// src: the path on cloud
// dst: local path
func Download(src, dst string) error {
    dstAbs := dst
    var err error
    if !filepath.IsAbs(dst) {
        dstAbs, err = filepath.Abs(dst)
        if err != nil {
            return err
        }
    } else {
        dstAbs = dst
    }

    log.Debugf("download %s to %s\n", src, dstAbs)

    dstIsDir, err := IsDir(dstAbs)
    if err != nil && !os.IsNotExist(err) {
        return err
    }

    statRlt, err := RemoteStat(NewStatCmd(src))
    if err != nil {
        return err
    }

    if !statRlt.IsDir {
        // case 1: src is a file and dst is directory
        if dstIsDir {
            _, file := filepath.Split(statRlt.Path)
            realDst := dstAbs + "/" + file
            log.Debugf("download src_path:%s dst_path:%s\n", statRlt.Path, realDst)
            if err := downloadFile(statRlt.Path, statRlt.Size, realDst); err != nil {
                return err
            }
        } else {
            // case 2: src is a file and dst is not exist
            // we download the file into current folder
            //realDst := filepath.Join(filepath.Dir(os.Args[0]), dstAbs)
            realDst := dstAbs
            log.Debugf("download src_path:%s dst_path:%s\n", statRlt.Path, realDst)

            if err := downloadFile(statRlt.Path, statRlt.Size, realDst); err != nil {
                return err
            }
        }

    } else {
        // case 3: src is a folder and dst not exist
        // new exception "folder not exist"
        if !dstIsDir {
            errInfo := fmt.Sprintf("destination directory: %s not exist. " +
                "please create it first", dstAbs)
            return errors.New(errInfo)
        }

        // case 4: src is a folder and dst is folder
        // download the file recursively

        // mlcloud fs get /test/test1/test2 /home
        // we will download the test2 folder to /home/test2
        // ls -l /home
        // 2017-09-28 10:30:45 d      102 /home/test2/Users
        // 2017-09-28 10:30:45 d      102 /home/test2/Users/hongtaozhang
        // 2017-09-28 12:26:18 d      306 /home/test2/Users/hongtaozhang/mysql
        // 2017-09-28 12:26:17 f       56 /home/test2/Users/hongtaozhang/mysql/auto.cnf
        // 2017-09-28 12:26:17 d      170 /home/test2/Users/hongtaozhang/mysql/demo
        // 2017-09-28 12:26:17 f       65 /home/test2/Users/hongtaozhang/mysql/demo/db.opt

        // [hongtaozhang@HongtaodeMacBook-Pro mlcloud]$ ./mlcloud fs ls -r /test/test1/test2
        // 2017-09-28 10:30:45 d      102 /test/test1/test2/Users
        // 2017-09-28 10:30:45 d      102 /test/test1/test2/Users/hongtaozhang
        // 2017-09-28 12:26:18 d      306 /test/test1/test2/Users/hongtaozhang/mysql
        // 2017-09-28 12:26:17 f       56 /test/test1/test2/Users/hongtaozhang/mysql/auto.cnf
        // 2017-09-28 12:26:17 d      170 /test/test1/test2/Users/hongtaozhang/mysql/demo
        // 2017-09-28 12:26:17 f       65 /test/test1/test2/Users/hongtaozhang/mysql/demo/db.opt

        // 1. split the source folder '/test/test1/test2' into two parts
        //    parent folder = /test/test1
        //    src folder = /test2
        // 2. the real destination folder is dstAbs + srcFolder
        // 3. list all the source files include folder
        //    1. if source is directory
        //       srcChildrenPath = trimprefix(srcMeta.path, realSrc)
        //       dstPath = filepath.join(realDstFolder, srcChildrenPath)
        //       mkdir(dstPath)
        //    2. if source is a file
        //       srcChildrenPath = trimprefix(srcMeta.path, realSrc)
        //       dstPath = filepath.join(realDstFolder, srcChildrenPath)
        //       dstFolder, _ = filepath.split(dstPath)
        //       mkdir(dstFolder)
        //       downloadFile(srcMeta.path, srcMeta.size, dstPath)

        var realSrc, realDstFolder string
        if src == "/" {
            realSrc = "/"
            realDstFolder = dstAbs
        } else {
            realSrc = strings.TrimRight(src, "/")
            _, srcFolder := filepath.Split(realSrc)
            realDstFolder = filepath.Join(dstAbs, srcFolder)
        }

        // list all the files
        lsRet, err := RemoteLs(NewLsCmd(true, src))
        if err != nil {
            return err
        }

        log.Debugf("ls returned: %#v", lsRet)
        var srcChildrenPath, dstPath string
        for _, srcMeta := range lsRet {
            log.Debugf("source meta data: %#v", srcMeta)
            if strings.HasPrefix(srcMeta.Path, realSrc) {
                srcChildrenPath = strings.TrimPrefix(srcMeta.Path, realSrc)
            }

            dstPath = filepath.Join(realDstFolder, srcChildrenPath)


            // if directory we create it if not exist
            if srcMeta.IsDir {
                log.Debugf("real source: %s, mkdir real dst: %s", srcMeta.Path, dstPath)
                if _, err := os.Stat(dstPath); os.IsNotExist(err) {
                    os.MkdirAll(dstPath, 0755)
                }
                continue
            }

            // if src is a file
            // create the parent folder if needed
            pDir, _ := filepath.Split(dstPath)
            if _, err := os.Stat(pDir); os.IsNotExist(err) {
                os.MkdirAll(pDir, 0755)
            }

            log.Debugf("download src_path: %s dst_path: %s\n", srcMeta.Path, dstPath)
            if err := downloadFile(srcMeta.Path, srcMeta.Size, dstPath); err != nil {
                return err
            }
        }
    }

    return nil
}
