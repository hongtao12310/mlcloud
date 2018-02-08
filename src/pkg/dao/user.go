package dao

import (
    "fmt"
    "time"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/deepinsight/mlcloud/src/utils"
    "github.com/deepinsight/mlcloud/src/pkg/auth"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/deepinsight/mlcloud/src/utils/errors"
)

func AddUser(user *models.User) (id int64, err error) {
    o := GetOrmer()
    id, err = o.Insert(user)
    return
}

// check if the user exists
func CheckUserExist(user *models.User) (bool, error) {
    qs := GetOrmer().QueryTable(&models.User{})
    if user.Id > 0 {
        qs = qs.Filter("id", user.Id)
    }

    if len(user.Username) > 0 {
        qs = qs.Filter("username", user.Username)
    }

    count, err := qs.Count()

    if err != nil {
        return false, err
    }

    if count == 0 {
        return false, nil
    }

    return true, nil
}

// Register is used for user to register, the password is encrypted before the record is inserted into database.
func Register(user *models.User) (int64, error) {
    o := GetOrmer()
    p, err := o.Raw("insert into user (username, password, salt, sys_admin, creation_time, update_time) values (?, ?, ?, ?, ?, ?)").Prepare()
    if err != nil {
        return 0, err
    }
    defer p.Close()

    salt := utils.GenerateRandomString()

    now := time.Now()
    r, err := p.Exec(user.Username, utils.Encrypt(user.Password, salt), salt, user.Sysadmin, now, now)

    if err != nil {
        return 0, err
    }
    userID, err := r.LastInsertId()
    if err != nil {
        return 0, err
    }

    return userID, nil
}

// CheckUserPassword checks whether the password is correct.
func CheckUserPassword(user *models.User) (bool, error) {
    currentUser, err := GetUser(user)
    if err != nil {
        return false, err
    }

    if currentUser == nil {
        return false, errors.New(USER_NOT_FOUND, "user: %s not found", currentUser.Username)
    }

    desiredPassword := utils.Encrypt(user.Password, currentUser.Salt)
    if currentUser.Password != desiredPassword {
        return false, nil
    }

    return true, nil
}

// Get user according to user name or user id
func GetUser(user *models.User) (*models.User, error) {

    o := GetOrmer()

    sql := `select id, username, password, salt,
		sys_admin, creation_time, update_time
		from user u where 1=1 `
    queryParam := make([]interface{}, 1)
    if user.Id != 0 {
        sql += ` and id = ? `
        queryParam = append(queryParam, user.Id)
    }

    if user.Username != "" {
        sql += ` and username = ? `
        queryParam = append(queryParam, user.Username)
    }

    var u []models.User
    n, err := o.Raw(sql, queryParam).QueryRows(&u)

    if err != nil {
        return nil, err
    }
    if n == 0 {
        return nil, nil
    }

    if n > 1 {
        return nil, fmt.Errorf("got more than one user when executing: %s param: %v", sql, queryParam)
    }

    return &u[0], nil
}

// LoginByDb is used for user to login with database auth mode.
func LoginByDb(auth *auth.LoginSpec) (*models.User, error) {
    o := GetOrmer()

    var users []models.User
    n, err := o.Raw(`select * from user where username = ?`,
        auth.Username).QueryRows(&users)
    if err != nil {
        return nil, err
    }
    if n == 0 {
        return nil, nil
    }

    user := users[0]

    if user.Password != utils.Encrypt(auth.Password, user.Salt) {
        return nil, nil
    }

    user.Password = "" //do not return the password

    return &user, nil
}

// GetTotalOfUsers ...
func GetTotalOfUsers() (int64, error) {
    return GetOrmer().QueryTable(&models.User{}).Count()
}

// ListUsers lists all users according to different conditions.
func ListUsers() ([]models.User, error) {
    users := []models.User{}
    _, err := GetOrmer().QueryTable(&models.User{}).
        Limit(-1).OrderBy("username").All(&users)
    return users, err
}


// DeleteUser ...
func DeleteUser(user *models.User) error {
    o := GetOrmer()

    u, err := GetUser(user)
    if u == nil {
        log.Debug("delete user failed. no user exist")
        return nil
    }

    _, err = o.Delete(u)
    log.Debugf("user #%d deleted", u.Id)
    return err
}


