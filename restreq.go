package restreq

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type result struct {
	*http.Response
}

type requester interface {
	Context(context.Context) requester
	AddHeader(string, string) requester
	AddCookie(*http.Cookie) requester
	AddJSONKeyValue(string, any) requester
	SetTimeoutSec(int) requester
	SetUserAgent(string) requester
	SetContentType(string) requester
	SetContentTypeJSON() requester
	SetJSONPayload(any) requester
	SetBasicAuth(username, password string) requester
	Post() (*result, error)
	Put() (*result, error)
	Patch() (*result, error)
	Get() (*result, error)
	Delete() (*result, error)
}

func (r *result) Header(s string) string {
	return r.Response.Header.Get(s)
}

func (r *result) GetDecodedJSON(s any) error {
	return json.NewDecoder(r.Body).Decode(&s)
}

type req struct {
	ctx         context.Context
	timeout     time.Duration
	url         string
	json        map[string]any
	headers     map[string]string
	cookies     map[string]*http.Cookie
	username    string
	password    string
	jsonPayload []byte
}

func New(u string) *req {
	return &req{
		url:     u,
		json:    make(map[string]any),
		headers: make(map[string]string),
		cookies: make(map[string]*http.Cookie),
	}
}

func (r *req) SetJSONPayload(p any) requester {
	w := bytes.NewBuffer([]byte{})
	json.NewEncoder(w).Encode(p)
	r.jsonPayload = w.Bytes()
	return r
}

func (r *req) SetBasicAuth(username, password string) requester {
	r.username = username
	r.password = password
	return r
}

func (r *req) AddCookie(c *http.Cookie) requester {
	r.cookies[c.Name] = c
	return r
}

func (r *req) SetContentType(s string) requester {
	r.headers["Content-Type"] = s
	return r
}

func (r *req) SetContentTypeJSON() requester {
	r.headers["Content-Type"] = "application/json"
	return r
}

func (r *req) SetUserAgent(s string) requester {
	r.headers["User-Agent"] = s
	return r
}

func (r *req) Context(ctx context.Context) requester {
	r.ctx = ctx
	return r
}

func (r *req) SetTimeoutSec(t int) requester {
	r.timeout = time.Second * time.Duration(t)
	return r
}

func (r *req) AddHeader(k string, v string) requester {
	r.headers[k] = v
	return r
}

func (r *req) AddJSONKeyValue(key string, value any) requester {
	if key == "" || value == "" {
		return r
	}

	r.json[key] = value
	return r
}

func (r *req) Post() (*result, error) {
	c := http.Client{
		Timeout: r.timeout,
	}
	return r.do("POST", &c)
}

func (r *req) Get() (*result, error) {
	c := http.Client{
		Timeout: r.timeout,
	}
	return r.do("GET", &c)
}

func (r *req) Delete() (*result, error) {
	c := http.Client{
		Timeout: r.timeout,
	}
	return r.do("DELETE", &c)
}

func (r *req) Patch() (*result, error) {
	c := http.Client{
		Timeout: r.timeout,
	}
	return r.do("PATCH", &c)
}

func (r *req) Put() (*result, error) {
	c := http.Client{
		Timeout: r.timeout,
	}
	return r.do("PUT", &c)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (r *req) do(method string, c HTTPClient) (*result, error) {
	payload := bytes.NewBuffer([]byte{})

	if len(r.jsonPayload) > 0 {
		payload.Write(r.jsonPayload)
	} else {
		if err := json.NewEncoder(payload).Encode(r.json); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, r.url, payload)
	if err != nil {
		return nil, err
	}

	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	if r.username != "" && r.password != "" {
		r.SetBasicAuth(r.username, r.password)
	}

	for _, v := range r.cookies {
		r.AddCookie(v)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return &result{
		Response: resp,
	}, nil
}
