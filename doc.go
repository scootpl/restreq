/*
RestReq is a wrapper around standard Go net/http client. In a simple call you can use json encoding, add headers
and parse result. This should be sufficient in most use cases.

# Examples

1) Simplest use

	resp, err := restreq.New("http://example.com").Post()

2) You can add a header

	resp, err := restreq.New("http://example.com").
		AddHeader("X-TOKEN", authToken).
		Post()

3) Use map with JSON payload

	p := map[string]any{
		"string": "string",
		"bool": true,
		"float": 2.34,
	}

	resp, err := restreq.New("http://example.com").
		SetContentTypeJSON().
		SetJSONPayload(p).
		Post()

4) JSON payload with KV

	resp, err := restreq.New("http://example.com").
		SetContentTypeJSON().
		SetUserAgent("Client 1.0").
		AddJSONKeyValue("string", "string").
		AddJSONKeyValue("bool", true).
		AddJSONKeyValue("float", 2.34).
		Post()

# Parsing response

1) Get header value

	value := resp.Header("token")

2) Decode JSON

	s := struct {
		Message string `json:"message,omitempty"`
	}{}

	err := resp.DecodeJSON(&s)
*/
package restreq
