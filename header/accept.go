package header

import "github.com/imulab/coldcall"

const (
	KeyAccept = "KeyAccept"
)

// Accept sets the given contentType as the "KeyAccept" header on the http.Request.
func Accept(contentType string) coldcall.Option {
	return Custom(KeyAccept, contentType)
}
