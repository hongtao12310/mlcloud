package crm

import (
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/deepinsight/mlcloud/src/pkg/cluster"
    "github.com/deepinsight/mlcloud/src/utils"
)

func ListPod(user *models.User) error  {
    clientSet, err := cluster.GetClientSet()
    if err != nil {
        return err
    }

    return  utils.ListPod(clientSet, user.Username)

}
