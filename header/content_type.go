package header

import "github.com/imulab/coldcall"

const (
	KeyContentType = "Content-Type"

	ContentTypeApplicationJSON           = "application/json"
	ContentTypeApplicationXML            = "application/xml"
	ContentTypeApplicationFormUrlEncoded = "application/x-www-form-urlencoded"
)

// ContentType sets the given contentType as the KeyContentType header on the http.Request.
func ContentType(contentType string) coldcall.Option {
	return Custom(KeyContentType, contentType)
}
