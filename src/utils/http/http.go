package http

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "errors"
    "github.com/deepinsight/mlcloud/src/utils/log"
)


// An HTTP response model
//
// swagger:model response
type Response struct {
    // the error message
    //
    // required: true
    Err     string      `json:"err"`

    // the error code
    //
    // required: true
    Code    int         `json:"code"`

    // the response body
    //
    // required: true
    Results interface{} `json:"results"`
}

func WriteJSONResponse(writer http.ResponseWriter, resp Response) {
    writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
    writer.WriteHeader(resp.Code)

    if err := json.NewEncoder(writer).Encode(resp); err != nil {
        log.Error(err)
    }
}

func ReadBody(request *http.Request) ([]byte, error) {
    body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        return nil, err
    }

    if err := request.Body.Close(); err != nil {
        return nil, err
    }

    return body, nil
}


func GetKeyFromBody(body []byte, key string) (interface{}, error) {
    maps := map[string]interface{}{}
    if err := json.Unmarshal(body, &maps); err != nil {
        return nil, err
    }

    if val, ok := maps[key]; ok {
        return val, nil
    }

    return nil, errors.New(fmt.Sprintf("no key: %s found from input", key))
}


func ReadEntity(request *http.Request, entity interface{}) error {
    bytes, err := ioutil.ReadAll(request.Body)
    if err != nil {
        log.Errorf("read body error: %#v", err)
        return err
    }

    defer request.Body.Close()

    return json.Unmarshal(bytes, entity)
}

func ReadForm(request *http.Request, entity interface{}) error {
    if err := request.ParseForm(); err != nil {
        return err
    }

    values := request.Form
    log.Infof("form values: %+v", values)
    return nil
}

