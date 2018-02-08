package utils

import (
    "fmt"

    "k8s.io/client-go/pkg/api/v1"
    meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
    rbac "k8s.io/client-go/pkg/apis/rbac/v1beta1"
    apierrors "k8s.io/apimachinery/pkg/api/errors"
    clientset "k8s.io/client-go/kubernetes"
    "github.com/deepinsight/mlcloud/src/utils/log"
)

func CreateOrUpdateNamespace(client clientset.Interface, ns *v1.Namespace) error {
    if _, err := client.CoreV1().Namespaces().Create(ns); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create namespace: %v", err)
        }

        if _, err := client.CoreV1().Namespaces().Update(ns); err != nil {
            return fmt.Errorf("unable to update namespace: %v", err)
        }
    }
    return nil
}

func DeleteNamespace(client clientset.Interface, ns string) error {
    if err := client.CoreV1().Namespaces().Delete(ns, &meta_v1.DeleteOptions{}); err != nil {
        if !apierrors.IsNotFound(err) {
            return err
        }
    }

    return nil
}

func CreateCephFSSecret(client clientset.Interface, namespace string, secret *v1.Secret) error {
    if _, err := client.CoreV1().Secrets(namespace).Create(secret); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create secret: %v", err)
        }
    }
    return nil
}

func CreateCephFSPV(client clientset.Interface, pv *v1.PersistentVolume) error {
    if _, err := client.CoreV1().PersistentVolumes().Create(pv); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create pv: %v", err)
        }
    }
    return nil
}

func DeleteCephFSPV(client clientset.Interface, name string) error {
    if err := client.CoreV1().PersistentVolumes().Delete(name, nil); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create pv: %v", err)
        }
    }
    return nil
}

func CreateCephFSPVC(client clientset.Interface, namespace string, pvc *v1.PersistentVolumeClaim) error {
    if _, err := client.CoreV1().PersistentVolumeClaims(namespace).Create(pvc); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create pvc: %v", err)
        }
    }
    return nil
}

func DeleteCephFSPVC(client clientset.Interface, namespace string, name string) error {
    if err := client.CoreV1().PersistentVolumeClaims(namespace).Delete(name, nil); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create pvc: %v", err)
        }
    }
    return nil
}

func ListPod(client clientset.Interface, ns string) error {
    podList, err := client.CoreV1().Pods(ns).List(meta_v1.ListOptions{})
    if err != nil {
        return err
    }

    log.Debugf("there are %d pods in namespace: %s", len(podList.Items), ns)

    for _, pod := range podList.Items {
        log.Debugf("pod name: %s", pod.Name)
    }

    return nil
}

// CreateOrUpdateConfigMap creates a ConfigMap if the target resource doesn't exist.
// If the resource exists already, this function will update the resource instead.
func CreateOrUpdateConfigMap(client clientset.Interface, cm *v1.ConfigMap) error {
    if _, err := client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Create(cm); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create configmap: %v", err)
        }

        if _, err := client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Update(cm); err != nil {
            return fmt.Errorf("unable to update configmap: %v", err)
        }
    }
    return nil
}

// CreateOrUpdateServiceAccount creates a ServiceAccount if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateServiceAccount(client clientset.Interface, sa *v1.ServiceAccount) error {
    if _, err := client.CoreV1().ServiceAccounts(sa.ObjectMeta.Namespace).Create(sa); err != nil {
        // Note: We don't run .Update here afterwards as that's probably not required
        // Only thing that could be updated is annotations/labels in .metadata, but we don't use that currently
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create serviceaccount: %v", err)
        }
    }
    return nil
}

func DeleteServiceAccount(client clientset.Interface, ns, saname string) error {
    if err := client.CoreV1().ServiceAccounts(ns).Delete(saname, &meta_v1.DeleteOptions{}); err != nil {
        if !apierrors.IsNotFound(err) {
            return err
        }
    }
    return nil
}

