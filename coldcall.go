package coldcall

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
)

type (
	// Option applies modification to the http.Request before it hits the wire.
	Option func(r *http.Request) error

	// Marshaller encodes object into raw data. The function signature is designed
	// to match json.Marshal and xml.Marshal so that they can be used directly.
	Marshaller func(v interface{}) ([]byte, error)

	// Condition returns true if the given http.Response meets its criteria. It is
	// used in conjunction with Producer to determine when to invoke the Producer
	// to read the response data.
	Condition func(*http.Response) bool

	// Producer can produce an object from the raw data. It is normally not used
	// directly, but rather composed by Constructor and Unmarshaler.
	//
	// See stock producers like body.JSONUnmarshal and body.XMLUnmarshal
	Producer func(raw []byte) (interface{}, error)

	// Unmarshaler decodes raw data into the given object, the object is normally
	// produced by Constructor. The function signature is designed to match json.Unmarshal
	// and xml.Unmarshal, so that they can be used directly.
	Unmarshaler func(raw []byte, v interface{}) error

	// Constructor returns a new object, ready to be used by Unmarshaler. Normally,
	// this is the only implementation that the user needs to provide.
	Constructor func() interface{}
)

// Request constructs the http.Request object with the method and given options.
func Request(ctx context.Context, method string, url string, options ...Option) (*http.Request, error) {
	r, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range options {
		if err := opt(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

// Get is a convenient wrapper for Request with the GET method.
func Get(ctx context.Context, url string, options ...Option) (*http.Request, error) {
	return Request(ctx, http.MethodGet, url, options...)
}

// Post is a convenient wrapper for Request with the POST method.
func Post(ctx context.Context, url string, options ...Option) (*http.Request, error) {
	return Request(ctx, http.MethodPost, url, options...)
}

// Put is a convenient wrapper for Request with the PUT method.
func Put(ctx context.Context, url string, options ...Option) (*http.Request, error) {
	return Request(ctx, http.MethodPut, url, options...)
}

// Patch is a convenient wrapper for Request with the PATCH method.
func Patch(ctx context.Context, url string, options ...Option) (*http.Request, error) {
	return Request(ctx, http.MethodPatch, url, options...)
}

// Delete is a convenient wrapper for Request with the DELETE method.
func Delete(ctx context.Context, url string, options ...Option) (*http.Request, error) {
	return Request(ctx, http.MethodDelete, url, options...)
}

// Response reads the http.Response and returns a builder object for further configuration.
// The function signature is designed so that it can be applied to the result of http.Client.Do
// directly.
func Response(response *http.Response, err error) *reader {
	return &reader{
		resp: response,
		err:  err,
	}
}

type reader struct {
	resp      *http.Response
	err       error
	criterion []produceOnCondition
	fallback  *produceOnCondition
}

// Original returns the original http.Response object.
func (r *reader) Original() *http.Response {
	return r.resp
}

// Error returns any error during the process.
func (r *reader) Error() error {
	return r.err
}

// Expect tells the builder that what type of response object (i.e. Producer) to
// parse the body with and when (i.e. Condition) to do that.
//
// The conditions and producers are applied in order. If conditions overlap, the
// first matching rule will be applied. If no condition matches the request, the
// raw body is returned.
func (r *reader) Expect(condition Condition, producer Producer) *reader {
	if r.criterion == nil {
		r.criterion = []produceOnCondition{}
	}

	r.criterion = append(r.criterion, produceOnCondition{
		condition: condition,
		producer:  producer,
	})

	return r
}

// Read executes the registered rules in sequence and return the
// parsed body, raw body and any error.
func (r *reader) Read() (interface{}, []byte, error) {
	if r.err != nil {
		return nil, nil, r.err
	}

	raw, err := ioutil.ReadAll(r.resp.Body)
	if err != nil {
		return nil, raw, err
	}

	for _, criteria := range r.criterion {
		if ok := criteria.condition(r.resp); !ok {
			continue
		}

		v, err := criteria.producer(raw)
		if err != nil {
			return nil, raw, err
		}

		return v, raw, nil
	}

	return raw, raw, nil
}

// Produce constructs a Producer with Constructor and Unmarshaler.
func Produce(constructor Constructor, unmarshaler Unmarshaler) Producer {
	return func(raw []byte) (interface{}, error) {
		v := constructor()
		err := unmarshaler(raw, v)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

type produceOnCondition struct {
	condition Condition
	producer  Producer
}

// URLValues is a convenient function to construct url.Values with a map
func URLValues(kv map[string]string) url.Values {
	values := url.Values{}
	for k, v := range kv {
		values.Set(k, v)
	}
	return values
}
