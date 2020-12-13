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

func ExampleRequest_get() {
	type Data struct {
		Message string `json:"message"`
	}

	var dataConstructor coldcall.Constructor = func() interface{} {
		return new(Data)
	}

	req, _ := coldcall.Get(context.Background(), "http://remote.com",
		addr.WithQueryMap(map[string]string{
			"foo": "bar",
		}),
	)

	data, raw, _ := coldcall.Response(http.DefaultClient.Do(req)).
		Expect(status.Is200, body.JSONUnmarshal(dataConstructor)).
		Read()

	println(data.(*Data).Message, raw)
}

func ExampleRequest_postJSON() {
	type (
		Greeting struct {
			Message string `json:"message"`
		}
		Echo struct {
			JSON Greeting `json:"json"`
		}
	)

	var echoConstructor coldcall.Constructor = func() interface{} {
		return new(Echo)
	}

	req, _ := coldcall.Post(context.Background(), "http://remote.com",
		header.ContentType(header.ContentTypeApplicationJSON),
		body.JSONMarshal(Greeting{Message: "hello world"}),
	)

	v, raw, _ := coldcall.Response(http.DefaultClient.Do(req)).
		Expect(status.Is200, body.JSONUnmarshal(echoConstructor)).
		Read()

	println(v.(*Echo).JSON.Message, raw)
}

func TestRequest_PostJSON(t *testing.T) {
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

func TestRequest_PostForm(t *testing.T) {
	type (
		Echo struct {
			Form map[string]string `json:"form"`
		}
	)

	var echoConstructor coldcall.Constructor = func() interface{} {
		return new(Echo)
	}

	req, err := coldcall.Post(context.Background(), "http://httpbin.org/post",
		header.ContentType(header.ContentTypeApplicationFormUrlEncoded),
		body.URLValuesMapEncode(map[string]string{
			"foo": "bar",
		}),
	)
	if err != nil {
		t.Error(err)
	}

	echo, _, err := coldcall.Response(http.DefaultClient.Do(req)).
		Expect(status.Is200, body.JSONUnmarshal(echoConstructor)).
		Read()
	if err != nil {
		t.Error(err)
	}

	if "bar" != echo.(*Echo).Form["foo"] {
		t.FailNow()
	}
}
