package storage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/caicloud/kubeflow-admin/pkg/apis/v1alpha1/clientset/client"
	//	"github.com/caicloud/kubeflow-admin/pkg/log"
)

type StorageClient struct {
	restClient *client.RESTClient
}

type GetStorageResourcesResult struct {
	StorageList *ListStorageResponse `json:"storageList"`
	Err         error                `json:"error"`
	StatusCode  *int                 `json:"statusCode"`
}

type ListStorageResponse struct {
	StorageResources []StorageResource `json:"storageResources"`
	SubUsers         []SubUserInfo     `json:"subUsers"`
}

type SubUserInfo struct {
	UserID   string `json:"userID,omitempty"`
	UserName string `json:"userName,omitempty"`
}

type StorageResource struct {
	ID           int64         `json:"id,omitempty"`
	Capacity     int64         `json:"capacity,omitempty"`
	Type         StorageType   `json:"type,omitempty"`
	Status       StorageStatus `json:"status,omitempty"`
	RootPath     string        `json:"rootPath,omitempty"`
	StorageName  string        `json:"storageName,omitempty"`
	CreateTime   time.Time     `json:"createTime,omitempty"`
	WebSocketURL string        `json:"webSocketURL,omitempty"`
}

type StorageType string

const (
	TypeIDStorageFileSystem              = 1
	TypeStrStorageFileSystem StorageType = "filesystem"
)

type StorageStatus string

const (
	StatusIDStorageNew       = 1 // next is Allocated
	StatusIDStorageAllocated = 2 // next is Free
	StatusIDStorageFree      = 3 // next is Available
	StatusIDStorageAvailable = 4 // next is Deleted
	StatusIDStorageDisable   = 5
	StatusIDStorageDeleted   = 6 // next is Recycling
	StatusIDStorageRecycling = 7 // next is Free

	StatusIDStoragePending = StatusIDStorageFree      // same to Free, named in user's view
	StatusIDStorageSuccess = StatusIDStorageAvailable // same to Available, named in user's view

	StatusStrStorageNew       StorageStatus = "new"       // just insert in db
	StatusStrStorageAllocated StorageStatus = "allocated" // real resource has been allocated ok
	StatusStrStorageFree      StorageStatus = "free"      // bounded with file server app and app running ok, this resource is free for user to allocate
	StatusStrStorageAvailable StorageStatus = "available" // bounded with user success, is available to use for user
	StatusStrStorageDisable   StorageStatus = "disable"   // disabled by admin
	StatusStrStorageDeleted   StorageStatus = "deleted"   // mark deleted by user
	StatusStrStorageRecycling StorageStatus = "recycling" // being recycling in background, will be set to free after finished

	StatusStrStoragePending StorageStatus = "pending" // same to Free, named in user's view
	StatusStrStorageSuccess StorageStatus = "success" // same to Available, named in user's view
)

type GetStorageResult struct {
	StorageInfo *StorageInfo `json:"storageInfo"`
	Err         error        `json:"error"`
	StatusCode  *int         `json:"statusCode"`
}

type StorageInfo struct {
	StorageID int64  `json:"storageID"`
	Namespace string `json:"namespace"`
	PVCName   string `json:"pvcName"`
	MountPath string `json:"mountPath"`
	ReadOnly  bool   `json:"readOnly"`
}

func NewStorageClient(host string) (*StorageClient, error) {
	restClient, err := client.NewRESTClient(host)
	if err != nil {
		return nil, err
	}

	return &StorageClient{
		restClient: restClient,
	}, nil
}

func (sc *StorageClient) GetStorageResources(UID string) *GetStorageResourcesResult {
	//	log.Infof(fmt.Sprintf("Try getting storageID for user(%s)", UID))
	storageResourcesResult := &GetStorageResourcesResult{}
	storageResourcesResult.StorageList = &ListStorageResponse{}
	storageResourcesResult.Err = sc.restClient.Get().
		SetHeader("Uid", UID).
		SetPath(fmt.Sprintf("/api/v3/storages")).
		Do().
		Into(storageResourcesResult.StorageList)

		//	log.Infof(fmt.Sprintf("storageID: %d", storageResourcesResult.StorageList.StorageResources[0].ID))

	return storageResourcesResult
}

func (sc *StorageClient) GetStorage(storageID int64, namespace string, UID string) *GetStorageResult {
	//	log.Infof(fmt.Sprintf("Try getting storage PVC for user(%s)", UID))
	storageResult := &GetStorageResult{}
	storageResult.StorageInfo = &StorageInfo{}
	storageResult.Err = sc.restClient.Post().
		SetHeader("Uid", UID).
		SetPath(fmt.Sprintf("/api/v3/storages/%s/usageBundle", strconv.Itoa(int(storageID)))).
		SetParam("namespace", namespace).
		Do().
		Into(storageResult.StorageInfo)

	return storageResult
}
