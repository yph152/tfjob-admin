package handlers

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	//	"fmt"
	tfjobtypes "tfjob-admin/pkg/api/v1/types"
	tfv1alpha2 "tfjob-admin/pkg/apis/tensorflow/v1alpha2"
	tfjobclient "tfjob-admin/pkg/tfjobs"
	//	"fmt"
)

var tfc *tfjobclient.TFJobClient

func NewClient(kubeconfig string) {
	tfc, _ = tfjobclient.NewTFJobClient(kubeconfig)
}

func convert(importconfig *tfjobtypes.ImportConfig) *tfv1alpha2.TFJob {

	var replica int32 = 1
	jobs := &tfv1alpha2.TFJob{}
	jobs.APIVersion = "kubeflow.org/v1alpha2"
	jobs.Kind = "TFjob"
	jobs.ObjectMeta.Name = importconfig.Jobid
	jobs.Namespace = importconfig.Partition
	jobs.Spec.TFReplicaSpecs["WORKER"].Replicas = &replica
	jobs.Spec.TFReplicaSpecs["WORKER"].Template.Spec.Containers[0].Name = "tfjob"
	jobs.Spec.TFReplicaSpecs["WORKER"].Template.Spec.Containers[0].Image = importconfig.Images
	jobs.Spec.TFReplicaSpecs["WORKER"].Template.Spec.Containers[0].Command = importconfig.Command
	jobs.Spec.TFReplicaSpecs["WORKER"].Template.Spec.RestartPolicy = "OnFailure"

	return jobs

}
func GetJob(ctx context.Context, cid, paritions, jobid string) (*tfv1alpha2.TFJob, error) {
	tfjob, err := tfc.GetJob(cid, paritions, jobid)

	return tfjob, err
}

func CreateJob(ctx context.Context, cid, paritions string, valuesReader io.ReadCloser) (*tfv1alpha2.TFJob, error) {

	defer valuesReader.Close()

	b, err := ioutil.ReadAll(valuesReader)

	if err != nil {
		return nil, err
	}

	var importconfig = &tfjobtypes.ImportConfig{}
	err = json.Unmarshal(b, importconfig)

	if err != nil {
		return nil, err
	}

	tfjobconfig := convert(importconfig)
	jobs, err := tfc.CreateJob(cid, paritions, tfjobconfig)

	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func DeleteJob(ctx context.Context, cid, partitions, jobid string) (bool, error) {
	err := tfc.DeleteJobs(cid, partitions, jobid)

	if err != nil {
		return false, err
	}
	return true, nil
}

func UpdateJob(ctx context.Context, cid, paritions string, valuesReader io.ReadCloser) (*tfv1alpha2.TFJob, error) {

	defer valuesReader.Close()

	b, err := ioutil.ReadAll(valuesReader)

	if err != nil {
		return nil, err
	}

	var importconfig = &tfjobtypes.ImportConfig{}
	err = json.Unmarshal(b, importconfig)

	if err != nil {
		return nil, err
	}

	tfjobconfig := convert(importconfig)
	jobs, err := tfc.UpdateJob(cid, paritions, tfjobconfig)

	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func ListJobs(ctx context.Context, cid, partitions string) (*tfv1alpha2.TFJobList, error) {
	listjobs, err := tfc.ListJobs(cid, partitions)

	return listjobs, err
}
