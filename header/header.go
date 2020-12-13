package header

import (
	"github.com/imulab/coldcall"
	"net/http"
)

// Custom sets the given name and value as header on the http.Request. Multiple calls with
// the same name will append the value.
func Custom(name string, value string) coldcall.Option {
	return func(r *http.Request) error {
		r.Header.Add(name, value)
		return nil
	}
}
