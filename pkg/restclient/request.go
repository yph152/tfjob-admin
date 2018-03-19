package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/caicloud/kubeflow-admin/pkg/log"
)

// HTTPClient is an interface for testing a request object.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Request allows for building up a request to a server in a chained fashion.
// Any clientset are stored until the end of your call, so you only have to
// check once.
type Request struct {
	// required
	client HTTPClient
	verb   string

	baseURL *url.URL

	// generic components accessible via method setters
	path    string
	params  url.Values
	headers http.Header

	// output
	err  error
	body io.Reader
}

// NewRequest creates a new request helper object for accessing kubernetes-admin api.
func NewRequest(client HTTPClient, verb string, baseURL *url.URL) *Request {
	r := &Request{
		client:  client,
		verb:    verb,
		baseURL: baseURL,
	}
	return r
}

// SetPath set request path
func (r *Request) SetPath(path string) *Request {
	r.path = path
	return r
}

// SetHeaders set request headers
func (r *Request) SetHeaders(headers map[string]string) *Request {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	for k, v := range headers {
		r.headers.Set(k, v)
	}
	return r
}

// SetHeader set request header
func (r *Request) SetHeader(key, value string) *Request {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Set(key, value)
	return r
}

// SetParams set request params
func (r *Request) SetParams(params map[string]string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	for paramName, value := range params {
		r.params[paramName] = append(r.params[paramName], value)
	}
	return r
}

// SetParam set request param
func (r *Request) SetParam(paramName, value string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[paramName] = append(r.params[paramName], value)
	return r
}

const (
	// ContentTypeJSON json content type
	ContentTypeJSON = "application/json"
)

// Body set request body
func (r *Request) Body(obj interface{}) *Request {
	if r.err != nil {
		return r
	}
	switch t := obj.(type) {
	case string:
		r.err = fmt.Errorf("[k8s-admin client] cannot use string as request body")
		return r
	case []byte:
		//		log.Debugf("Request Body: %#v", string(t))
		r.body = bytes.NewReader(t)
	case io.Reader:
		r.body = t
	default:
		// callers may pass typed interface pointers, therefore we must check nil with reflection
		if reflect.ValueOf(t).IsNil() {
			return r
		}
		data, err := encode(t)
		if err != nil {
			//			log.Errorf("failed to encode: %#+v", t)
			r.err = err
			return r
		}
		//		log.Debugf("Request Body: %#v", string(data))
		r.body = bytes.NewReader(data)
		r.SetHeader("Content-Type", ContentTypeJSON)
	}
	return r
}

// request wraps building a http request, sending and closing it, and calling fn to handle the result
func (r *Request) request(fn func(*http.Request, *http.Response)) error {
	c := r.client
	if c == nil {
		c = http.DefaultClient
	}

	url := r.URL().String()
	//	log.Println(r.path, url)
	req, err := http.NewRequest(r.verb, url, r.body)
	if err != nil {
		return err
	}
	req.Header = r.headers

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	// Ensure the response body is fully read and closed.
	defer func() {
		const maxBodySlurpSize = 2 * 1024
		if resp.ContentLength <= maxBodySlurpSize {
			io.Copy(ioutil.Discard, &io.LimitedReader{R: resp.Body, N: maxBodySlurpSize}) // nolint: errcheck
		}
		resp.Body.Close() // nolint: errcheck
	}()

	// Callback to handle result
	fn(req, resp)
	return nil
}

// Do will send out request and put decode http response to Result object
func (r *Request) Do() Result {
	var result Result
	err := r.request(func(req *http.Request, resp *http.Response) {
		result = r.transformResponse(resp, req)
	})
	if err != nil {
		return Result{err: err}
	}
	return result
}

