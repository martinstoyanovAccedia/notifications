package pooling

import (
	"net/http"
	"time"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type PoolClient struct {
	client       http.Client
	maxPoolSize  int
	cSemaphore   chan int
	reqPerSecond int
	rateLimiter  <-chan time.Time
}

// Do 'overloads' the http.Client Do method
func (c *PoolClient) Do(req *http.Request) (*http.Response, error) {
	return c.DoPool(req)
}

// DoPool does the required synchronization with channels to
// perform the request in accordance to the PoolClient's maximum pool size
func (c *PoolClient) DoPool(req *http.Request) (*http.Response, error) {
	if c.maxPoolSize > 0 {
		c.cSemaphore <- 1 // connection from our pool
		defer func() {
			<-c.cSemaphore // release connection back to the pool
		}()
	}

	if c.reqPerSecond > 0 {
		<-c.rateLimiter // wait for the client to emmit
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	// do not close resp.Body here

	return resp, nil
}
