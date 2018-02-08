package auth

import "github.com/deepinsight/mlcloud/src/pkg/models"


// AuthManager is used for user authentication management.
type AuthManager interface {
	// Login authenticates user based on provided LoginSpec and returns LoginResponse. LoginResponse contains
	// generated token and list of non-critical errors such as 'Failed authentication'.
	Login(*LoginSpec) (*LoginResponse, error)
}

// TokenManager is responsible for generating and decrypting tokens used for authorization. Authorization is handled
// by K8S apiserver. Token contains AuthInfo structure used to create K8S api client.
type TokenManager interface {
	// Generate secure token based on AuthInfo structure and save it it tokens' payload.
	Generate(*models.User) (*LoginResponse, error)
	// Decrypt generated token and return AuthInfo structure that will be used for K8S api client creation.
	Validate(string) (bool, error)

    GetUserInfo(tokenStr string) (*models.User, error)
}

// user login spec.
//
// swagger:parameters userSpec
type LoginSpec struct {
	// Username is the username for basic authentication to the ML Cloud.
    //
    // in: body
    // required: true
	Username string `json:"username"`

	// Password is the password for basic authentication to the ML Cloud.
    //
    // in: body
    // required: true
    Password string `json:"password"`
}


// LoginResponse is returned from our backend as a response for login request. It contains generated JWEToken and a list
// of non-critical errors such as 'Failed authentication'.
type LoginResponse struct {
	// JWEToken is a token generated during login request that contains AuthInfo data in the payload.
	Token string `json:"token"`
    //Cookie http.Cookie `json:"cookie"`
	// Errors are a list of non-critical errors that happened during login request.
	//Errors []error `json:"errors"`
}
