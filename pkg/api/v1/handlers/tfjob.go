package handlers

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"k8s.io/api/core/v1"

	//	"fmt""
	tfjobtypes "tfjob-admin/pkg/api/v1/types"
	tfv1alpha2 "tfjob-admin/pkg/apis/tensorflow/v1alpha2"
	storageclient "tfjob-admin/pkg/storage"
	tfjobclient "tfjob-admin/pkg/tfjobs"
	//	"fmt"
)

var tfc *tfjobclient.TFJobClient
var sc *storageclient.StorageClient

func NewClient(kubeconfig, storagehost string) {
	tfc, _ = tfjobclient.NewTFJobClient(kubeconfig)
	sc, _ = storageclient.NewStorageClient(storagehost)
}

func convert(importconfig *tfjobtypes.ImportConfig) *tfv1alpha2.TFJob {

	storageSpec, _ := getsStorage(importconfig.UserID, importconfig.Partition)
	var replica int32 = 1
	jobs := new(tfv1alpha2.TFJob)
	//	jobs.APIVersion = "kubeflow.org/v1alpha2"
	//	jobs.Kind = "tfjob"
	jobs.ObjectMeta.Name = importconfig.Jobid
	jobs.Namespace = importconfig.Partition
	jobs.Spec.TFReplicaSpecs = make(map[tfv1alpha2.TFReplicaType]*tfv1alpha2.TFReplicaSpec)
	jobs.Spec.TFReplicaSpecs["Worker"] = &tfv1alpha2.TFReplicaSpec{}
	jobs.Spec.TFReplicaSpecs["Worker"].Replicas = &replica
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.Containers = make([]v1.Container, 1)
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.Containers[0].Name = "tfjob"
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.Containers[0].Image = importconfig.Images
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.Containers[0].Command = importconfig.Command
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.RestartPolicy = "OnFailure"
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.Volumes = make([]v1.Volume, 1)
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.Volumes[0].Name = "path"
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.Volumes[0].VolumeSource.PersistentVolumeClaim = &v1.PersistentVolumeClaimVolumeSource{}
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.Volumes[0].VolumeSource.PersistentVolumeClaim.ClaimName = storageSpec.PersistentVolumeClaim
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.Containers[0].VolumeMounts = make([]v1.VolumeMount, 1)
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.Containers[0].VolumeMounts[0].Name = "path"
	jobs.Spec.TFReplicaSpecs["Worker"].Template.Spec.Containers[0].VolumeMounts[0].MountPath = storageSpec.MountPath
	return jobs

}

func getsStorage(userid, partition string) (*tfjobtypes.StorageSpec, error) {
	storageResourcesResult := sc.GetStorageResources(userid)

	if storageResourcesResult.Err != nil {
		return nil, storageResourcesResult.Err
	}

	storageID := storageResourcesResult.StorageList.StorageResources[0].ID

	storageResult := sc.GetStorage(storageID, partition, userid)

	if storageResult.Err != nil {
		return nil, storageResult.Err
	}

	storage := &tfjobtypes.StorageSpec{
		StorageID:             storageID,
		MountPath:             storageResult.StorageInfo.MountPath,
		ReadOnly:              storageResult.StorageInfo.ReadOnly,
		PersistentVolumeClaim: storageResult.StorageInfo.PVCName,
	}
	return storage, nil
}
func GetJob(ctx context.Context, cid, partitions, jobid string) (*tfv1alpha2.TFJob, error) {
	tfjob, err := tfc.GetJob(cid, partitions, jobid)

	return tfjob, err
}

func CreateJob(ctx context.Context, cid, partitions string, valuesReader io.ReadCloser) (*tfv1alpha2.TFJob, error) {

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
	jobs, err := tfc.CreateJob(cid, partitions, tfjobconfig)

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

func UpdateJob(ctx context.Context, cid, partitions, jobid string, valuesReader io.ReadCloser) (*tfv1alpha2.TFJob, error) {

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
	jobs, err := tfc.UpdateJob(cid, partitions, tfjobconfig)

	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func ListJobs(ctx context.Context, cid, partitions string) (*tfv1alpha2.TFJobList, error) {
	listjobs, err := tfc.ListJobs(cid, partitions)

	return listjobs, err
}
