package crm

import (
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "k8s.io/client-go/pkg/apis/rbac/v1beta1"
    "github.com/deepinsight/mlcloud/src/utils"
    "github.com/deepinsight/mlcloud/src/pkg/cluster"
)

func CreateUserRolebinding(user *models.User) error {
    clientSet, err := cluster.GetClientSet()
    if err != nil {
        return err
    }

    roleBinding := &v1beta1.RoleBinding{}
    roleBinding.Name = user.Username + "-rb"
    roleBinding.Namespace = user.Username
    roleBinding.RoleRef = v1beta1.RoleRef{Name: "admin", Kind: "ClusterRole"}
    roleBinding.Subjects = []v1beta1.Subject{}
    subject := v1beta1.Subject{Kind: "User", Name: user.Username}
    roleBinding.Subjects = append(roleBinding.Subjects, subject)

    return utils.CreateOrUpdateRoleBinding(clientSet, roleBinding)
}

func DeleteUserRolebinding(user *models.User) error {
    clientSet, err := cluster.GetClientSet()
    if err != nil {
        return err
    }

    rolebindingName := user.Username + "-rb"
    ns := user.Username
    return utils.DeleteRoleBinding(clientSet, ns, rolebindingName)
}


func CreateUserClusterRolebinding(user *models.User) error {
    clientSet, err := cluster.GetClientSet()
    if err != nil {
        return err
    }

    clusterRoleBinding := &v1beta1.ClusterRoleBinding{}
    clusterRoleBinding.Name = user.Username + "-crb"
    clusterRoleBinding.RoleRef = v1beta1.RoleRef{Name: "cluster-admin", Kind: "ClusterRole"}
    clusterRoleBinding.Subjects = []v1beta1.Subject{}
    subject := v1beta1.Subject{Kind: "User", Name: user.Username}
    clusterRoleBinding.Subjects = append(clusterRoleBinding.Subjects, subject)

    return utils.CreateOrUpdateClusterRoleBinding(clientSet, clusterRoleBinding)
}

func DeleteUserClusterRolebinding(user *models.User) error {
    clientSet, err := cluster.GetClientSet()
    if err != nil {
        return err
    }

    clusterRoleBindingName := user.Username + "-crb"
    return utils.DeleteClusterRoleBinding(clientSet, clusterRoleBindingName)
}
