package crm

import (
	"github.com/deepinsight/mlcloud/src/pkg/models"
	"github.com/deepinsight/mlcloud/src/pkg/cluster"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"github.com/deepinsight/mlcloud/src/utils"
)

const (
	cephFSType  = "kubernetes.io/rbd"
	// Ceph secret will be stored by base64 encoding
	cephFSSecret = "AQDHz/VYaIaeNBAAt4AcFtC/ErbppUJmWVENoA=="
	CephFSSecretName = "cephfs-secret"
	CephFSUser = "admin"
)

var CephFS = []string {
	"10.214.160.5:6789",
	"10.214.160.6:6789",
	"10.214.160.7:6789",
}

func CreateUserNamespaceCephFSSecret(user *models.User) error {

	if user == nil {
		return fmt.Errorf("user is nil")
	}

	clientSet, err := cluster.GetClientSet()
	if err != nil {
		return err
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: CephFSSecretName,
			Namespace: user.Username,
		},
		StringData: map[string]string{"key": cephFSSecret},
		Type: cephFSType,
	}

	return utils.CreateCephFSSecret(clientSet, user.Username, secret)
}

