package client

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	assert := assert.New(t)
	tempStr := "test"
	tempClient := http.DefaultClient

	tempBase := &url.URL{
		Scheme:     "test",
		Opaque:     "test",
		User:       &url.Userinfo{},
		Host:       "test",
		Path:       "test",
		RawPath:    "test",
		ForceQuery: false,
		RawQuery:   "",
		Fragment:   "",
	}

	aReq := &Request{
		client:  tempClient,
		verb:    tempStr,
		baseURL: tempBase,
	}

	bReq := NewRequest(tempClient, tempStr, tempBase)
	assert.Equal(aReq, bReq)
}

func newTempReq() *Request {
	tempClient := http.DefaultClient

	tempBase := &url.URL{
		Scheme:     "test",
		Opaque:     "test",
		User:       &url.Userinfo{},
		Host:       "test",
		Path:       "test",
		RawPath:    "test",
		ForceQuery: false,
		RawQuery:   "",
		Fragment:   "",
	}

	return &Request{
		client:  tempClient,
		verb:    "test",
		baseURL: tempBase,
	}
}

func TestSetPath(t *testing.T) {
	aReq := newTempReq()
	bReq := newTempReq()
	tempPath := "test"

	aReq.path = tempPath
	bReq = bReq.SetPath(tempPath)

	assert.Equal(t, aReq, bReq)
}

func TestSetHeaders(t *testing.T) {
	aReq := newTempReq()
	bReq := aReq
	tempHeaders := map[string]string{
		"a": "aa",
		"b": "bb",
	}

	if aReq.headers == nil {
		aReq.headers = http.Header{}
	}
	for k, v := range tempHeaders {
		aReq.headers.Set(k, v)
	}

	bReq = bReq.SetHeaders(tempHeaders)
	assert.Equal(t, aReq, bReq)
}

func TestSetHeader(t *testing.T) {
	aReq := newTempReq()
	bReq := aReq
	tempK := "testKey"
	tempV := "testVal"

	if aReq.headers == nil {
		aReq.headers = http.Header{}
	}
	aReq.headers.Set(tempK, tempV)

	bReq = bReq.SetHeader(tempK, tempV)
	assert.Equal(t, aReq, bReq)
}

func TestSetParams(t *testing.T) {
	aReq := newTempReq()
	bReq := aReq
	tempParams := map[string]string{
		"a": "aa",
		"b": "bb",
	}

	if aReq.headers == nil {
		aReq.params = make(url.Values)
	}
	for p, v := range tempParams {
		aReq.params[p] = append(aReq.params[p], v)
	}

	bReq = bReq.SetParams(tempParams)
	assert.Equal(t, aReq, bReq)
}

//TODO: Add tests for CRUD functions if CI environment supported clusters
