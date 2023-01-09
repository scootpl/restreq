# RestReq

[![Go Reference](https://pkg.go.dev/badge/github.com/scootpl/restreq.svg)](https://pkg.go.dev/github.com/scootpl/restreq)  [![Go Report Card](https://goreportcard.com/badge/github.com/scootpl/restreq)](https://goreportcard.com/report/github.com/scootpl/restreq)

RestReq is a wrapper around standard Go net/http client. In a simple call you can use json encoding, add headers
and parse result. This should be sufficient in most use cases.

## Features

- Simple syntax
- Only stdlib (no external dependencies)
- JSON parsing
- Debug logging

## Quick Start

```go
import "github.com/scootpl/restreq"

resp, err := restreq.New("http://example.com").
	AddHeader("X-TOKEN", authToken).
	Post()
```

## Examples and Documentation

See [GoDoc](https://godoc.org/github.com/scootpl/restreq) for more details.
