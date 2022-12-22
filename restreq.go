package restreq

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type result struct {
	*http.Response
}

type requester interface {
	Ctx(context.Context) requester
	SetTimeoutSec(int) requester
	SetHeader(string) requester
	JSON(string) requester
	Post() (*result, error)
	Get() (*result, error)
	Delete() (*result, error)
}

func (result) JSON(string) (any, bool) {
	return nil, false
}

func (r *result) Header(s string) string {
	return r.Response.Header.Get(s)
}

type req struct {
	ctx     context.Context
	timeout time.Duration
	url     string
	json    map[string]any
	headers map[string]string
}

func New(u string) *req {
	r := req{
		url:     u,
		json:    make(map[string]any),
		headers: make(map[string]string),
	}
	return &r
}

func (r *req) Ctx(ctx context.Context) requester {
	r.ctx = ctx
	return r
}

func (r *req) SetTimeoutSec(t int) requester {
	r.timeout = time.Second * time.Duration(t)
	return r
}

func (r *req) SetHeader(h string) requester {
	key, value, ok := strings.Cut(h, "=")
	if ok && key != "" && value != "" {
		r.headers[key] = value
	}
	return r
}

func (r *req) JSON(i string) requester {
	key, value, ok := strings.Cut(i, ":=")
	if ok && (key == "" || value == "") {
		return r
	}
	if ok {
		vb, err := strconv.ParseBool(value)
		if err == nil {
			r.json[key] = vb
			return r
		}

		vi, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			r.json[key] = vi
			return r
		}

		vf, err := strconv.ParseFloat(value, 64)
		if err == nil {
			r.json[key] = vf
		}
		return r
	}

	key, value, ok = strings.Cut(i, "=")
	if ok && key != "" && value != "" {
		r.json[key] = value
	}
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

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (r *req) do(method string, c HTTPClient) (*result, error) {
	var (
		request *http.Request
		err     error
	)

	payload := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(payload).Encode(r.json); err != nil {
		return nil, err
	}

	if r.ctx != nil {
		request, err = http.NewRequestWithContext(r.ctx, method, r.url, payload)
	} else {
		request, err = http.NewRequest(method, r.url, payload)
	}

	if err != nil {
		return nil, err
	}

	for k, v := range r.headers {
		request.Header.Set(k, v)
	}

	resp, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	return &result{
		Response: resp,
	}, nil
}
