package crm

import (
    "fmt"
    extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
    apierrors "k8s.io/apimachinery/pkg/api/errors"
    clientset "k8s.io/client-go/kubernetes"
)

// CreateOrUpdateDeployment creates a Deployment if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
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

// CreateOrUpdateDaemonSet creates a DaemonSet if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
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
