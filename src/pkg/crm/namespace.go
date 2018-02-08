// package crm - cluster resource management
// used to talk to kubernetes cluster to manager user's resources

package crm

import (
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/deepinsight/mlcloud/src/utils"
    "github.com/deepinsight/mlcloud/src/pkg/cluster"
    "k8s.io/client-go/pkg/api/v1"
)

func CreateUserNamespace(user *models.User) error  {

    clientSet, err := cluster.GetClientSet()
    if err != nil {
        return err
    }

    ns := &v1.Namespace{}
    ns.Name = user.Username
    return utils.CreateOrUpdateNamespace(clientSet, ns)
}

func DeleteUserNamespace(user *models.User) error  {
    clientSet, err := cluster.GetClientSet()
    if err != nil {
        return err
    }
    return utils.DeleteNamespace(clientSet, user.Username)
}