// transformResponse converts an API response into a structured API object
func (r *Request) transformResponse(resp *http.Response, req *http.Request) Result {
	var body []byte
	if resp.Body != nil {
		if data, err := ioutil.ReadAll(resp.Body); err == nil {
			body = data
		}
	}

	//	log.Debugf("Response Body: %s", string(body))

	// verify the content type is accurate
	contentType := resp.Header.Get("Content-Type")
	if contentType != ContentTypeJSON {
		if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent {
			return Result{
				statusCode: resp.StatusCode,
				err:        r.transformUnstructuredResponseError(resp, req, body),
			}
		}
		return Result{
			body:        body,
			contentType: contentType,
			statusCode:  resp.StatusCode,
		}
	}
	//	log.Debugf("Response Content Type: %s", contentType)

	// If the server gave us a status back, look at what it was.
	success := resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusPartialContent
	if !success {
		apiError := &Error{}
		err := decode(body, apiError)
		if err != nil {
			// failed to decode the error, just return result
			//			log.Errorf("Failed to decode api error, body: %s", string(body))
			return Result{
				body:        body,
				contentType: contentType,
				statusCode:  resp.StatusCode,
			}
		}
		// "Failed" requests are clearly just an error and it makes sense to return them as such.
		return Result{err: apiError}
	}

	return Result{
		body:        body,
		contentType: contentType,
		statusCode:  resp.StatusCode,
	}
}

// transformUnstructuredResponseError handles an error from the server that is not in a structured form.
// It is expected to transform any response that is not recognizable as a clear server sent error from the
// K8S API using the information provided with the request. In practice, HTTP proxies and client libraries
// introduce a level of uncertainty to the responses returned by servers that in common use result in
// unexpected responses. The rough structure is:
//
// 1. Assume the server sends you something sane - JSON + well defined error objects + proper codes
//    - this is the happy path
//    - when you get this output, trust what the server sends
// 2. Guard against empty fields / bodies in received JSON and attempt to cull sufficient info from them to
//    generate a reasonable facsimile of the original failure.
//    - Be sure to use a distinct error type or flag that allows a client to distinguish between this and error 1 above
// 3. Handle true disconnect failures / completely malformed data by moving up to a more generic client error
// 4. Distinguish between various connection failures like SSL certificates, timeouts, proxy clientset, unexpected
//    initial contact, the presence of mismatched body contents from posted content types
//    - Give these a separate distinct error type and capture as much as possible of the original message
//
// TODO: introduce transformation of generic http.Client.Do() clientset that separates 4.
func (r *Request) transformUnstructuredResponseError(resp *http.Response, req *http.Request, body []byte) error {
	if body == nil && resp.Body != nil {
		if data, err := ioutil.ReadAll(resp.Body); err == nil {
			body = data
		}
	}

	//	log.Debugf("Response Body: %#v", string(body))

	message := "unknown"
	if isTextResponse(resp) {
		message = strings.TrimSpace(string(body))
	}

	return fmt.Errorf("server response %v, method: %v, message: %v",
		resp.StatusCode,
		req.Method,
		message)
}

// isTextResponse returns true if the response appears to be a textual media type.
func isTextResponse(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	if len(contentType) == 0 {
		return true
	}
	media, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return strings.HasPrefix(media, "text/")
}

// URL returns the current working URL.
func (r *Request) URL() *url.URL {
	p := r.path
	if !strings.HasPrefix(p, "/") {
		r.path = "/" + r.path
	}

	finalURL := &url.URL{}
	if r.baseURL != nil {
		*finalURL = *r.baseURL
	}
	finalURL.Path = finalURL.Path + r.path

	query := url.Values{}
	for key, values := range r.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	finalURL.RawQuery = query.Encode()
	return finalURL
}

// Result contains the result of calling Request.Do().
type Result struct {
	body        []byte
	contentType string
	err         error
	statusCode  int
}

// Into stores the result into obj, if possible. If obj is nil it is ignored.
func (r Result) Into(obj interface{}) error {
	if r.err != nil {
		return r.err
	}
	return decode(r.body, obj)
}

func (r Result) Error() error {
	return r.err
}

// StatusCode returns the HTTP status code of the request. (Only valid if no
// error was returned.)
func (r Result) StatusCode(statusCode *int) Result {
	*statusCode = r.statusCode
	return r
}

// GetBody return
func (r Result) GetBody() []byte {
	return r.body
}
