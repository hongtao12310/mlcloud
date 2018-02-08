package api

import (
    "net/http"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/gorilla/mux"
    "encoding/json"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "mime/multipart"
    "io"
    httputil "github.com/deepinsight/mlcloud/src/utils/http"
    fscmd "github.com/deepinsight/mlcloud/src/cmd/mlcloud/fs"
    "github.com/deepinsight/mlcloud/src/pkg/common/syscfg"
    "github.com/deepinsight/mlcloud/src/pkg/common"
    "path/filepath"
)

type FSHandler struct {
    *BaseHandler
}


// NewFSHandler created FSHandler instance.
func NewFSHandler(baseHandler *BaseHandler) *FSHandler {
    return &FSHandler{
        BaseHandler: baseHandler,
    }
}

func (self *FSHandler) Register(router *mux.Router) {
    bs := router.PathPrefix("/api/v1/fs").Subrouter()
    bs.HandleFunc("/files", self.GetFilesHandler).
        Methods("GET")

    bs.HandleFunc("/files", self.PostFilesHandler).
        Methods("GET", "POST")

    bs.HandleFunc("/files", self.DeleteFilesHandler).
        Methods("DELETE")

    bs.HandleFunc("/chunks", self.GetChunkMetaHandler).
        Methods("GET")

    bs.HandleFunc("/storage/chunks", self.GetChunkHandler).
        Methods("GET")

    bs.HandleFunc("/storage/chunks", self.PostChunkHandler).
        Methods("POST")

}

// swagger:route GET /api/v1/fs/files getFiles lsCmd
//
// Handler for files list/stat.
//
// the response body can be fetched from response body
//
// Responses:
// 		        200: response
// 		        400: response
// 		        405: response
// 		        500: response
func (self *FSHandler) GetFilesHandler(writer http.ResponseWriter, request *http.Request) {
    response := httputil.Response{}
    var command fscmd.Command
    var err error

    user, err := self.getUserInfo(writer, request)
    if err != nil {
        response.Err = err.Error()
        response.Code = http.StatusUnauthorized
        httputil.WriteJSONResponse(writer, response)
        return
    }

    rawQuery := request.URL.RawQuery
    method := request.URL.Query().Get("method")

    config := syscfg.GetSysConfig()
    root := config.Get(common.FSBasePath).(string)

    if !user.Sysadmin {
        root = filepath.Join(root, user.Username)
    }

    switch method {
    case "ls":
        command, err = fscmd.NewLsCmdFromURLParam(root, rawQuery, user)
    case "md5sum":
    // TODO
    // err := md5Handler(w, r)
    case "stat":
        command, err = fscmd.NewStatCmdFromURLParam(root, rawQuery, user)
    default:
        response.Code = http.StatusMethodNotAllowed
        httputil.WriteJSONResponse(writer, response)
        return
    }

    if err != nil {
        response.Err = err.Error()
        response.Code = http.StatusInternalServerError
        httputil.WriteJSONResponse(writer, response)
        return
    }

    result, err := command.Run()
    if err != nil {
        response.Code = http.StatusInternalServerError
        response.Err = err.Error()
        httputil.WriteJSONResponse(writer, response)
        return
    }

    response.Results = result
    response.Code = http.StatusOK
    httputil.WriteJSONResponse(writer, response)

}

