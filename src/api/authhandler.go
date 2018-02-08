package api

import (
    "net/http"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/gorilla/mux"
    "github.com/deepinsight/mlcloud/src/pkg/auth"
    "github.com/deepinsight/mlcloud/src/pkg/dao"
    "fmt"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/deepinsight/mlcloud/src/pkg/crm"
    "github.com/deepinsight/mlcloud/src/pkg/fs"
    httputil "github.com/deepinsight/mlcloud/src/utils/http"
)

type AuthHandler struct {
    *BaseHandler
}

// register handlers for login and sign up API
func (self *AuthHandler) Register(router *mux.Router) {
    bs := router.PathPrefix("/api/v1").Subrouter()
    bs.HandleFunc("/login", self.handleLogin).
        Methods("POST").
        Headers("Content-Type", "application/json")

    bs.HandleFunc("/signup", self.handleSignUp).
        Methods("POST").
        Headers("Content-Type", "application/json")
}

// swagger:route POST /api/v1/login login userSpec
//
// Handler for user login.
//
// the response code can be fetched from response body
//
// Responses:
// 		        200: response
// 		        400: response
// 		        401: response
// 		        500: response
func (self *BaseHandler) handleLogin(writer http.ResponseWriter, request *http.Request) {
    resp := httputil.Response{}
    loginSpec := auth.LoginSpec{}
    if err := httputil.ReadEntity(request, &loginSpec); err != nil {
        resp.Code = http.StatusBadRequest
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    log.Debugf("login spec: %+v", loginSpec)

    user := models.User{Username: loginSpec.Username, Password: loginSpec.Password}
    // check user exist
    exist, _ := dao.CheckUserExist(&user)
    if !exist {
        errorInfo := fmt.Sprintf("user: %s doesn't exist. please register first", user.Username)
        resp.Code = http.StatusUnauthorized
        resp.Err = errorInfo
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // check the password
    isOk, _ := dao.CheckUserPassword(&user)
    if !isOk {
        resp.Code = http.StatusUnauthorized
        resp.Err = "user password is incorrect"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // get the full information of the user from DAO（Data Access Object）
    currentUser, _ := dao.GetUser(&user)
    // invoke the token manager to generate the token
    loginResponse, err := self.tManager.Generate(currentUser)
    if err != nil {
        resp.Code = http.StatusInternalServerError
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // all things goes well
    resp.Code = http.StatusOK
    resp.Results = loginResponse
    httputil.WriteJSONResponse(writer, resp )
}

// swagger:route POST /api/v1/signup signup userSpec
//
// Handler for user signup.
//
// the response code can be fetched from response body
//
// Responses:
// 		        200: response
// 		        400: response
// 		        403: response
// 		        405: response
// 		        500: response
func (self *BaseHandler) handleSignUp(writer http.ResponseWriter, request *http.Request) {
    resp := httputil.Response{}

    // read user spec
    user := new(models.User)
    if err := httputil.ReadEntity(request, &user); err != nil {
        resp.Code = http.StatusBadRequest
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // validate user name and password
    if err := models.Validate(user); err != nil {
        resp.Code = http.StatusNotAcceptable
        resp.Err = err.Error()
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // check if user exists
    exist, _ := dao.CheckUserExist(user)
    if exist {
        errorInfo := fmt.Sprintf("user: %s exist. please change another name", user.Username)
        resp.Code = http.StatusForbidden
        resp.Err = errorInfo
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // register user
    _, err := dao.Register(user)
    if err != nil {
        errorInfo := fmt.Sprintf("failed to register user: %s", user.Username)
        resp.Code = http.StatusForbidden
        resp.Err = errorInfo
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // create related resources in kubernetes cluster
    // including namespace, rolebindings
    // so that the user could access the cluster.
    if err := crm.CreateUserNamespace(user); err != nil {
        errorInfo := fmt.Sprintf("failed to create user namespace: %s", user.Username)
        resp.Code = http.StatusInternalServerError
        resp.Err = errorInfo
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // create rolebindings
    // for now. we don't create user rolebindings on that namespace.
    // we just use admin user
/*    if err := crm.CreateUserRolebinding(user); err != nil {
        errorInfo := fmt.Sprintf("failed to create rolebinding for user: %s", user.Username)
        resp.Code = http.StatusInternalServerError
        resp.Err = errorInfo
        httputil.WriteJSONResponse(writer, resp)
        return
    }*/

    // create PV / PVC for user
    if err := crm.CreateUserNamespaceCephFSSecret(user); err != nil {
        resp.Code = http.StatusInternalServerError
        resp.Err = "failed to create ceph secret"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    if err := crm.CreateUserPV(user); err != nil {
        resp.Code = http.StatusInternalServerError
        resp.Err = "failed to create ceph pv"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    if err := crm.CreateUserPVC(user); err != nil {
        resp.Code = http.StatusInternalServerError
        resp.Err = "failed to create ceph pvc"
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // create base fs directory for user
    if err := fs.CreateBaseDirForUser(user); err != nil {
        errorInfo := fmt.Sprintf("failed to create base dir for user: %s", user.Username)
        resp.Code = http.StatusInternalServerError
        resp.Err = errorInfo
        httputil.WriteJSONResponse(writer, resp)
        return
    }

    // every thing goes well
    resp.Code = http.StatusOK
    resp.Results = fmt.Sprintf("user: %s register successfully!!", user.Username)
    httputil.WriteJSONResponse(writer, resp)

}

// NewAuthHandler created AuthHandler instance.
func NewAuthHandler(baseHandler *BaseHandler) *AuthHandler {
    return &AuthHandler{
        BaseHandler: baseHandler,
    }
}