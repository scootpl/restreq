package restreq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (r *Request) do(method string) (*Response, error) {
	var c httpClient

	if r.client == nil {
		c = &http.Client{
			Timeout: r.timeout,
		}
	} else {
		c = r.client
	}

	payload := &bytes.Buffer{}

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

	body := &bytes.Buffer{}
	if !r.bodyReader {
		if _, err = io.Copy(body, resp.Body); err != nil {
			return nil, err
		}
		resp.Body.Close()
	}

	return &Response{
		Response: resp,
		Body:     body.Bytes(),
	}, nil
}
