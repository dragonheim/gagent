package router_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	gs "github.com/dragonheim/gagent/internal/gstructs"
	"github.com/dragonheim/gagent/internal/router"
)

func TestRouterMain(t *testing.T) {
	config := gs.GagentConfig{
		Name:        "test-config",
		Mode:        "router",
		UUID:        "test-uuid",
		ListenAddr:  "127.0.0.1",
		ClientPort:  1234,
		RouterPort:  5678,
		WorkerPort:  9012,
		ChainDBPath: "test-chaindb-path",
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go router.Main(wg, config)

	// Allow router to start before sending HTTP requests
	time.Sleep(time.Millisecond * 100)

	// Test GET request
	resp := makeRequest(t, "GET", "http://localhost:1234/hello")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test POST request
	resp = makeRequest(t, "POST", "http://localhost:1234/hello")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test OPTIONS request
	resp = makeRequest(t, "OPTIONS", "http://localhost:1234/hello")
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	// Test unsupported method
	resp = makeRequest(t, "PUT", "http://localhost:1234/hello")
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}

	wg.Wait()
}

func makeRequest(t *testing.T, method, url string) *http.Response {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(router.AnswerClient)
	handler.ServeHTTP(rr, req)

	return rr.Result()
}
