package api

import (
    "github.com/deepinsight/mlcloud/src/pkg/cluster"
    "github.com/deepinsight/mlcloud/src/pkg/auth"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/gorilla/mux"
    "net/http"
)

// APIHandler is a representation of API handler. Structure contains client, Heapster client and client configuration.
type BaseHandler struct {
    clientset cluster.ClientSet
    tManager auth.TokenManager
}

/**
func (self *BaseHandler) Register(ws *restful.WebService)  {
    ws.Route(
        ws.GET("/version").
            To(self.handleVersion))
}
**/

func (self *BaseHandler) Register(router *mux.Router)  {
    bs := router.PathPrefix("/api/v1").Subrouter()
    bs.HandleFunc("/health", self.handleHealth).
        Methods("GET")
}

// swagger:route GET /api/v1/health health healthCheck
//
// Handler for health check.
//
// the response body can be fetched from response body
//
// Responses:
// 		        200: "health OK"
func (self *BaseHandler) handleHealth(writer http.ResponseWriter, request *http.Request) {
    log.Debug("request host: ", request.Host)
    writer.Write([]byte("health OK"))
}


func (self *BaseHandler) getUserInfo(writer http.ResponseWriter, request *http.Request) (*models.User, error) {
    tokenStr := request.Header.Get("token")
    user, err := self.tManager.GetUserInfo(tokenStr)
    if err != nil {
        return nil, err
    }

    return user, nil
}

// NewBaseHandler created base handler instance.
func NewBaseHandler(clientset cluster.ClientSet, tokenManager auth.TokenManager) *BaseHandler {
    return &BaseHandler{
            clientset: clientset,
            tManager: tokenManager,
        }
}
