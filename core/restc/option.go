package restc

import (
	"net/http"
	"time"
)

func WithHeaders(headers http.Header) Opt {
	return func(c *client) error {
		c.headers = headers
		return nil
	}
}

func WithRetryTimes(times int) Opt {
	return func(c *client) error {
		c.retryTimes = times
		return nil
	}
}

func WithRetryDelay(time time.Duration) Opt {
	return func(c *client) error {
		c.retryDelay = time
		return nil
	}
}

func WithClient(c *http.Client) Opt {
	return func(client *client) error {
		client.client = c
		return nil
	}
}

func WithRequestMiddleware(middleware RequestMiddleware) Opt {
	return func(c *client) error {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.beforeRequest = append(c.beforeRequest, middleware)
		return nil
	}
}
