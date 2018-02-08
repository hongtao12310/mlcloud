package syscfg

import (
    "os"
    "github.com/deepinsight/mlcloud/src/pkg/common"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/deepinsight/mlcloud/src/pkg/common/syscfg/store"
    "github.com/deepinsight/mlcloud/src/pkg/common/syscfg/store/json"
)

var (
    defaultJSONCfgStorePath = os.Getenv("HOME") + "/mlcloud/config/config.json"
    defaultKubeCertDIR = os.Getenv("HOME") + "/mlcloud/certs"
)

var (
    // all configurations need read from environment variables
    allEnvs = map[string]string {
        common.MySQLHost:    "MYSQL_HOST",
        common.MySQLPort:    "MYSQL_PORT",
        common.MySQLUsername: "MYSQL_USER",
        common.MySQLPassword: "MYSQL_PASSWORD",
        common.MySQLDatabase: "MYSQL_DATABASE",
        common.KubeAPIServer: "KUBE_APISERVER",
        common.KubeClusterName: "KUBE_CLUSTER_NAME",
        common.KubeCertDir: "KUBE_CERT_DIR",
        common.KubeInCluster: "KUBE_IN_CLUSTER",
        common.AdminKubeConfigPath: "ADMIN_KUBE_CONFIG",
        common.FSBasePath: "FS_BASE_PATH",
    }
)

type parser struct {
    // the name of env
    env   string
    // parse the value of env, e.g. parse string to int or
    // parse string to bool
    parse func(string) (interface{}, error)
}


type SysConfig struct {
    Contents    map[string]interface{}

    // A storage driver that configurations
    // can be read from and wrote to
    StoreDriver store.Driver
}

func (self *SysConfig) Init() {
    self.Contents = make(map[string]interface{})
    for k, _ := range allEnvs {
        self.Contents[k] = ""
    }
}

func (self *SysConfig) Load() (err error) {
    contents, err := self.StoreDriver.Read()
    if contents != nil {
        self.Contents = contents
    }

    // override the config from environments
    self.LoadFromEnv()

    // give the default value for cert dir
    if _, ok := self.Contents[common.KubeCertDir]; !ok {
        self.Contents[common.KubeCertDir] = defaultKubeCertDIR
    }

    return nil
}

func (self *SysConfig) LoadFromEnv() {
    for k, v := range allEnvs {
        val := os.Getenv(v)
        if len(val) > 0 {
            self.Contents[k] = val
        }
    }
}

func (self *SysConfig) Get(key string) interface{} {
    if val, ok := self.Contents[key]; ok {
        return val
    }

    return nil
}

func NewSysConfig() (common.SYSConfig, error) {
    // init store driver
    path := os.Getenv("JSON_CFG_STORE_PATH")
    if len(path) == 0 {
        path = defaultJSONCfgStorePath
    }

    log.Infof("the path of json configuration storage: %s", path)

    store, err := json.NewCfgStore(path)
    if err != nil {
        return nil, err
    }

    c := &SysConfig{
        StoreDriver: store,
    }

    // init all configuration. this will give the configuration a default value
    c.Init()

    // load configurations
    if err := c.Load(); err != nil {
        return nil, err
    }

    return c, nil
}

var config common.SYSConfig

func GetSysConfig() common.SYSConfig {
    return config
}

func Init() error {
    c, err := NewSysConfig()
    if err != nil {
        return err
    }

    config = c
    return nil
}