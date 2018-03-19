package client

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newRESTClientTestByURL(t *testing.T, u string) {
	assert := assert.New(t)

	aBaseURL, aErr := url.Parse(u)
	bRestClient, bErr := NewRESTClient(u)

	if nil == aErr && nil == bErr {
		aRestClient := &RESTClient{
			base:   aBaseURL,
			Client: http.DefaultClient,
		}
		assert.Equal(aRestClient, bRestClient)
	} else {
		assert.Equal(aErr, bErr)
	}
}

func TestNewRESTClient(t *testing.T) {
	newRESTClientTestByURL(t, "//valid")
	newRESTClientTestByURL(t, "//||invalid")
}

//TODO: Add tests for CRUD functions if CI environment supported clusters
