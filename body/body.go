package body

import (
	"bytes"
	"github.com/imulab/coldcall"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Read option sets the io.Reader body on the http.Request.
func Read(body io.Reader) coldcall.Option {
	return func(req *http.Request) error {
		rc, ok := body.(io.ReadCloser)
		if !ok && body != nil {
			rc = ioutil.NopCloser(body)
		}

		req.Body = rc

		if body != nil {
			switch v := body.(type) {
			case *bytes.Buffer:
				req.ContentLength = int64(v.Len())
				buf := v.Bytes()
				req.GetBody = func() (io.ReadCloser, error) {
					r := bytes.NewReader(buf)
					return ioutil.NopCloser(r), nil
				}
			case *bytes.Reader:
				req.ContentLength = int64(v.Len())
				snapshot := *v
				req.GetBody = func() (io.ReadCloser, error) {
					r := snapshot
					return ioutil.NopCloser(&r), nil
				}
			case *strings.Reader:
				req.ContentLength = int64(v.Len())
				snapshot := *v
				req.GetBody = func() (io.ReadCloser, error) {
					r := snapshot
					return ioutil.NopCloser(&r), nil
				}
			default:
				// This is where we'd set it to -1 (at least
				// if body != NoBody) to mean unknown, but
				// that broke people during the Go 1.8 testing
				// period. People depend on it being 0 I
				// guess. Maybe retry later. See Issue 18117.
			}
			// For client requests, Request.ContentLength of 0
			// means either actually 0, or unknown. The only way
			// to explicitly say that the ContentLength is zero is
			// to set the Body to nil. But turns out too much code
			// depends on NewRequest returning a non-nil Body,
			// so we use a well-known ReadCloser variable instead
			// and have the http package also treat that sentinel
			// variable to mean explicitly zero.
			if req.GetBody != nil && req.ContentLength == 0 {
				req.Body = http.NoBody
				req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
			}
		}

		return nil
	}
}

// Marshal option encodes the body using the supplied coldcall.Marshaller and sets the rendered
// raw data on http.Request as body. The supplied body must not be nil, otherwise coldcall.ErrNoBody is
// returned. If the coldcall.Marshaller returned an error, coldcall.ErrBadBody is returned.
func Marshal(body interface{}, marshaller coldcall.Marshaller) coldcall.Option {
	return func(r *http.Request) error {
		if body == nil {
			return coldcall.ErrNoBody
		}

		raw, err := marshaller(body)
		if err != nil {
			return coldcall.ErrBadBody
		}

		return Read(bytes.NewReader(raw))(r)
	}
}
