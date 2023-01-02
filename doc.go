/*
RestReq is a wrapper around standard Go net/http client. In a simple call you can use json encoding, add headers
and parse result. This should be sufficient in most use cases.

# Examples

- Simplest use

	resp, err := restreq.New("http://example.com").Post()

- You can add a header

	resp, err := restreq.New("http://example.com").
		AddHeader("X-TOKEN", authToken).
		Post()

- Use map with JSON payload

	p := map[string]any{
		"string": "string",
		"bool": true,
		"float": 2.34,
	}

	resp, err := restreq.New("http://example.com").
		SetContentTypeJSON().
		SetJSONPayload(p).
		Post()

- JSON payload with KV

	resp, err := restreq.New("http://example.com").
		SetContentTypeJSON().
		SetUserAgent("Client 1.0").
		AddJSONKeyValue("string", "string").
		AddJSONKeyValue("bool", true).
		AddJSONKeyValue("float", 2.34).
		Post()

# Parsing response

- In the default behavior, the body of the response is copied to Response.Body, and you don't have to
call http.Response.Body.Close()

	resp, err := restreq.New("http://example.com").Post()

	if err == nil {
		fmt.Printf("%s\n", resp.Body)
	}

- Default behavior is convenient but not optimal, due to redundant copying. If you need high performance,
you can disable this behavior and direct access to the io.Reader. Don't forget to call Response.Body.Close()

	resp, err := restreq.New("http://example.com").
		WithBodyReader().
		Post()

	if err == nil {
		defer resp.Response.Body.Close()
		b := bytes.NewBuffer([]byte{})
		b.ReadFrom(resp.Response.Body)
		fmt.Println(b.String())
	}

- Get header value

	value := resp.Header("token")

- Decode JSON

	s := struct {
		Message string `json:"message,omitempty"`
	}{}

	err := resp.DecodeJSON(&s)
*/
package restreq
