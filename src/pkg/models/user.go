package models

import (
    "time"
    "fmt"
    "strings"
    "github.com/astaxie/beego/orm"
)

func init() {
    orm.RegisterModel(new(User))
}

// user model
type User struct {
    // user ID
    Id   int    `orm:"pk;auto;column(id)" json:"id"`

    // user Name
    Username string `orm:"column(username)" json:"username"`

    // user password
    Password string `orm:"column(password)" json:"password"`

    // the key to encrypt the user password
    Salt string `orm:"column(salt)" json:"-"`

    // rolename will not be
    Sysadmin bool `orm:"column(sys_admin)" json:"sys_admin"`

    // user creation time
    CreationTime time.Time `orm:"column(creation_time)" json:"creation_time"`

    // user update time
    UpdateTime   time.Time `orm:"column(update_time)" json:"update_time"`
}

// validate only validate when user register
func Validate(user *User) error {
    if isIllegalLength(user.Username, 1, 20) {
        return fmt.Errorf("username: %s with illegal length", user.Username)
    }
    if isContainIllegalChar(user.Username, []string{",", "~", "#", "$", "%"}) {
        return fmt.Errorf("username: %s contains illegal characters", user.Username)
    }
    if isIllegalLength(user.Password, 8, 20) {
        return fmt.Errorf("user password with illegal length")
    }

    return nil
}

func isIllegalLength(s string, min int, max int) bool {
    if min == -1 {
        return (len(s) > max)
    }
    if max == -1 {
        return (len(s) <= min)
    }
    return (len(s) < min || len(s) > max)
}

func isContainIllegalChar(s string, illegalChar []string) bool {
    for _, c := range illegalChar {
        if strings.Index(s, c) >= 0 {
            return true
        }
    }
    return false
}
