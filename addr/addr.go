package addr

import (
	"github.com/imulab/coldcall"
	"net/http"
	"net/url"
)

// WithQuery is a convenient wrapper for URL that sets the params as the raw query component.
func WithQuery(params url.Values) coldcall.Option {
	return func(r *http.Request) error {
		r.URL.RawQuery = params.Encode()
		return nil
	}
}

// WithQueryMap is a convenience wrapper for WithQuery with map values to construct url.Values.
func WithQueryMap(params map[string]string) coldcall.Option {
	return WithQuery(coldcall.URLValues(params))
}
