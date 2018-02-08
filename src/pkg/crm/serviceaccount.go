package crm

import (
    "k8s.io/client-go/pkg/api/v1"
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/deepinsight/mlcloud/src/pkg/cluster"
    "github.com/deepinsight/mlcloud/src/utils"
)

// CreateOrUpdateServiceAccount creates a ServiceAccount if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateMLCloudServiceAccount(user *models.User) error {
    clientSet, err := cluster.GetClientSet()
    if err != nil {
        return err
    }

    sa := v1.ServiceAccount{}
    sa.Namespace = user.Username
    sa.Name = "mlcloud"

    return utils.CreateOrUpdateServiceAccount(clientSet, &sa)
}

func DeleteMLCloudServiceAccount(user *models.User) error {
    clientSet, err := cluster.GetClientSet()
    if err != nil {
        return err
    }

    ns := user.Username
    saname := "mlcloud"

    return utils.DeleteServiceAccount(clientSet, ns, saname)
}