package body

import (
	"encoding/json"
	"github.com/imulab/coldcall"
)

// JSONMarshal option marshals the given body into JSON and sets it
// as the body on http.Request.
//
// The supplied body parameter must not be nil, otherwise ErrNoBody
// error is returned. In addition, the body parameter must be capable of
// marshalling into JSON, otherwise, ErrBadBody error is returned.
//
// Although typically used in a "application/json" content type scenario,
// this function does not automatically set the "Content-Type" header.
func JSONMarshal(body interface{}) coldcall.Option {
	return Marshal(body, json.Marshal)
}

// JSONUnmarshal returns a coldcall.Producer which can unmarshal the raw data into the object
// returned by the constructor in JSON format.
func JSONUnmarshal(constructor coldcall.Constructor) coldcall.Producer {
	return coldcall.Produce(constructor, json.Unmarshal)
}
