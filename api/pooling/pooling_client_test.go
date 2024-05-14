package pooling

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func init() {
	http.HandleFunc("/",
		func(respOut http.ResponseWriter, reqIn *http.Request) {
			defer reqIn.Body.Close()
		},
	)

	go func() {
		log.Fatalf("%s\n", http.ListenAndServe(":8080", nil))
	}()
}

func doPoolTest(client Client) error {
	testURL, err := url.Parse("http://127.0.0.1:8080/")
	if err != nil {
		return err
	}

	var doFunc func(req *http.Request) (*http.Response, error)

	switch cp := client.(type) {
	case *PoolClient:
		doFunc = cp.DoPool
	case *http.Client:
		doFunc = cp.Do
	}

	resp, err := doFunc(&http.Request{URL: testURL})
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return nil
}

// TestNewPClient tests creating a PClient
func TestNewPClient(t *testing.T) {
	t.Logf("[INFO] Starting TestNewPClient")

	defaultClient := http.Client{}

	_ = NewPoolClient(defaultClient, 25, 200, 10*time.Second) // Max 25 connections, 200 requests-per-second
	_ = NewPoolClient(defaultClient, 0, 0, 10*time.Second)    // Why do this? Just use http.Client
	_ = NewPoolClient(defaultClient, -1, -1, 10*time.Second)  // What

	t.Logf("[INFO] Completed TestNewPClient")
}

// TestPClient_Do tests performing a drop-in http.Client with pooling
func TestPClient_Do(t *testing.T) {
	t.Logf("[INFO] Starting TestPClient_Do")
	pClient := NewPoolClient(http.Client{}, 0, 0, time.Minute)

	if err := doPoolTest(pClient); err != nil {
		t.Error("pool: ", err)
	}
	t.Logf("[INFO] Completed TestClient_Do")
}

// TestPClient_DoPool tests performing a request with the pooling logic
func TestPClient_DoPool(t *testing.T) {
	t.Logf("[INFO] Starting TestPClient_DoPool")
	client := NewPoolClient(http.Client{}, 0, 0, time.Second)

	if err := doPoolTest(client); err != nil {
		t.Error("pool: ", err)
	}

	client2 := NewPoolClient(http.Client{}, 25, 200, 10*time.Second)

	if err := doPoolTest(client2); err != nil {
		t.Error("pool: ", err)
	}

	client3 := NewPoolClient(http.Client{}, -1, -1, 20*time.Second)

	if err := doPoolTest(client3); err != nil {
		t.Error("pool: ", err)
	}
	t.Logf("[INFO] Completed TestPClient_DoPool")
}
