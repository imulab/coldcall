package body

import (
	"encoding/xml"
	"github.com/imulab/coldcall"
)

// XMLMarshal option marshals the given body into XML and sets it
// as the body on http.Request.
//
// The supplied body parameter must not be nil, otherwise ErrNoBody
// error is returned. In addition, the body parameter must be capable of
// marshalling into XML, otherwise, ErrBadBody error is returned.
//
// Although typically used in a "application/xml" content type scenario,
// this function does not automatically set the "Content-Type" header.
func XMLMarshal(body interface{}) coldcall.Option {
	return Marshal(body, xml.Marshal)
}

// XMLUnmarshal returns a coldcall.Producer which can unmarshal raw data into
// the object returned by constructor in XML format.
func XMLUnmarshal(constructor coldcall.Constructor) coldcall.Producer {
	return coldcall.Produce(constructor, xml.Unmarshal)
}
