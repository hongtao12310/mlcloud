package cluster

import (
    "testing"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/deepinsight/mlcloud/src/pkg/common"
)

func TestBuildKubeConfigFile(t *testing.T) {
    log.Debug("== Test Build Kubeconfig Files ==")

    certDir := config.Get(common.KubeCertDir).(string)
    apiServer := config.Get(common.KubeAPIServer).(string)
    spec := kubeConfigSpec {
        CertDir: certDir,
        ClientName: "test-user",
        APIServer: apiServer,
    }

    log.Debugf("kubeconfig spec: %+v", spec)

    if err := BuildKubeConfigFileIfNotExist(&spec); err != nil {
        log.Fatal("failed to build kubeconfig file")
    }
}

func TestDeleteKubeConfigFile(t *testing.T) {

}


