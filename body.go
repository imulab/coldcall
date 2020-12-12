package codecall

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// JSONBody option marshals the given body into JSON and sets it
// as the body on http.Request.
//
// The supplied body parameter must not be nil, otherwise ErrNoBody
// error is returned. In addition, the body parameter must be capable of
// marshalling into JSON, otherwise, ErrBadBody error is returned.
//
// Although typically used in a "application/json" content type scenario,
// this function does not automatically set the "Content-Type" header.
func JSONBody(body interface{}) Option {
	return &marshallerOption{
		body:       body,
		marshaller: json.Marshal,
	}
}

// XMLBody option marshals the given body into XML and sets it
// as the body on http.Request.
//
// The supplied body parameter must not be nil, otherwise ErrNoBody
// error is returned. In addition, the body parameter must be capable of
// marshalling into XML, otherwise, ErrBadBody error is returned.
//
// Although typically used in a "application/xml" content type scenario,
// this function does not automatically set the "Content-Type" header.
func XMLBody(body interface{}) Option {
	return &marshallerOption{
		body:       body,
		marshaller: xml.Marshal,
	}
}

type marshallerOption struct {
	body       interface{}
	marshaller func(v interface{}) ([]byte, error)
}

func (opt *marshallerOption) Apply(r *http.Request) error {
	if opt.body == nil {
		return ErrNoBody
	}

	raw, err := opt.marshaller(opt.body)
	if err != nil {
		return ErrBadBody
	}

	return ReaderBody(bytes.NewReader(raw)).Apply(r)
}

// URLValuesBody option encodes the given url values and sets it as the body on http.Request.
//
// Although typically used in the "application/x-www-form-urlencoded" scenario, this function
// does not automatically set the "Content-Type" header.
func URLValuesBody(values url.Values) Option {
	return ReaderBody(strings.NewReader(values.Encode()))
}

// ReaderBody option sets the io.Reader body on the http.Request.
func ReaderBody(body io.Reader) Option {
	return &readerBodyOption{body: body}
}

type readerBodyOption struct {
	body io.Reader
}

func (opt *readerBodyOption) Apply(req *http.Request) error {
	body := opt.body

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
