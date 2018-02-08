// go-swagger document for mlcloud.
//
// The purpose of this application is to provide a machine learning platform to adaptor
// Tensorflow / MXNet / Caffe based on kubernetes container cloud.
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /
//     Version: 1.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
    "net/http"
    "fmt"
    "github.com/deepinsight/mlcloud/src/api"
    "github.com/deepinsight/mlcloud/src/pkg/auth/jwt"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/deepinsight/mlcloud/src/pkg/dao"
    "github.com/deepinsight/mlcloud/src/pkg/common/syscfg"
    "github.com/deepinsight/mlcloud/src/pkg/cluster"
    "os"
    "flag"
)

func main() {
    port := flag.Int("port", 9090, "port of server")
    ip := flag.String("ip", "0.0.0.0", "ip of server")

    // for glog
    flag.Parse()

    // init config. panic if error
    if err := syscfg.Init(); err != nil {
        log.Fatal(err)
    }

    // init database. panic if error
    if err := dao.Init(); err != nil {
        log.Fatal(err)
    }

    clientset, err := cluster.GetClientSet()
    if err != nil {
        log.Error(err)
        os.Exit(1)
    }
    tManager := jwt.NewJWTTokenManager()

    // add swagger comments
    apiHandler, err := api.NewRouter(*clientset, tManager)
    if err != nil {
        panic(err)
    }

    //http.HandleFunc("/version", func(res http.ResponseWriter, req *http.Request) {res.Write([]byte("version 1.0"))})
    //http.Handle("/api/", apiHandler)
    // Listen for http and https
    addr := fmt.Sprintf("%s:%d", *ip, *port)
    log.Infof("Serving insecurely on HTTP server: %s", addr)
    go func() {
        log.Fatal(http.ListenAndServe(addr, apiHandler ))
    }()

    select {}

}