package restreq

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Response inherits from http.Response, so you can use almost every
// field and method of http.Response.
//
// http.Response.Body is an exception. You cannot use it, because
// content of http.Response.Body is copied to Response.Body.
// You don't have to call http.Response.Body.Close()
type Response struct {
	*http.Response
	Body []byte
}

// Header returns header
func (r *Response) Header(s string) string {
	return r.Response.Header.Get(s)
}

// DecodeJSON decodes JSON
func (r *Response) DecodeJSON(s any) error {
	b := bytes.NewReader(r.Body)
	return json.NewDecoder(b).Decode(&s)
}

type requester interface {
	Context(context.Context) requester
	SetHTTPClient(httpClient) requester
	AddHeader(string, string) requester
	AddCookie(*http.Cookie) requester
	AddJSONKeyValue(string, any) requester
	SetTimeoutSec(int) requester
	SetUserAgent(string) requester
	SetContentType(string) requester
	SetContentTypeJSON() requester
	SetJSONPayload(any) requester
	SetBasicAuth(username, password string) requester
	Debug(*log.Logger, DebugFlag) requester
	WithBodyReader() requester
	Post() (*Response, error)
	Put() (*Response, error)
	Patch() (*Response, error)
	Get() (*Response, error)
	Delete() (*Response, error)
}

// Request contains all methods to operate on REST API
type Request struct {
	ctx         context.Context
	timeout     time.Duration
	url         string
	json        map[string]any
	headers     map[string]string
	cookies     map[string]*http.Cookie
	username    string
	password    string
	jsonPayload []byte
	client      httpClient
	debugFlags  int32
	logger      *log.Logger
	bodyReader  bool
}

func New(u string) *Request {
	return &Request{
		url:     u,
		json:    make(map[string]any),
		headers: make(map[string]string),
		cookies: make(map[string]*http.Cookie),
	}
}

type DebugFlag int32

// DebugFlags to control logger behavior.
const (
	// Debug request body
	ReqBody DebugFlag = 1 << iota
	// Debug request headers
	ReqHeaders
	// Debug request cookies
	ReqCookies
	// Debug response body
	RespBody
	// Debug response header
	RespHeaders
	// Debug response cookies
	RespCookies
)

// WithBodyReader allows direct reading from http.Response.Body without
// copying to restreq.Response.Body
func (r *Request) WithBodyReader() requester {
	r.bodyReader = true
	return r
}

// SetHTTPClient sets external http client.
func (r *Request) SetHTTPClient(c httpClient) requester {
	r.client = c
	return r
}

// Debug sets logger and debug flags.
// You can combine flags, ReqBody+ReqHeader etc.
func (r *Request) Debug(logger *log.Logger, flags DebugFlag) requester {
	r.debugFlags = int32(flags)
	r.logger = logger
	return r
}

// SetJSONPayload encodes map or struct to json byte array.
func (r *Request) SetJSONPayload(p any) requester {
	w := bytes.NewBuffer([]byte{})
	json.NewEncoder(w).Encode(p)
	r.jsonPayload = w.Bytes()
	return r
}

// SetBasicAuth sets basic auth with username and password.
func (r *Request) SetBasicAuth(username, password string) requester {
	r.username = username
	r.password = password
	return r
}

// AddCookie adds cookie to request.
func (r *Request) AddCookie(c *http.Cookie) requester {
	r.cookies[c.Name] = c
	return r
}

// SetContentType sets Content-Type.
func (r *Request) SetContentType(s string) requester {
	r.headers["Content-Type"] = s
	return r
}

// SetContentTypeJSON sets Content-Type to application/json.
func (r *Request) SetContentTypeJSON() requester {
	r.headers["Content-Type"] = "application/json"
	return r
}

// SetUserAgent sets User-Agent header.
func (r *Request) SetUserAgent(s string) requester {
	r.headers["User-Agent"] = s
	return r
}

// Context sets context to ctx
func (r *Request) Context(ctx context.Context) requester {
	r.ctx = ctx
	return r
}

// SetTimeoutSec sets connection timeout.
func (r *Request) SetTimeoutSec(t int) requester {
	r.timeout = time.Second * time.Duration(t)
	return r
}

// AddHeader adds header with value.
func (r *Request) AddHeader(k string, v string) requester {
	r.headers[k] = v
	return r
}

// AddJSONKeyValue converts KV to json byte array.
// You can add many KV, they will be added to the map
// and converted to an byte array when the request is sent.
func (r *Request) AddJSONKeyValue(key string, value any) requester {
	if key == "" || value == "" {
		return r
	}

	r.json[key] = value
	return r
}

// Post executes the post method
func (r *Request) Post() (*Response, error) {
	return r.do("POST")
}

// Get executes the get method
func (r *Request) Get() (*Response, error) {
	return r.do("GET")
}

// Delete executes the delete method
func (r *Request) Delete() (*Response, error) {
	return r.do("DELETE")
}

// Patch executes the patch method
func (r *Request) Patch() (*Response, error) {
	return r.do("PATCH")
}

// Put executes the put method
func (r *Request) Put() (*Response, error) {
	return r.do("PUT")
}

func (r *Request) debug(f DebugFlag, s string) {
	if r.logger == nil || r.debugFlags&(1<<(f-1)) == 0 {
		return
	}

	r.logger.Printf("%s\n", s)
}
