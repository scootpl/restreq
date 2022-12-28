# RestReq

[![Go Reference](https://pkg.go.dev/badge/pkg.go.dev/github.com/scootpl/restreq.svg)](https://pkg.go.dev/pkg.go.dev/github.com/scootpl/restreq)

RestReq is a wrapper around standard Go net/http client. In a simple call you can use json encoding, add headers
and parse result. This should be sufficient in most use cases.

## Examples

- Simplest use

```go
    resp, err := restreq.New("http://example.com").Post()
```

- You can add a header

```go
    resp, err := restreq.New("http://example.com").
		AddHeader("X-TOKEN", authToken).
		Post()
```

- Use map with JSON payload

 ```go
    p := map[string]any{
		"string": "string",
		"bool": true,
		"float": 2.34,
	}

	resp, err := restreq.New("http://example.com").
		SetContentTypeJSON().
		SetJSONPayload(p).
		Post()
```

- JSON payload with KV

 ```go
	resp, err := restreq.New("http://example.com").
		SetContentTypeJSON().
		SetUserAgent("Client 1.0").
		AddJSONKeyValue("string", "string").
		AddJSONKeyValue("bool", true).
		AddJSONKeyValue("float", 2.34).
		Post()
```

## Parsing response

- Get header value

```go
	value := resp.Header("token")
```

- Decode JSON

```go
	s := struct {
		Message string `json:"message,omitempty"`
	}{}

	err := resp.DecodeJSON(&s)
```

## Request Methods

- Add a context to the request
```go
    Context(context.Context)
```

- Set an external httpClient. SetTimeoutSec() doesn't work in this case
```go
	SetHTTPClient(httpClient)
```

- Add header, you can repeat this method to add multiple headers
```go
	AddHeader(string, string)
```

- Add cookie, you can repeat this method to add multiple cookies
```go
	AddCookie(*http.Cookie)
```
- Add JSON key-value, you can repeat this method to add multiple tokens
```go
	AddJSONKeyValue(string, any)
```

- Set connection timeout in seconds
```go
	SetTimeoutSec(int)
```

- Set User-Agent header
```go
    SetUserAgent(string)
```

 - Set Content-Type header
 ```go
    SetContentType(string)
```
 
 - Set Content-Type to application/json
 ```go
    SetContentTypeJSON()
```

- Set JSON payload, encodes map or struct to json byte array.
```go
	SetJSONPayload(any) requester
```

- Set basic auth with username and password
```go
	SetBasicAuth(username, password string)
```

- Set logger and debug level
```go
	Debug(*log.Logger, DebugFlag)
```

```go
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
```

- Execute a method
```go
	Post() (*Response, error)
	Put() (*Response, error)
	Patch() (*Response, error)
	Get() (*Response, error)
	Delete() (*Response, error)
```

## Response methods

- Get header value
```go
    Header(s string) string 
```

- Decode JSON reply
```go
    DecodeJSON(s any) error
```





