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

// NewPoolClient returns a *PoolClient that wraps an *http.Client and sets the
// maximum pool size as well as the requests per second as integers.
func NewPoolClient(stdClient http.Client, maxPoolSize int, reqPerSec int, timeout time.Duration) *PoolClient {
	var semaphore chan int = nil
	if maxPoolSize > 0 {
		semaphore = make(chan int, maxPoolSize) // Buffered channel to act as a semaphore
	}

	var emitter <-chan time.Time = nil
	if reqPerSec > 0 {
		emitter = time.NewTicker(time.Second / time.Duration(reqPerSec)).C // x req/s == 1s/x req (inverse)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 5
	transport.MaxIdleConnsPerHost = 5
	transport.MaxConnsPerHost = 10

	stdClient.Transport = transport
	stdClient.Timeout = timeout

	return &PoolClient{
		client:       stdClient,
		maxPoolSize:  maxPoolSize,
		cSemaphore:   semaphore,
		reqPerSecond: reqPerSec,
		rateLimiter:  emitter,
	}
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
