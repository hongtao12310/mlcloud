package cluster

import (
    "bytes"
    "crypto/x509"
    "fmt"
    "os"
    "path/filepath"
    "crypto/rsa"
    "k8s.io/client-go/tools/clientcmd"
    clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
    certutil "k8s.io/client-go/util/cert"
    util "github.com/deepinsight/mlcloud/src/utils"
    "github.com/deepinsight/mlcloud/src/pkg/certs/pkiutil"
    "github.com/deepinsight/mlcloud/src/utils/log"

    "github.com/astaxie/beego/utils"
)

// clientCertAuth struct holds info required to build a client certificate to provide authentication info in a kubeconfig object
type clientCertAuth struct {
    CAKey         *rsa.PrivateKey
    Organizations []string
}

// tokenAuth struct holds info required to use a token to provide authentication info in a kubeconfig object
type tokenAuth struct {
    Token string
}

// kubeConfigSpec struct holds info required to build a KubeConfig object
type kubeConfigSpec struct {
    CACert         *x509.Certificate
    CertDir        string
    APIServer      string
    ClientName     string
    GroupName      string
    ClusterName    string
    TokenAuth      *tokenAuth
    ClientCertAuth *clientCertAuth
}

// createKubeConfigFiles creates all the requested kubeconfig files.
// If kubeconfig files already exists, they are used only if evaluated equal; otherwise an error is returned.
func BuildKubeConfigFileIfNotExist(spec *kubeConfigSpec) error {
    kubeConfigFileName := spec.ClientName + ".kubeconfig"
    kubeConfigPath := filepath.Join(spec.CertDir, kubeConfigFileName)
    if utils.FileExists(kubeConfigPath) {
        log.Debugf("kubeconfig file: %s already exist", kubeConfigPath)
        return nil
    }

    // builds the KubeConfig object
    config, err := buildKubeConfigFromSpec(spec)
    if err != nil {
        return err
    }

    // writes the KubeConfig to disk if it not exists
    return createKubeConfigFileIfNotExists(kubeConfigPath, config)

}

func DeleteKubeConfigFile(spec *kubeConfigSpec) error {
    kubeConfigFileName := spec.ClientName + ".kubeconfig"
    kubeConfigPath := filepath.Join(spec.CertDir, kubeConfigFileName)
    if !utils.FileExists(kubeConfigPath) {
        log.Debugf("kubeconfig file: %s doesn't already exist", kubeConfigPath)
        return nil
    }

    return os.Remove(kubeConfigPath)
}

func loadCACertFiles(spec *kubeConfigSpec) error {
    caCert, caKey, err := pkiutil.TryLoadCertAndKeyFromDisk(spec.CertDir, "ca")
    if err != nil {
        return err
    }

    spec.CACert = caCert
    spec.ClientCertAuth = &clientCertAuth{
        CAKey:         caKey,
        //Organizations: []string{cfg.GroupName},
    }

    return nil
}

// buildKubeConfigFromSpec creates a kubeconfig object for the given kubeConfigSpec
func buildKubeConfigFromSpec(spec *kubeConfigSpec) (*clientcmdapi.Config, error) {
    if err := loadCACertFiles(spec); err != nil {
        return nil, err
    }

    // If this kubeconfig should use token
    if spec.TokenAuth != nil {
        // create a kubeconfig with a token
        return util.CreateWithToken(
            spec.APIServer,
            spec.ClusterName,
            spec.ClientName,
            certutil.EncodeCertPEM(spec.CACert),
            spec.TokenAuth.Token,
        ), nil
    }

    // otherwise, create a client certs
    clientCertConfig := certutil.Config{
        CommonName:   spec.ClientName,
        Organization: spec.ClientCertAuth.Organizations,
        Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
    }
    clientCert, clientKey, err := pkiutil.NewCertAndKey(spec.CACert, spec.ClientCertAuth.CAKey, clientCertConfig)
    if err != nil {
        return nil, fmt.Errorf("failure while creating %s client certificate: %v", spec.ClientName, err)
    }

    // create a kubeconfig with the client certs
    return util.CreateWithCerts(
        spec.APIServer,
        spec.ClusterName,
        spec.ClientName,
        certutil.EncodeCertPEM(spec.CACert),
        certutil.EncodePrivateKeyPEM(clientKey),
        certutil.EncodeCertPEM(clientCert),
    ), nil
}

// createKubeConfigFileIfNotExists saves the KubeConfig object into a file if there isn't any file at the given path.
// If there already is a KubeConfig file at the given path; we tries to load it and check if the values in the
// existing and the expected config equals. If they do; we will just skip writing the file as it's up-to-date,
// but if a file exists but has old content or isn't a kubeconfig file, this function returns an error.
func createKubeConfigFileIfNotExists(kubeConfigFilePath string, config *clientcmdapi.Config) error {
    // Check if the file exist, and if it doesn't, just write it to disk
    if _, err := os.Stat(kubeConfigFilePath); os.IsNotExist(err) {
        err = util.WriteToDisk(kubeConfigFilePath, config)
        if err != nil {
            return fmt.Errorf("failed to save kubeconfig file %s on disk: %v", kubeConfigFilePath, err)
        }
        return nil
    }

    // The kubeconfig already exists, let's check if it has got the same CA and server URL
    currentConfig, err := clientcmd.LoadFromFile(kubeConfigFilePath)
    if err != nil {
        return fmt.Errorf("failed to load kubeconfig file %s that already exists on disk: %v", kubeConfigFilePath, err)
    }

    expectedCtx := config.CurrentContext
    expectedCluster := config.Contexts[expectedCtx].Cluster
    currentCtx := currentConfig.CurrentContext
    currentCluster := currentConfig.Contexts[currentCtx].Cluster

    // If the current CA cert on disk doesn't match the expected CA cert, error out because we have a file, but it's stale
    if !bytes.Equal(currentConfig.Clusters[currentCluster].CertificateAuthorityData, config.Clusters[expectedCluster].CertificateAuthorityData) {
        return fmt.Errorf("a kubeconfig file %q exists already but has got the wrong CA cert", kubeConfigFilePath)
    }
    // If the current API Server location on disk doesn't match the expected API server, error out because we have a file, but it's stale
    if currentConfig.Clusters[currentCluster].Server != config.Clusters[expectedCluster].Server {
        return fmt.Errorf("a kubeconfig file %q exists already but has got the wrong API Server URL", kubeConfigFilePath)
    }

    return nil
}

