package jwt

import (
    "github.com/deepinsight/mlcloud/src/pkg/auth"
    "github.com/dgrijalva/jwt-go"
    "time"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "fmt"
)

const (
    SECRET = "Wanda Machine Learning Cloud Platform"
)

type JWTTokenManager struct {
    Secret string `json:"secret"`
}

type UserClaims struct {
    User models.User `json:"user"`
    jwt.StandardClaims
}

// generate jwt token "github.com/dgrijalva/jwt-go"
func (self *JWTTokenManager) Generate(user *models.User) (*auth.LoginResponse, error) {
    expireToken := time.Now().Add(time.Hour * 24).Unix()
    //expireCookie := time.Now().Add(time.Hour * 24)
    iat := time.Now().Unix()

    claims := UserClaims {
        *user,
        jwt.StandardClaims {
            ExpiresAt: expireToken,
            Issuer: "wanda.com",
            IssuedAt: iat,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    signedToken, _ := token.SignedString([]byte(SECRET))

    //cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}

    return &auth.LoginResponse{Token: signedToken }, nil
}

// check if the token has the permission to access our resources
func (self *JWTTokenManager) Validate(tokenStr string) (bool, error) {
    _, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error){
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
            return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
        }
        return []byte(SECRET), nil
    })

    if err != nil {
        return false, err
    }

    return true, nil
}

// get user info by token
func (self *JWTTokenManager) GetUserInfo(tokenStr string) (*models.User, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error){
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
            return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
        }
        return []byte(SECRET), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
        return &claims.User, nil
    } else {
        return nil, jwt.NewValidationError("claims invalid", jwt.ValidationErrorClaimsInvalid)
    }
}

func NewJWTTokenManager() auth.TokenManager {
    return &JWTTokenManager{Secret: SECRET}
}