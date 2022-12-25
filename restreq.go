/*
Restreq is a wrapper for net/http client.
You can easily create requests.

Body content is copied to Response.Body, so you don't have to call resp.Body.Close()

Examples:

1)

	resp, err := restreq.New("http://").Post()

2)

	resp, err := restreq.New("http://").
		AddHeader(token, authToken).
		Post()

3)

	p := map[string]any{
		"key": "value",
	}

	resp, err := restreq.New("http://").
		Context(ctx).
		SetTimeoutSec(30).
		SetContentTypeJSON().
		SetJSONPayload(p).
		Post()
*/
package restreq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Response inherits from *http.Response
// and adds some new methods
type Response struct {
	*http.Response
	Body []byte
}

// Header returns s header
func (r *Response) Header(s string) string {
	return r.Response.Header.Get(s)
}

// DecodeJSON decodes JSON from the response body
// to the s struct
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
	Post() (*Response, error)
	Put() (*Response, error)
	Patch() (*Response, error)
	Get() (*Response, error)
	Delete() (*Response, error)
}

// Request contains all the methods to operate on REST API
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

const (
	ReqBody DebugFlag = 1 << iota
	ReqHeaders
	ReqCookies
	RespBody
	RespHeaders
	RespCookies
)

// SetHTTPClient sets http client
func (r *Request) SetHTTPClient(c httpClient) requester {
	r.client = c
	return r
}

// Debug sets logger and debug flags
func (r *Request) Debug(logger *log.Logger, flags DebugFlag) requester {
	r.debugFlags = int32(flags)
	r.logger = logger
	return r
}

// SetJSONPayload encodes json
func (r *Request) SetJSONPayload(p any) requester {
	w := bytes.NewBuffer([]byte{})
	json.NewEncoder(w).Encode(p)
	r.jsonPayload = w.Bytes()
	return r
}

// SetBasicAuth sets basic auth with username and password
func (r *Request) SetBasicAuth(username, password string) requester {
	r.username = username
	r.password = password
	return r
}

// AddCookie adds cookie to request
func (r *Request) AddCookie(c *http.Cookie) requester {
	r.cookies[c.Name] = c
	return r
}

// SetContentType sets Content-Type
func (r *Request) SetContentType(s string) requester {
	r.headers["Content-Type"] = s
	return r
}

// SetContentTypeJSON sets Content-Type to application/json
func (r *Request) SetContentTypeJSON() requester {
	r.headers["Content-Type"] = "application/json"
	return r
}

// SetUserAgent sets User-Agent to s
func (r *Request) SetUserAgent(s string) requester {
	r.headers["User-Agent"] = s
	return r
}

// Context sets context to ctx
func (r *Request) Context(ctx context.Context) requester {
	r.ctx = ctx
	return r
}

// SetTimeoutSec sets connection timeout to t seconds
func (r *Request) SetTimeoutSec(t int) requester {
	r.timeout = time.Second * time.Duration(t)
	return r
}

// AddHeader adds k header with v value
func (r *Request) AddHeader(k string, v string) requester {
	r.headers[k] = v
	return r
}

// AddJSONKeyValue converts key/value to json
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

func (r *Request) do(method string) (*Response, error) {
	var c httpClient

	if r.client == nil {
		c = &http.Client{
			Timeout: r.timeout,
		}
	} else {
		c = r.client
	}

	payload := bytes.NewBuffer([]byte{})

	if len(r.jsonPayload) > 0 {
		payload.Write(r.jsonPayload)
	} else {
		if err := json.NewEncoder(payload).Encode(r.json); err != nil {
			return nil, err
		}
	}

	r.debug(ReqBody, fmt.Sprintf("Body: %s", strings.TrimRight(payload.String(), "\n")))

	req, err := http.NewRequest(method, r.url, payload)
	if err != nil {
		return nil, err
	}

	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
		r.debug(ReqHeaders, fmt.Sprintf("Header: %s: %s", k, v))
	}

	if r.username != "" && r.password != "" {
		r.SetBasicAuth(r.username, r.password)
	}

	for k, v := range r.cookies {
		r.AddCookie(v)
		r.debug(ReqCookies, fmt.Sprintf("Cookie: %s: %s", k, v))
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer([]byte{})
	if _, err = io.Copy(body, resp.Body); err != nil {
		return nil, err
	}
	resp.Body.Close()

	return &Response{
		Response: resp,
		Body:     body.Bytes(),
	}, nil
}
