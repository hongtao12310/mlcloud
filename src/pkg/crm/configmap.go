package crm

import (
    "fmt"

    "k8s.io/client-go/pkg/api/v1"
    apierrors "k8s.io/apimachinery/pkg/api/errors"
    clientset "k8s.io/client-go/kubernetes"
)

// CreateOrUpdateConfigMap creates a ConfigMap if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
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
