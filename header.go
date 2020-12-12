package codecall

import (
	"net/http"
)

const (
	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"
	HeaderAccept        = "Accept"
)

const (
	ContentTypeApplicationJSON           = "application/json"
	ContentTypeApplicationXML            = "application/xml"
	ContentTypeApplicationFormUrlEncoded = "application/x-www-form-urlencoded"
)

// ContentType sets the given contentType as the "Content-Type" header on the http.Request.
func ContentType(contentType string) Option {
	return Header(HeaderContentType, contentType)
}

// Accept sets the given contentType as the "Accept" header on the http.Request.
func Accept(contentType string) Option {
	return Header(HeaderAccept, contentType)
}

// Header sets the given name and value as header on the http.Request. Multiple calls with
// the same name will append the value.
func Header(name string, value string) Option {
	return &kvHeaderOption{key: name, val: value}
}

type kvHeaderOption struct {
	key string
	val string
}

func (opt *kvHeaderOption) Apply(r *http.Request) error {
	r.Header.Add(opt.key, opt.val)
	return nil
}
