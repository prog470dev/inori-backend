package base

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func defaultServer() *Server {
	s := New()
	s.Init("../dbconfig.yaml") // testの場合は対象ファイルからの相対パス
	return s
}

func TestHealthCheck(t *testing.T) {
	s := defaultServer()

	ts := httptest.NewServer(s.router)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Fatalf("ping failed %s", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if string(b) != "pong" {
		t.Fatalf("expected: %s, actual: %s", "pong", b)
	}
	defer resp.Body.Close()
}
