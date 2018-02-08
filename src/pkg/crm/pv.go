package crm

import (
	"github.com/deepinsight/mlcloud/src/pkg/models"
	"fmt"
	"github.com/golang/glog"
	"github.com/deepinsight/mlcloud/src/pkg/cluster"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	"github.com/deepinsight/mlcloud/src/utils"
)

var capacity, _ = k8sresource.ParseQuantity("100G")


func CreateUserPV(user *models.User) error {
	if user == nil {
		glog.V(4).Infof("user is nil")
		return fmt.Errorf("user is nil")
	}

	clientSet, err := cluster.GetClientSet()
	if err != nil {
		return err
	}

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mlcloud-" + user.Username,
			Labels: map[string]string{
				"app": "mlcloud",
				"user": user.Username,
			},
		},
		Spec: v1.PersistentVolumeSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			Capacity: map[v1.ResourceName]k8sresource.Quantity{
				v1.ResourceStorage: capacity,
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				CephFS: &v1.CephFSVolumeSource{
					Monitors: CephFS,
					Path: "/mlcloud/" + user.Username,
					SecretRef: &v1.LocalObjectReference{
						Name: CephFSSecretName,
					},
					User: CephFSUser,
				},
			},
			PersistentVolumeReclaimPolicy: v1.PersistentVolumeReclaimRetain,
		},
	}

	return utils.CreateCephFSPV(clientSet, pv)
}

func DeleteUserPV(user *models.User) error {
	if user == nil {
		glog.V(4).Infof("user is nil")
		return fmt.Errorf("user is nil")
	}

	clientSet, err := cluster.GetClientSet()
	if err != nil {
		return err
	}

	return utils.DeleteCephFSPV(clientSet, "mlcloud-"+user.Username)
}

func CreateUserPVC(user *models.User) error {
	if user == nil {
		glog.V(4).Infof("user is nil")
		return fmt.Errorf("user is nil")
	}

	clientSet, err := cluster.GetClientSet()
	if err != nil {
		return err
	}

	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mlcloud-" + user.Username,
			Namespace: user.Username,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			Resources: v1.ResourceRequirements{
				Requests: map[v1.ResourceName]k8sresource.Quantity{
					v1.ResourceStorage: capacity,
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "mlcloud",
					"user": user.Username,
				},
			},
		},
	}

	return utils.CreateCephFSPVC(clientSet, user.Username, pvc)
}

func DeleteUserPVC(user *models.User) error {
	if user == nil {
		glog.V(4).Infof("user is nil")
		return fmt.Errorf("user is nil")
	}

	clientSet, err := cluster.GetClientSet()
	if err != nil {
		return err
	}

	return utils.DeleteCephFSPVC(clientSet, user.Username, "mlcloud-"+user.Username)
}