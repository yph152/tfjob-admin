package storage

import (
	"testing"

	"github.com/caicloud/kubeflow-admin/pkg/apis/v1alpha1/clientset/client"
	"github.com/stretchr/testify/assert"
)

func newNewStorageClientTestByURL(t *testing.T, u string) {
	assert := assert.New(t)

	aRestClient, aErr := client.NewRESTClient(u)
	bStorageClient, bErr := NewStorageClient(u)

	if nil == aErr && nil == bErr {
		aStorageClient := &StorageClient{
			restClient: aRestClient,
		}
		assert.Equal(aStorageClient, bStorageClient)
	} else {
		assert.Equal(aErr, bErr)
	}
}

func TestNewStorageClient(t *testing.T) {
	newNewStorageClientTestByURL(t, "//valid")
	newNewStorageClientTestByURL(t, "//||invalid")
}

//TODO: Add tests for CRUD functions if CI environment supported clusters
