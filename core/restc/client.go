package restc

import (
	"net/http"
	"sync"
	"time"
)

type Client interface {
	Verb(verb string) *Request
	SetHeader(headers http.Header)
}

type Opt func(client *client) error

type client struct {
	lock *sync.RWMutex
	addr string

	retryTimes int
	retryDelay time.Duration

	headers http.Header

	// Set specific behavior of the client.  If not set http.DefaultClient will be used.
	client *http.Client

	// middleware
	beforeRequest []RequestMiddleware
}

func (c *client) SetHeader(headers http.Header) {
	c.headers = headers
}

type RequestMiddleware func(Client, *Request) error

func (c *client) requestMiddlewares() []RequestMiddleware {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.beforeRequest
}

func (c *client) executeRequestMiddlewares(req *Request) (err error) {
	for _, f := range c.requestMiddlewares() {
		if err = f(c, req); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) Verb(verb string) *Request {
	return NewRequest(c).Verb(verb)
}

func NewClient(addr string, ops ...Opt) (Client, error) {
	c := &client{
		lock: &sync.RWMutex{},
		addr: addr,
	}

	for _, op := range ops {
		if err := op(c); err != nil {
			return nil, err
		}
	}

	if c.client == nil {
		c.client = &http.Client{}
	}
	if c.headers == nil {
		c.headers = make(http.Header)
	}
	return c, nil
}
