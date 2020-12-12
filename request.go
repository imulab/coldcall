package codecall

import (
	"context"
	"net/http"
)

func Request(ctx context.Context, method string, url string, options ...Option) (*http.Request, error) {
	r, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range options {
		if err := opt.Apply(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func Get(ctx context.Context, url string, options ...Option) (*http.Request, error) {
	return Request(ctx, http.MethodGet, url, options...)
}

func Post(ctx context.Context, url string, options ...Option) (*http.Request, error) {
	return Request(ctx, http.MethodPost, url, options...)
}

func Put(ctx context.Context, url string, options ...Option) (*http.Request, error) {
	return Request(ctx, http.MethodPut, url, options...)
}

func Patch(ctx context.Context, url string, options ...Option) (*http.Request, error) {
	return Request(ctx, http.MethodPatch, url, options...)
}

func Delete(ctx context.Context, url string, options ...Option) (*http.Request, error) {
	return Request(ctx, http.MethodDelete, url, options...)
}
