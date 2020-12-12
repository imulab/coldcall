package addr

import (
	"fmt"
	"github.com/imulab/coldcall"
	"net/http"
	"net/url"
	"strings"
)

// URL returns a coldcall.Option that sets the request address.
func URL(u *url.URL) coldcall.Option {
	return func(r *http.Request) error {
		r.URL = u
		removeEmptyPort(u.Host)
		r.Host = u.Host
		return nil
	}
}

// String is a convenient wrapper for URL that uses plain string as address.
func String(addr string) coldcall.Option {
	return func(r *http.Request) error {
		u, err := url.Parse(addr)
		if err != nil {
			return err
		}
		return URL(u)(r)
	}
}

// SPrintf is a convenient wrapper for String that allows users to render url with fmt.Sprintf.
func SPrintf(tmpl string, args ...interface{}) coldcall.Option {
	return String(fmt.Sprintf(tmpl, args...))
}

// WithQuery is a convenient wrapper for URL that sets the params as the raw query component.
func WithQuery(base string, params url.Values) coldcall.Option {
	return func(r *http.Request) error {
		u, err := url.Parse(base)
		if err != nil {
			return err
		}

		u.RawQuery = params.Encode()
		return URL(u)(r)
	}
}

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

// removeEmptyPort strips the empty port in ":port" to ""
// as mandated by RFC 3986 Section 6.2.3.
func removeEmptyPort(host string) string {
	if hasPort(host) {
		return strings.TrimSuffix(host, ":")
	}
	return host
}