// CreateOrUpdateDeployment creates a Deployment if the target resource doesn't exist.
// If the resource exists already, this function will update the resource instead.
func CreateOrUpdateDeployment(client clientset.Interface, deploy *extensions.Deployment) error {
    if _, err := client.ExtensionsV1beta1().Deployments(deploy.ObjectMeta.Namespace).Create(deploy); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create deployment: %v", err)
        }

        if _, err := client.ExtensionsV1beta1().Deployments(deploy.ObjectMeta.Namespace).Update(deploy); err != nil {
            return fmt.Errorf("unable to update deployment: %v", err)
        }
    }
    return nil
}

// CreateOrUpdateDaemonSet creates a DaemonSet if the target resource doesn't exist.
// If the resource exists already, this function will update the resource instead.
func CreateOrUpdateDaemonSet(client clientset.Interface, ds *extensions.DaemonSet) error {
    if _, err := client.ExtensionsV1beta1().DaemonSets(ds.ObjectMeta.Namespace).Create(ds); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create daemonset: %v", err)
        }

        if _, err := client.ExtensionsV1beta1().DaemonSets(ds.ObjectMeta.Namespace).Update(ds); err != nil {
            return fmt.Errorf("unable to update daemonset: %v", err)
        }
    }
    return nil
}

// CreateOrUpdateRole creates a Role if the target resource doesn't exist.
// If the resource exists already, this function will update the resource instead.
func CreateOrUpdateRole(client clientset.Interface, role *rbac.Role) error {
    if _, err := client.RbacV1beta1().Roles(role.ObjectMeta.Namespace).Create(role); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create RBAC role: %v", err)
        }

        if _, err := client.RbacV1beta1().Roles(role.ObjectMeta.Namespace).Update(role); err != nil {
            return fmt.Errorf("unable to update RBAC role: %v", err)
        }
    }
    return nil
}

// CreateOrUpdateRoleBinding creates a RoleBinding if the target resource doesn't exist.
// If the resource exists already, this function will update the resource instead.
func CreateOrUpdateRoleBinding(client clientset.Interface, roleBinding *rbac.RoleBinding) error {
    if _, err := client.RbacV1beta1().RoleBindings(roleBinding.ObjectMeta.Namespace).Create(roleBinding); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create RBAC rolebinding: %v", err)
        }

        if _, err := client.RbacV1beta1().RoleBindings(roleBinding.ObjectMeta.Namespace).Update(roleBinding); err != nil {
            return fmt.Errorf("unable to update RBAC rolebinding: %v", err)
        }
    }
    return nil
}

func DeleteRoleBinding(client clientset.Interface, ns, rolebindName string) error {
    if err := client.RbacV1beta1().RoleBindings(ns).Delete(rolebindName, &meta_v1.DeleteOptions{}); err != nil {
        if !apierrors.IsNotFound(err) {
            return err
        }
    }
    return nil
}

// CreateOrUpdateClusterRole creates a ClusterRole if the target resource doesn't exist.
// If the resource exists already, this function will update the resource instead.
func CreateOrUpdateClusterRole(client clientset.Interface, clusterRole *rbac.ClusterRole) error {
    if _, err := client.RbacV1beta1().ClusterRoles().Create(clusterRole); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create RBAC clusterrole: %v", err)
        }

        if _, err := client.RbacV1beta1().ClusterRoles().Update(clusterRole); err != nil {
            return fmt.Errorf("unable to update RBAC clusterrole: %v", err)
        }
    }
    return nil
}


func DeleteClusterRoleBinding(client clientset.Interface, rolebindName string) error {
    if err := client.RbacV1beta1().ClusterRoleBindings().Delete(rolebindName, &meta_v1.DeleteOptions{}); err != nil {
        if !apierrors.IsNotFound(err) {
            return err
        }
    }
    return nil
}
// CreateOrUpdateClusterRoleBinding creates a ClusterRoleBinding if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateClusterRoleBinding(client clientset.Interface, clusterRoleBinding *rbac.ClusterRoleBinding) error {
    if _, err := client.RbacV1beta1().ClusterRoleBindings().Create(clusterRoleBinding); err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("unable to create RBAC clusterrolebinding: %v", err)
        }

        if _, err := client.RbacV1beta1().ClusterRoleBindings().Update(clusterRoleBinding); err != nil {
            return fmt.Errorf("unable to update RBAC clusterrolebinding: %v", err)
        }
    }
    return nil
}
