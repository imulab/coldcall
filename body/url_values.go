package body

import (
	"github.com/imulab/coldcall"
	"net/url"
	"strings"
)

// URLValuesEncode option encodes the given url values and sets it as the body on http.Request.
//
// Although typically used in the "application/x-www-form-urlencoded" scenario, this function
// does not automatically set the "Content-Type" header.
func URLValuesEncode(values url.Values) coldcall.Option {
	return Read(strings.NewReader(values.Encode()))
}
