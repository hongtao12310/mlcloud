package api

import (
    "github.com/deepinsight/mlcloud/src/pkg/cluster"
    "net/http"
    "github.com/deepinsight/mlcloud/src/pkg/auth"
    "github.com/gorilla/mux"
)

// CreateHTTPAPIHandler creates a new HTTP handler that handles all requests to the API of the backend.
func NewRouter(clientset cluster.ClientSet, tManager auth.TokenManager) (http.Handler, error) {
    router := mux.NewRouter()

    // register base handler
    baseHandler := NewBaseHandler(clientset, tManager)
    baseHandler.Register(router)

    // register AUTH handler
    authHandler := NewAuthHandler(baseHandler)
    authHandler.Register(router)

    // register user handler
    userHandler := NewUserHandler(baseHandler)
    userHandler.Register(router)

    // register FS handler
    fsHandler := NewFSHandler(baseHandler)
    fsHandler.Register(router)

    // register job handler
    jobHandler := NewJobHandler(baseHandler)
    jobHandler.Register(router)

    return router, nil
}