package api

import (
    "net/http"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/deepinsight/mlcloud/src/pkg/dao"
    "github.com/deepinsight/mlcloud/src/pkg/crm"
    "fmt"
    "github.com/gorilla/mux"
    httputil "github.com/deepinsight/mlcloud/src/utils/http"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/deepinsight/mlcloud/src/pkg/fs"
)

type UserHandler struct {
    *BaseHandler
}

func (self *UserHandler) Register(router *mux.Router)  {
    bs := router.PathPrefix("/api/v1").Subrouter()

    bs.HandleFunc("/users", self.ListUser).
        Methods("GET")

    bs.HandleFunc("/users/{username}", self.DeleteUser).
        Methods("DELETE")

    bs.HandleFunc("/users/{username}/changepassword", self.ChangeUserPassword).
        Methods("POST")
}

// List All the users
func (self *BaseHandler) ListUser(writer http.ResponseWriter, request *http.Request) {
    resp := httputil.Response{}
    tokenStr := request.Header.Get("token")
    if len(tokenStr) == 0 {
        resp.Code = http.StatusUnauthorized
        resp.Err = "no authentication info provided"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    user, err := self.getUserInfo(writer, request)
    if err != nil {
        resp.Code = http.StatusUnauthorized
        resp.Err = "authentication failed"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    if !user.Sysadmin {
        resp.Code = http.StatusForbidden
        resp.Err = "user list need administrator privilidge"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // return the total users.
    // TODO: get all the user info from database
    // should be use Pagination
    users, err := dao.ListUsers()
    if err != nil {
        log.Errorf("cannot list users: %v", err)
        resp.Code = http.StatusInternalServerError
        resp.Err = "cannot list users"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    resp.Code = http.StatusOK
    resp.Results = users
    httputil.WriteJSONResponse(writer, resp)
}

// delete user API
func (self *UserHandler) DeleteUser(writer http.ResponseWriter, request *http.Request) {
    resp := httputil.Response{}
    tokenStr := request.Header.Get("token")
    if len(tokenStr) == 0 {
        resp.Code = http.StatusUnauthorized
        resp.Err = "no authentication info provided"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    user, err := self.getUserInfo(writer, request)
    if err != nil {
        resp.Code = http.StatusUnauthorized
        resp.Err = "authentication failed"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    if !user.Sysadmin {
        resp.Code = http.StatusForbidden
        resp.Err = "user list need administrator privilidge"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // get the username to be deleted
    tUserName := mux.Vars(request)["username"]
    if tUserName == "admin" || len(tUserName) == 0 {
        resp.Code = http.StatusForbidden
        resp.Err = "admin user don't allow to be deleted"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // target user
    tUser := &models.User{
        Username: tUserName,
    }

    // check if user exist
    exist, _ := dao.CheckUserExist(tUser)
    if !exist {
        resp.Code = http.StatusBadRequest
        resp.Err = fmt.Sprintf("user: %s doesn't exist.", tUserName)
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // delete user kubernetes namespace resources
    // such as jobs, deployments, services, pods, roles, rolebindings, namespaces
    // deleteUserClusterResources(user)
    if err := crm.DeleteUserNamespace(tUser); err != nil {
        errorInfo := fmt.Sprintf("failed to delete user namespace: %s", tUser.Username)
        resp.Code = http.StatusInternalServerError
        resp.Err = errorInfo
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    if err := crm.DeleteUserPV(tUser); err != nil {
        resp.Code = http.StatusInternalServerError
        resp.Err = "failed to delete user's pv: " + tUser.Username
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // delete user's file system
    if err := fs.DeleteBaseDirForUser(tUser); err != nil {
        resp.Code = http.StatusInternalServerError
        resp.Err = "failed to delete user base dir"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // delete user record from databases
    if err := dao.DeleteUser(tUser); err != nil {
        errorInfo := fmt.Sprintf("failed to delete user: %s", tUser.Username)
        resp.Code = http.StatusInternalServerError
        resp.Err = errorInfo
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    resp.Code = http.StatusOK
    resp.Results = fmt.Sprintf("user: %s was deleted", tUserName)
    httputil.WriteJSONResponse(writer, resp)

}

// change the password for user
func (self *UserHandler) ChangeUserPassword(writer http.ResponseWriter, request *http.Request) {

}


// NewAuthHandler created AuthHandler instance.
func NewUserHandler(baseHandler *BaseHandler) *UserHandler {
    return &UserHandler {
        BaseHandler: baseHandler,
    }
}