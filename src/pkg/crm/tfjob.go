package crm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/deepinsight/mlcloud/src/pkg/models"
	"github.com/deepinsight/tf-operator/pkg/spec"
	"github.com/gogo/protobuf/proto"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	DefaultTfImage    = "10.199.192.16/tensorflow/tensorflow:1.2.1"
	DefaultTfGPUImage = "10.199.192.16/tensorflow/tensorflow:1.2.1-gpu"
)

func GenerateTFJob(job *models.Job) (*spec.TfJob, error) {

	if job == nil {
		return nil, fmt.Errorf("Generate TFJob: job is nil")
	}

	if job.Tensorflow == nil {
		return nil, fmt.Errorf("Generate TFJob: tensorflow job is nil")
	}

	typeMeta := metav1.TypeMeta{
		Kind:       spec.CRDKind,
		APIVersion: spec.ApiVersion(),
	}

	var command []string
	var args []string
	var workdir string = "/mlcloud"
	var trainScript []string
	if job.Tensorflow.Command != "" {
		trainScript = strings.Split(job.Tensorflow.Command, "/")
	}
	l := len(trainScript)
	for i := 0; i < l-1; i++ {
		if len(trainScript[i]) > 0 {
			workdir += "/" + trainScript[i]
		}
	}
	command = append(command, "python")
	command = append(command, trainScript[l-1])
	initArgs := strings.Split(job.Tensorflow.Arguments, " ")
	for _, arg := range initArgs {
		argTrim := strings.TrimSpace(arg)
		if len(argTrim) > 0 {
			args = append(args, argTrim)
		}
	}

	volumeName := "mlcloud-" + job.User.Username
	volumeMounts := []v1.VolumeMount{{Name: volumeName, MountPath: "/mlcloud"}}
	volumes := []v1.Volume{{
		Name: volumeName,
		VolumeSource: v1.VolumeSource{
			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
				ClaimName: volumeName,
			},
		},
	}}

	// Set image
	image := DefaultTfImage
	if job.Tensorflow.NumGPU > 0 {
		image = DefaultTfGPUImage
	}
	if job.Tensorflow.Image != "" {
		image = job.Tensorflow.Image
	}

	resources := v1.ResourceRequirements{}
	tolerations := []v1.Toleration{}

	if job.Tensorflow.NumGPU > 0 {
		// Set GPU resources
		quantity, err := k8sresource.ParseQuantity(strconv.Itoa(job.Tensorflow.NumGPU))
		if err != nil {
			return nil, fmt.Errorf("Cannot parse GPU: %v", err)
		}

		resources = v1.ResourceRequirements{
			Limits: v1.ResourceList{
				"alpha.kubernetes.io/nvidia-gpu": quantity,
			},
		}

		// Set tolerations
		toleration := v1.Toleration{
			Key:    "gpu",
			Effect: v1.TaintEffectNoSchedule,
		}
		tolerations = append(tolerations, toleration)
	}

	// Container
	container := v1.Container{
		Name:         "tensorflow",
		Image:        image,
		Command:      command,
		Args:         args,
		WorkingDir:   workdir,
		VolumeMounts: volumeMounts,
		Resources:    resources,
	}

	// podSpec
	podSpec := v1.PodSpec{
		Containers:    []v1.Container{container},
		RestartPolicy: v1.RestartPolicyOnFailure,
		Volumes:       volumes,
		Tolerations:   tolerations,
	}

	// TFJobSpec
	var masterSpec *spec.TfReplicaSpec
	if job.Tensorflow.HasMaster {
		masterSpec = &spec.TfReplicaSpec{
			Replicas: proto.Int32(1),
			Template: &v1.PodTemplateSpec{
				Spec: podSpec,
			},
			TfPort:        proto.Int32(2222),
			TfReplicaType: spec.MASTER,
		}
	}

	var psSpec *spec.TfReplicaSpec
	if job.Tensorflow.NumPs > 0 {
		psSpec = &spec.TfReplicaSpec{
			Replicas: proto.Int32(int32(job.Tensorflow.NumPs)),
			Template: &v1.PodTemplateSpec{
				Spec: podSpec,
			},
			TfPort:        proto.Int32(2222),
			TfReplicaType: spec.PS,
		}
	}

	var workerSpec *spec.TfReplicaSpec
	if job.Tensorflow.NumWorkers > 0 {
		workerSpec = &spec.TfReplicaSpec{
			Replicas: proto.Int32(int32(job.Tensorflow.NumWorkers)),
			Template: &v1.PodTemplateSpec{
				Spec: podSpec,
			},
			TfPort:        proto.Int32(2222),
			TfReplicaType: spec.WORKER,
		}
	}

	replicas := []*spec.TfReplicaSpec{}
	if masterSpec != nil {
		replicas = append(replicas, masterSpec)
	}
	if psSpec != nil {
		replicas = append(replicas, psSpec)
	}
	if workerSpec != nil {
		replicas = append(replicas, workerSpec)
	}

	// Generate Tensorboard if requested
	var tensorboardSpec *spec.TensorBoardSpec
	if job.Tensorflow.Tensorboard {
		job.Tensorflow.TensorboardHost = job.Name + "." + job.User.Username + ".test-cloud.bigdata.wanda.cn"
		tensorboardSpec = &spec.TensorBoardSpec{
			Image:        "10.199.192.16/tensorflow/tensorflow:1.2.1",
			LogDir:       job.Tensorflow.LogDir,
			Host:         job.Tensorflow.TensorboardHost,
			VolumeMounts: volumeMounts,
			Volumes:      volumes,
		}
	}

	tfJobSpec := spec.TfJobSpec{
		ReplicaSpecs: replicas,
		TensorBoard:  tensorboardSpec,
	}

	tfjob := spec.TfJob{
		TypeMeta: typeMeta,
		Metadata: metav1.ObjectMeta{
			Name:      job.Name,
			Namespace: job.User.Username,
		},
		Spec: tfJobSpec,
	}

	return &tfjob, nil
}
