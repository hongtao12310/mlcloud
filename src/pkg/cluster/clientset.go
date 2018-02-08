package cluster

import (
	mxk8sutil "github.com/deepinsight/mxnet-operator/pkg/util/k8sutil"
	tfk8sutil "github.com/deepinsight/tf-operator/pkg/util/k8sutil"
	"k8s.io/client-go/kubernetes"
)

var clientset *ClientSet

type ClientSet struct {
	*kubernetes.Clientset
	tfclient *tfk8sutil.TfJobRestClient
	mxclient *mxk8sutil.MxJobRestClient
}

func (clientset *ClientSet) TFJob() tfk8sutil.TfJobClient {
	if clientset == nil {
		return nil
	}

	return clientset.tfclient
}

func (clientset *ClientSet) MXJob() mxk8sutil.MxJobClient {
	if clientset == nil {
		return nil
	}

	return clientset.mxclient
}

func newClientSet() (*ClientSet, error) {

	config, err := tfk8sutil.GetClusterConfig()
	if err != nil {
		return nil, err
	}

	k8sclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	tfclient, err := tfk8sutil.NewTfJobClientExternal(config)
	if err != nil {
		return nil, err
	}

	mxclient, err := mxk8sutil.NewMxJobClientExternal(config)
	if err != nil {
		return nil, err
	}

	clientset = &ClientSet{k8sclientset, tfclient, mxclient}
	return clientset, nil
}

func GetClientSet() (*ClientSet, error) {
	if clientset == nil {
		return newClientSet()
	} else {
		return clientset, nil
	}
}
