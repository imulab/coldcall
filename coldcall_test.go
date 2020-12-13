package coldcall_test

import (
	"context"
	"github.com/imulab/coldcall"
	"github.com/imulab/coldcall/addr"
	"github.com/imulab/coldcall/body"
	"github.com/imulab/coldcall/header"
	"github.com/imulab/coldcall/status"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	type (
		Greeting struct {
			Message string `json:"message"`
		}
		Echo struct {
			JSON Greeting `json:"json"`
		}
	)

	req, err := coldcall.Post(context.Background(), "http://httpbin.org/post",
		addr.WithQueryMap(map[string]string{
			"foo": "bar",
		}),
		header.ContentType(header.ContentTypeApplicationJSON),
		body.JSONMarshal(Greeting{Message: "hello world"}),
	)
	if err != nil {
		t.Error(err)
	}

	println(req.URL.String())

	var newEcho coldcall.Constructor = func() interface{} {
		return new(Echo)
	}

	v, raw, err := coldcall.Response(http.DefaultClient.Do(req)).
		Expect(status.Is200, body.JSONUnmarshal(newEcho)).
		Read()
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s\n", string(raw))
	t.Logf("echoed greeting is '%s'\n", v.(*Echo).JSON.Message)
}
