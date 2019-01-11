package runner

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type server struct{}

func (*server) Serve(l net.Listener) error {
	panic("implement me")
}

func (*server) Close() error {
	panic("implement me")
}

func TestServerRunner_Start_NoServer(t *testing.T) {
	runner := &ServerRunner{}

	err := runner.Start()
	if err == nil {
		t.Error("runner should not start without a server")
	}
}

func TestServerRunner_Start_NoListener(t *testing.T) {
	runner := &ServerRunner{
		Server: &server{},
	}

	err := runner.Start()
	if err == nil {
		t.Error("runner should not start without a listener")
	}
}

func TestServerRunner_Start(t *testing.T) {
	testServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))

	runner := NewServerRunner(testServer.Config, testServer.Listener)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := runner.Start()
		if err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	resp, err := http.Get("http://" + runner.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "ok" {
		t.Errorf("unexpected body: %s", string(body))
	}

	runner.Stop(nil)
	wg.Wait()
}
