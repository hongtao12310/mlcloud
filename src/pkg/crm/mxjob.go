package crm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/deepinsight/mlcloud/src/pkg/models"
	"github.com/deepinsight/mxnet-operator/pkg/spec"
	"github.com/gogo/protobuf/proto"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	DefaultMxImage    = "10.199.192.16/library/mxnet-gpu-dist:0.12.0"
	DefaultMxGPUImage = "10.199.192.16/library/mxnet-gpu-dist:0.12.0"
)

func GenerateMXJob(job *models.Job) (*spec.MxJob, error) {

	if job == nil {
		return nil, fmt.Errorf("Generate MXJob: job is nil")
	}
	if job.Mxnet == nil {
		return nil, fmt.Errorf("Generate MXJob: mxnet job is nil")
	}

	typeMeta := metav1.TypeMeta{
		Kind:       spec.CRDKind,
		APIVersion: spec.CRDApiVersion,
	}

	var command []string
	var args []string
	var workdir string = "/mlcloud"
	var trainScript []string
	if job.Mxnet.Command != "" {
		trainScript = strings.Split(job.Mxnet.Command, "/")
	}
	l := len(trainScript)
	for i := 0; i < l-1; i++ {
		if len(trainScript[i]) > 0 {
			workdir += "/" + trainScript[i]
		}
	}
	command = append(command, "python")
	command = append(command, trainScript[l-1])
	initArgs := strings.Split(job.Mxnet.Arguments, " ")
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
	image := DefaultMxImage
	if job.Mxnet.NumGPU > 0 {
		image = DefaultMxGPUImage
	}
	if job.Mxnet.Image != "" {
		image = job.Mxnet.Image
	}

	resources := v1.ResourceRequirements{}
	tolerations := []v1.Toleration{}

	if job.Mxnet.NumGPU > 0 {
		// Set GPU resources
		quantity, err := k8sresource.ParseQuantity(strconv.Itoa(job.Mxnet.NumGPU))
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
		Name:         "mxnet",
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

	// MXJobSpec
	var schedulerSpec *spec.MxReplicaSpec
	if job.Mxnet.Mode == "dist" {
		schedulerSpec = &spec.MxReplicaSpec{
			Replicas: proto.Int32(1),
			Template: &v1.PodTemplateSpec{
				Spec: podSpec,
			},
			PsRootPort:    proto.Int32(9091),
			MxReplicaType: spec.SCHEDULER,
		}
	}

	var psSpec *spec.MxReplicaSpec
	if job.Mxnet.NumPs > 0 {
		psSpec = &spec.MxReplicaSpec{
			Replicas: proto.Int32(int32(job.Mxnet.NumPs)),
			Template: &v1.PodTemplateSpec{
				Spec: podSpec,
			},
			MxReplicaType: spec.SERVER,
		}
	}

	var workerSpec *spec.MxReplicaSpec
	if job.Mxnet.NumWorkers > 0 {
		workerSpec = &spec.MxReplicaSpec{
			Replicas: proto.Int32(int32(job.Mxnet.NumWorkers)),
			Template: &v1.PodTemplateSpec{
				Spec: podSpec,
			},
			MxReplicaType: spec.WORKER,
		}
	}

	replicas := []*spec.MxReplicaSpec{}
	if schedulerSpec != nil {
		replicas = append(replicas, schedulerSpec)
	}
	if psSpec != nil {
		replicas = append(replicas, psSpec)
	}
	if workerSpec != nil {
		replicas = append(replicas, workerSpec)
	}

	var mxJobSpec spec.MxJobSpec
	if job.Mxnet.Mode == "local" {
		mxJobSpec = spec.MxJobSpec{
			ReplicaSpecs: replicas,
			JobMode:      spec.LocalJob,
		}
	}

	if job.Mxnet.Mode == "dist" {
		mxJobSpec = spec.MxJobSpec{
			ReplicaSpecs: replicas,
			JobMode:      spec.DistJob,
		}
	}

	mxjob := spec.MxJob{
		TypeMeta: typeMeta,
		Metadata: metav1.ObjectMeta{
			Name:      job.Name,
			Namespace: job.User.Username,
		},
		Spec: mxJobSpec,
	}

	return &mxjob, nil
}
