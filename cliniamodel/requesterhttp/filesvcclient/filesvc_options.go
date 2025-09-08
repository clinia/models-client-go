package filesvcclient

import (
	"context"
	"net/http"
)

// WithAPIKey adds `Authorization : Api-Key <key>` to every request.
func WithAPIKey(key string) ClientOption {
	return func(c *Client) error {
		if key == "" {
			return nil
		}
		c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, r *http.Request) error {
			r.Header.Set("Authorization", "Api-Key "+key)
			return nil
		})
		return nil
	}
}

// WithBearer adds `Authorization: Bearer <token>` to every request.
func WithBearer(token string) ClientOption {
	return func(c *Client) error {
		if token == "" {
			return nil
		}
		c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, r *http.Request) error {
			r.Header.Set("Authorization", "Bearer "+token)
			return nil
		})
		return nil
	}
}

// WithHeader adds a static header to every request (useful for User-Agent, tracing, etc.).
func WithHeader(name, value string) ClientOption {
	return func(c *Client) error {
		if name == "" || value == "" {
			return nil
		}
		c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, r *http.Request) error {
			r.Header.Set(name, value)
			return nil
		})
		return nil
	}
}
