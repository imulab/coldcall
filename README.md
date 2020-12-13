# Cold Call
Utility to reduce boilerplate of making HTTP API calls in Go, includes only what I need.

[![Go Report Card](https://goreportcard.com/badge/github.com/imulab/coldcall)](https://goreportcard.com/report/github.com/imulab/coldcall)
[![Version](https://img.shields.io/badge/version-0.1.0-blue)](https://img.shields.io/badge/version-0.1.0-blue)

![cold-call](./assets/cold-call.png)

## Motivation

> Go's HTTP library is awesome, except being somewhat verbose at times. I wanted a thin wrapper around it to expose a succinct API to make it easy to create `http.Request` and read `http.Response`, and still being able to do complex stuff under the same set of API.
> 
> Hence, this library provides extensible functional APIs and out-of-box implementations for what I need for now.

## Install

```bash
go get github.com/imulab/coldcall
```

## Usage

To create a request

```go
req, err := coldcall.Post(context.Background(), "http://remote.com",
    addr.WithQueryMap(map[string]string{
        "foo": "bar",
    }),
	header.ContentType(header.ContentTypeApplicationJSON),
    body.JSONMarshal(Greeting{
	    Message: "hello world"
    }),
)
```

coldcall returns `*http.Request`, so you can execute it with your `*http.Client`.

To read a request

```go
type Data struct { Message string `json:message` }
func newData coldcall.Constructor = func() interface{} { return new(Data) }

// coldcall.Response creates a builder object on which you can register parsing rules.
resp := coldcall.Response(http.DefaultClient.Do(req))

// Use Expect to register parsing rules. Here, we say: "when status is 200, marshal the body as JSON into Data". Chain multiple Expect rules together!
//
// When done, call Read to get results. "data" will be *Data if status is indeed 200, otherwise it will be
// []byte, "raw" is always []byte in case you need the original data, "err" is any error encountered.
data, raw, err := resp.Expect(status.Is200, body.JSONUnmarshal(newData)).Read()

// If you want the original http.Response object
httpResp, err := resp.Original(), resp.Error()
```

## Main Concepts

### `Options`

A function to customize an aspect of a `http.Request`. Currently, out-of-box options include:
- `addr.WithQuery`
- `body.JSONMarshal`
- `body.XMLMarshal`
- `body.URLValuesEncode`
- `header.Custom`
- `header.ContentType`
- `header.Accept`
- and their convenience wrappers

Didn't include what you want? Implement `Options` directly to roll your own!

### `Condition` and `Producer`

`Condition` is a function that tells coldcall when to read the `http.Response` using the corresponding criteria. Conditions are
registered using the `Expect` function and are evaluated in order of registration. 

Currently, out-of-box `Condition` include:
- `status.Is`
- `status.Is200`
- `status.InRange` 
- `status.IsSuccess`
- `status.IsFailure`

**Try to avoid condition overlaps!** The first matching condition always wins.

`Producer` is a function that converts the response body bytes into an object. It is registered along with `Condition` as the actual
processing logic. 

Feel free to implement this, however, you may not need to. See `Unmarshaler` and `Constructor`.

### `Unmarshaller` and `Constructor`

`Unmarshaller` is a function to unmarshal bytes into an object. It has the same function signature with `json.Unmarshal`
and `xml.Unmarshal` so they can be used directly.

Currently, out-of-box `Unmarshaller` includes:
- `body.JSONUnmarshal`
- `body.XMLUnmarshal`

The unmarshalling process requires a made object, and here's where the `Constructor` comes in. The `Constructor` is a function
that returns a new object to be unmarshalled.

Together, `Unmarshaller` and `Constructor` makes up a `Producer`. And in most cases, you may only need to write a `Constructor`.

## Maintenance and Contribution

This library is maintained by me by myself. Updates to it are driven by real needs from my projects.

Contributions to `Options` implementation for common use cases in `header`, `body` and `addr` are welcomed. 

Before opening PRs, Please file an issue to discuss it first.