package codecall

import "net/http"

// Option applies modification to the http.Request before it hits the wire.
type Option interface {
	// Apply modifies the http.Request in place, returns any error encountered.
	Apply(r *http.Request) error
}
