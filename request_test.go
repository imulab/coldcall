package codecall

import (
	"context"
	"testing"
)

func TestRequest(t *testing.T) {
	_, _ = Post(context.Background(), "/httpbin",
		ContentType(ContentTypeApplicationJSON),
		JSONBody(struct {
			Greeting string `json:"greeting"`
		}{
			Greeting: "hello world",
		}),
	)
}
