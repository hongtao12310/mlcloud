package cluster

import (
    "testing"
    "github.com/deepinsight/mlcloud/src/pkg/common"
    "github.com/deepinsight/mlcloud/src/pkg/common/syscfg"
    "os"
    "errors"
    "github.com/deepinsight/mlcloud/src/utils/log"
)

var (
    defaultJSONCfgStorePath  = os.Getenv("HOME") + "/mlcloud/config/config.json"
    defaultKubeCertDIR  = os.Getenv("HOME") + "/mlcloud/certs"
)

func setENVs()  {
    os.Setenv("MYSQL_HOST", "localhost")
    os.Setenv("MYSQL_PORT", "3306")
    os.Setenv("MYSQL_USER", "mlcloud")
    os.Setenv("MYSQL_PASSWORD", "mlcloud")
    os.Setenv("MYSQL_DATABASE", "mlcloud")
    os.Setenv("KUBE_APISERVER", "https://10.214.160.9:6443")
    os.Setenv("KUBE_CLUSTER_NAME", "kubernetes")
    os.Setenv("KUBE_CERT_DIR", defaultKubeCertDIR)

}


var config common.SYSConfig


func initSysConfigForTest() error {
    setENVs()
    log.Debugf("MYSQL_HOST: %s", os.Getenv("MYSQL_HOST"))
    c := make(map[string]interface{})
    c["MYSQL_HOST"] = os.Getenv("MYSQL_HOST")
    config = syscfg.GetSysConfig()
    if config == nil {
        return errors.New("failed to init sysconfig")
    }

    return nil
}

func TestMain(m *testing.M) {
    // init system config
    if err := initSysConfigForTest(); err != nil {
        log.Fatal(err)
    }

    // init kubernetes client


    result := m.Run()
    if result != 0 {
        os.Exit(result)
    }
}