func rmHandler(writer http.ResponseWriter, user *models.User, body []byte) {
    log.Infof("begin proc rmHandler\n")
    cmd := &fscmd.RmCmd{}

    resp := httputil.Response{}
    if err := json.Unmarshal(body, &cmd); err != nil {
        resp.Code = http.StatusMethodNotAllowed
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    config := syscfg.GetSysConfig()
    root := config.Get(common.FSBasePath).(string)

    if !user.Sysadmin {
        root = filepath.Join(root, user.Username)
    }

    cmd.Root = root
    cmd = fscmd.FormatUserRmCmd(cmd, user)

    result, err := cmd.Run()
    if err != nil {
        resp.Code = http.StatusInternalServerError
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    resp.Results = result
    resp.Code = http.StatusOK
    httputil.WriteJSONResponse(writer, resp)
    log.Infof("end proc handler\n")
    return
}

func mkdirHandler(writer http.ResponseWriter, user *models.User, body []byte) {
    log.Debugf("begin proc mkdir\n")
    cmd := &fscmd.MkdirCmd{}

    resp := httputil.Response{}
    if err := json.Unmarshal(body, &cmd); err != nil {
        resp.Code = http.StatusMethodNotAllowed
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    config := syscfg.GetSysConfig()
    root := config.Get(common.FSBasePath).(string)

    if !user.Sysadmin {
        root = filepath.Join(root, user.Username)
    }
    cmd.Root = root

    cmd = fscmd.FormatUserMkdirCmd(cmd, user)

    result, err := cmd.Run()
    if err != nil {
        resp.Code = http.StatusInternalServerError
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    resp.Results = result
    resp.Code = http.StatusOK
    httputil.WriteJSONResponse(writer, resp)

    log.Debugf("end proc mkdir\n")

    return
}

func touchHandler(writer http.ResponseWriter, user *models.User,  body []byte) {
    log.Infof("begin proc touch\n")
    cmd := &fscmd.TouchCmd{}

    resp := httputil.Response{}
    if err := json.Unmarshal(body, &cmd); err != nil {
        resp.Code = http.StatusMethodNotAllowed
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    config := syscfg.GetSysConfig()
    root := config.Get(common.FSBasePath).(string)

    if !user.Sysadmin {
        root = filepath.Join(root, user.Username)
    }
    cmd.Root = root

    cmd = fscmd.FormatUserTouchCmd(cmd)

    result, err := cmd.Run()
    if err != nil {
        resp.Code = http.StatusInternalServerError
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    resp.Results = result
    resp.Code = http.StatusOK
    httputil.WriteJSONResponse(writer, resp)
    return
    log.Infof("end proc touch\n")
}

// swagger:route POST /api/v1/fs/files postFiles mkdirCmd
//
// Handler for mkdir.
//
// the response body can be fetched from response body
//
// Responses:
// 		        200: response
// 		        400: response
// 		        401: response
func (self *FSHandler) PostFilesHandler(writer http.ResponseWriter, request *http.Request) {
    resp := httputil.Response{}

    user, err := self.getUserInfo(writer, request)
    if err != nil {
        resp.Code = http.StatusUnauthorized
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    body, err := httputil.ReadBody(request)
    if err != nil {
        resp.Code = http.StatusBadRequest
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
    }

    method, err := httputil.GetKeyFromBody(body, "method")
    if err != nil {
        resp.Code = http.StatusBadRequest
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    log.Debugf("method: %v", method)

    switch method.(string) {
    case "mkdir":
        mkdirHandler(writer, user, body)
    case "touch":
        touchHandler(writer, user, body)
    default:
        resp.Code = http.StatusBadRequest
        resp.Err = "illegal method"
        httputil.WriteJSONResponse(writer, resp)
        return
    }
}


// swagger:route DELETE /api/v1/fs/files deleteFiles rmCmd
//
// Handler for delete files or directory.
//
// the response body can be fetched from response body
//
// Responses:
// 		        200: response
// 		        400: response
// 		        401: response
func (self *FSHandler) DeleteFilesHandler(writer http.ResponseWriter, request *http.Request) {
    resp := httputil.Response{}

    // read user info from token
    user, err := self.getUserInfo(writer, request)
    if err != nil {
        resp.Code = http.StatusUnauthorized
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    body, err := httputil.ReadBody(request)
    if err != nil {
        resp.Code = http.StatusBadRequest
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
    }

    rmHandler(writer, user, body)

}

func (self *FSHandler)getChunkMetaHandler(writer http.ResponseWriter, request *http.Request) {
    log.Debugf("begin proc getChunkMeta\n")
    resp := httputil.Response{}

    user, err := self.getUserInfo(writer, request)
    if err != nil {
        resp.Code = http.StatusUnauthorized
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    rawQuery := request.URL.RawQuery

    config := syscfg.GetSysConfig()
    root := config.Get(common.FSBasePath).(string)

    if !user.Sysadmin {
        root = filepath.Join(root, user.Username)
    }

    cmd, err := fscmd.NewChunkMetaCmdFromURLParam(root, rawQuery, user)

    if err != nil {
        resp.Err = err.Error()
        resp.Code = http.StatusInternalServerError
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    result, err := cmd.Run()
    if err != nil {
        resp.Code = http.StatusInternalServerError
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    resp.Code = http.StatusOK
    resp.Results = result
    httputil.WriteJSONResponse(writer, resp)
    log.Debugf("end proc getChunkMeta\n")
}

// GetChunkMetaHandler processes GET ChunkMeta  request.
func (self *FSHandler)GetChunkMetaHandler(writer http.ResponseWriter, request *http.Request) {
    method := request.URL.Query().Get("method")
    resp := httputil.Response{}
    switch method {
    case "GetChunkMeta":
        self.getChunkMetaHandler(writer, request)
    default:
        resp.Code = http.StatusMethodNotAllowed
        httputil.WriteJSONResponse(writer, resp)
    }
}

// GetChunkHandler processes GET Chunk  request.
func (self *FSHandler)GetChunkHandler(writer http.ResponseWriter, request *http.Request) {
    log.Debugf("begin proc GetChunkHandler")
    resp := httputil.Response{}

    user, err := self.getUserInfo(writer, request)
    if err != nil {
        resp.Code = http.StatusUnauthorized
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    config := syscfg.GetSysConfig()
    root := config.Get(common.FSBasePath).(string)

    if !user.Sysadmin {
        root = filepath.Join(root, user.Username)
    }

    chunk, err := fscmd.ParseChunk(request.URL.RawQuery)
    if err != nil {
        log.Errorf("parse chunk Error: %s", err.Error())
        resp.Code = http.StatusBadRequest
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }


    mWriter := multipart.NewWriter(writer)
    mWriter.SetBoundary(fscmd.DefaultMultiPartBoundary)

    fileName := chunk.ToURLParam().Encode()

    part, err := mWriter.CreateFormFile("chunk", fileName)
    if err != nil {
        log.Error(err)
        resp.Code = http.StatusInternalServerError
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // before load chunk data. we will reformat file path
    // according to the root
    chunk.Path = filepath.Join(root, chunk.Path)
    if err = chunk.LoadChunkData(part); err != nil {
        log.Error(err)
        resp.Code = http.StatusInternalServerError
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    defer mWriter.Close()

    log.Debug("end proc GetChunkHandler")
}

// PostChunkHandler processes POST Chunk request.
func (self *FSHandler)PostChunkHandler(writer http.ResponseWriter, request *http.Request) {
    log.Debugf("begin proc PostChunksHandler\n")

    resp := httputil.Response{}

    user, err := self.getUserInfo(writer, request)
    if err != nil {
        resp.Code = http.StatusUnauthorized
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    config := syscfg.GetSysConfig()
    root := config.Get(common.FSBasePath).(string)

    if !user.Sysadmin {
        root = filepath.Join(root, user.Username)
    }

    defer request.Body.Close()

    //partReader, err := request.MultipartReader()
    partReader := multipart.NewReader(request.Body, fscmd.DefaultMultiPartBoundary)
    if err != nil {
        resp.Code = http.StatusBadRequest
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    for {
        part, err := partReader.NextPart()
        if err == io.EOF {
            break
        }

        // check part read error
        if err != nil {
            log.Errorf("part reader Error: %#v", err)
            resp.Err = err.Error()
            resp.Code = http.StatusInternalServerError
            httputil.WriteJSONResponse(writer, resp)
            return
        }

        if part.FormName() != "chunk" {
            continue
        }

        chunk, err := fscmd.ParseChunk(part.FileName())
        if err != nil {
            resp.Err = err.Error()
            resp.Code = http.StatusInternalServerError
            httputil.WriteJSONResponse(writer, resp)
            return
        }


        // before load chunk data. we will reformat file path
        // according to the root
        chunk.Path = filepath.Join(root, chunk.Path)

        log.Debugf("recv chunk: %#v\n", chunk)

        if err := chunk.SaveChunkData(part); err != nil {
            resp.Err = err.Error()
            resp.Code = http.StatusInternalServerError
            httputil.WriteJSONResponse(writer, resp)
            return
        }

    }

    resp.Code = http.StatusOK
    httputil.WriteJSONResponse(writer, resp)

    log.Debugf("end proc PostChunksHandler\n")
}
