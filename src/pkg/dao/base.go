package dao

import (
    "fmt"
    "strconv"
    "strings"
    "sync"
    "github.com/astaxie/beego/orm"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/deepinsight/mlcloud/src/pkg/common/syscfg"
    "github.com/deepinsight/mlcloud/src/pkg/common"
    "github.com/deepinsight/mlcloud/src/utils"
)

// Database is an interface of different databases
type Database interface {
    // Name returns the name of database
    Name() string
    // String returns the details of database
    String() string
    // Register registers the database which will be used
    Register(alias ...string) error
}

// Init DataBase From config
func Init() error {
    // get system config
    config := syscfg.GetSysConfig()

    port, err := utils.ParseStringToInt(config.Get(common.MySQLPort).(string))
    if err != nil {
        return err
    }

    database := &models.Database{
        Type: "mysql",
        MySQL: &models.MySQL{
            Host:     config.Get(common.MySQLHost).(string),
            Port:     port,
            Username: config.Get(common.MySQLUsername).(string),
            Password: config.Get(common.MySQLPassword).(string),
            Database: config.Get(common.MySQLDatabase).(string),
        },
    }

    return initDatabase(database)
}

// InitDatabase initializes the database
func initDatabase(database *models.Database) error {
    db, err := getDatabase(database)
    if err != nil {
        return err
    }

    log.Infof("initializing database: %s", db.String())
    if err := db.Register(); err != nil {
        return err
    }

    // sync the default db. basically it will create or update table schema
    //orm.RunSyncdb("default", false, true)

    // register admin user
    log.Info("Register system admin user")
    if err := RegisterSysAdmin(); err != nil {
        return err
    }

    return nil
}


// TODO load the user name and password from system config
const ADMIN_USER_NAME = "admin"
const ADMIN_USER_PASSWORD = "mlcloud"

func RegisterSysAdmin() error {
    user := &models.User {
        Username: ADMIN_USER_NAME,
        Password: ADMIN_USER_PASSWORD,
        Sysadmin: true,
    }

    exist, err := CheckUserExist(user)
    if err != nil {
        return err
    }

    if !exist {
        _, err := Register(user)
        return err
    }

    return nil
}

func getDatabase(database *models.Database) (db Database, err error) {
    switch database.Type {
    case "", "mysql":
        db = NewMySQL(database.MySQL.Host,
            strconv.Itoa(database.MySQL.Port),
            database.MySQL.Username,
            database.MySQL.Password,
            database.MySQL.Database)
    default:
        err = fmt.Errorf("invalid database: %s", database.Type)
    }
    return
}

var globalOrm orm.Ormer
var once sync.Once

// GetOrmer :set ormer singleton
func GetOrmer() orm.Ormer {
    once.Do(func() {
        globalOrm = orm.NewOrm()
    })
    return globalOrm
}

// ClearTable is the shortcut for test cases, it should be called only in test cases.
func ClearTable(table string) error {
    o := GetOrmer()
    sql := fmt.Sprintf("delete from %s where 1=1", table)
    _, err := o.Raw(sql).Exec()
    return err
}

func paginateForRawSQL(sql string, limit, offset int64) string {
    return fmt.Sprintf("%s limit %d offset %d", sql, limit, offset)
}

func escape(str string) string {
    str = strings.Replace(str, `%`, `\%`, -1)
    str = strings.Replace(str, `_`, `\_`, -1)
    return str
}

