package dao

import (
    "testing"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "os"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "time"
    "k8s.io/apimachinery/pkg/util/json"
)

var defaultRegistered = false

func initDatabaseForTest(db *models.Database) {
    database, err := getDatabase(db)
    if err != nil {
        panic(err)
    }

    log.Infof("initializing database: %s", database.String())

    alias := database.Name()
    if !defaultRegistered {
        defaultRegistered = true
        alias = "default"
    }
    if err := database.Register(alias); err != nil {
        panic(err)
    }

    if alias != "default" {
        if err = globalOrm.Using(alias); err != nil {
            log.Fatalf("failed to create new orm: %v", err)
        }
    }

    // if we will create database and table in mysql docker container
    // then there is no need to run syncdb
    // sync db
    //orm.RunSyncdb("default", false, true)
}

func TestMain(m *testing.M) {
    database := &models.Database{
        Type: "mysql",
        MySQL: &models.MySQL{
            Host:     "127.0.0.1",
            Port:     3306,
            Username: "mlcloud",
            Password: "mlcloud",
            Database: "mlcloud",
        },
    }

    initDatabaseForTest(database)

    result := m.Run()
    if result != 0 {
        os.Exit(result)
    }
}

/*const USERNAME string  = "test"
const PASSWORD string = "test"

func TestDeleteUser(t *testing.T)  {
    user := models.User{
        Username: USERNAME,
    }

    err := DeleteUser(&user)
    if err != nil {
        t.Log(err)
    }
}


func TestAddUser(t *testing.T) {
    user := models.User{
        Username: USERNAME,
        Password: USERNAME,
        Sysadmin: false,
        CreationTime: time.Now(),
        UpdateTime: time.Now(),
    }
    id , err := AddUser( &user )

    t.Logf("user id: %d", id)
    if err != nil {
        log.Fatal(err)
    }


}*/

const REGISTER_USER string = "register_user"
const REGISTER_PASSWORD string = "register_password"
func TestRegisterUser(t *testing.T) {
    user := models.User{
        Username: REGISTER_USER,
        Password: REGISTER_PASSWORD,
        Sysadmin: false,
        CreationTime: time.Now(),
        UpdateTime: time.Now(),
    }

    exist, err := CheckUserExist(&user)
    if err != nil {
        t.Error(err)
    }

    if exist {
        return
    }

    uid, err := Register(&user)
    if err != nil {
        t.Error(err)
    }

    t.Logf("register user successfully, user id: %d", uid)
}

func TestCheckUser(t *testing.T) {
    user := models.User{
        Username: REGISTER_USER,
    }

    exist, err := CheckUserExist(&user)
    if err != nil {
        t.Error(err)
    }

    if !exist {
        t.Errorf("user %s has been registered. but it doesn't exist", user.Username)
    }
}

func TestCheckUserPassword(t *testing.T) {
    user := models.User{
        Username: REGISTER_USER,
        Password: REGISTER_PASSWORD,
    }

    isOk, err := CheckUserPassword(&user)
    if err != nil {
        t.Error(err)
    }

    if !isOk {
        t.Error("user password is not correct")
    }
}

func TestMarshalUser(t *testing.T) {
    user := models.User{
        UserID: 1,
        Username: REGISTER_USER,
        Password: REGISTER_PASSWORD,
        Sysadmin: false,
        Salt: "MLCloud",
        CreationTime: time.Now(),
        UpdateTime: time.Now(),
    }

    userBytes, err := json.Marshal(&user)
    if err != nil {
        t.Error("failed to marshal user")
    }

    targetUser := models.User{}
    if err := json.Unmarshal(userBytes, &targetUser); err != nil {
        t.Error(err)
    }

    if len(targetUser.Salt) != 0 {
        t.Errorf("the user salt should be null after marshal. but actually it is: %s", targetUser.Salt)
    }
}