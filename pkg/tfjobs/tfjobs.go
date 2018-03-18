package tfjobs

import (
	//	"encoding/json"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	tfv1alpha2 "tfjob-admin/pkg/apis/tensorflow/v1alpha2"
	tfjobclientset "tfjob-admin/pkg/client/clientset/versioned"
)

type TFJobClient struct {
	tfjobclientset.Interface
}

func NewTFJobClient(kubeconfig string) (*TFJobClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		return nil, err
	}
	TFJobClientSet, err := tfjobclientset.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return &TFJobClient{TFJobClientSet}, nil
}
func (cs *TFJobClient) GetJob(cid, partition, jobid string) (*tfv1alpha2.TFJob, error) {
	tfjob, err := cs.KubeflowV1alpha2().TFJobs(partition).Get(jobid, v1.GetOptions{})

	if err != nil {
		return nil, err
	}
	return tfjob, nil
}

func (cs *TFJobClient) CreateJob(cid, partition string, tfjobconfig *tfv1alpha2.TFJob) (*tfv1alpha2.TFJob, error) {

	tfjob, err := cs.KubeflowV1alpha2().TFJobs(partition).Create(tfjobconfig)

	if err != nil {
		return nil, err
	}
	return tfjob, nil
}

func (cs *TFJobClient) UpdateJob(cid, partition string, tfjobconfig *tfv1alpha2.TFJob) (*tfv1alpha2.TFJob, error) {

	tfjob, err := cs.KubeflowV1alpha2().TFJobs(partition).Update(tfjobconfig)

	if err != nil {
		return nil, err
	}
	return tfjob, nil
}

func (cs *TFJobClient) ListJobs(cid, partition string) (*tfv1alpha2.TFJobList, error) {

	tfjoblist, err := cs.KubeflowV1alpha2().TFJobs(partition).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return tfjoblist, nil
}

func (cs *TFJobClient) DeleteJobs(cid, partition, jobid string) error {
	err := cs.KubeflowV1alpha2().TFJobs(partition).Delete(jobid, &v1.DeleteOptions{})

	return err
}
