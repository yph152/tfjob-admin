// /*
// Copyright 2017 caicloud authors. All rights reserved.
// */

package client

import (
	"encoding/json"
	"net/http"
	"net/url"
	//"github.com/asaskevich/govalidator"
)

// RESTClient is a wrapper of http.Client,
// making it easier to write client for rest apis in ai-serving
type RESTClient struct {
	base *url.URL

	// Set specific behavior of the client.
	// If not set http.DefaultClient will be used.
	Client *http.Client
}

// encode encodes object using json
func encode(object interface{}) ([]byte, error) {
	return json.Marshal(object)
}

// decode decodes bytes into object using json
func decode(data []byte, object interface{}) error {
	return json.Unmarshal(data, object)
}

// NewRESTClient create new client.
// Endpoint is base url to access other api,
// eg. http://127.0.0.1:7799
func NewRESTClient(endpoint string) (*RESTClient, error) {
	// ret := govalidator.IsURL(endpoint)
	// if !ret {
	// 	return nil, fmt.Errorf("endpoint %s is not valid URL", endpoint)
	// }
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return &RESTClient{
		base:   baseURL,
		Client: http.DefaultClient,
	}, nil
}

const (
	// HTTPVerbPost post verb
	HTTPVerbPost string = "POST"
	// HTTPVerbPut put verb
	HTTPVerbPut string = "PUT"
	// HTTPVerbGet get verb
	HTTPVerbGet string = "GET"
	// HTTPVerbDelete delete verb
	HTTPVerbDelete string = "DELETE"
	// HTTPVerbPatch patch verb
	HTTPVerbPatch string = "PATCH"
)

// Verb set request verb
func (c *RESTClient) Verb(verb string) *Request {
	return NewRequest(c.Client, verb, c.base)
}

// Post send a post request
func (c *RESTClient) Post() *Request {
	return c.Verb(HTTPVerbPost)
}

// Put send a put request
func (c *RESTClient) Put() *Request {
	return c.Verb(HTTPVerbPut)
}

// Get send a get request
func (c *RESTClient) Get() *Request {
	return c.Verb(HTTPVerbGet)
}

// Delete delete a get request
func (c *RESTClient) Delete() *Request {
	return c.Verb(HTTPVerbDelete)
}

// Patch delete a patch request
func (c *RESTClient) Patch() *Request {
	return c.Verb(HTTPVerbPatch)
}
