package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	env "github.com/caarlos0/env/v6"
	gstructs "github.com/dragonheim/gagent/internal/gstructs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the WaitGroup to avoid actually waiting in tests
type MockWaitGroup struct {
	mock.Mock
}

func (m *MockWaitGroup) Add(delta int) {
	m.Called(delta)
}

func (m *MockWaitGroup) Done() {
	m.Called()
}

func (m *MockWaitGroup) Wait() {
	m.Called()
}

// Mocking the config loader function to inject test configurations
func mockInitConfig() {
	config = gstructs.GagentConfig{
		Mode:        "client",
		MonitorPort: 8080,
		// Populate other required fields as needed
	}
}

func TestMainFunction(t *testing.T) {
	var wg MockWaitGroup
	wg.On("Add", 1).Return()
	wg.On("Wait").Return()
	wg.On("Done").Return()

	mockInitConfig()

	// Test the client mode
	config.Mode = "client"
	main()
	wg.AssertCalled(t, "Add", 1)

	// Test the router mode
	config.Mode = "router"
	main()
	wg.AssertCalled(t, "Add", 1)

	// Test the worker mode
	config.Mode = "worker"
	main()
	wg.AssertCalled(t, "Add", 1)

	// Test the setup mode
	config.Mode = "setup"
	main()
	wg.AssertCalled(t, "Add", 1)

	// Test an invalid mode
	config.Mode = "invalid"
	assert.Panics(t, func() { main() }, "Expected main() to panic with invalid mode")
}

func TestInitFunction(t *testing.T) {
	// Backup original stdout and defer restoration
	origStdout := os.Stdout
	defer func() { os.Stdout = origStdout }()

	// Capture stdout output to test log output
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	// Test init
	init()

	// Assertions
	assert.Contains(t, logOutput.String(), "[DEBUG] Arguments are")
	assert.NotEmpty(t, config.Version, "Config version should not be empty")
	assert.NotEmpty(t, config.UUID, "Config UUID should not be empty")
}

func TestPrometheusMetricsExporter(t *testing.T) {
	mockInitConfig()
	config.MonitorPort = 8080

	req, err := http.NewRequest("GET", "/metrics", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")
	assert.Contains(t, rr.Body.String(), "go_gc_duration_seconds", "Expected metrics output")
}

func TestEnvironmentParsing(t *testing.T) {
	cfg := environment
	err := env.Parse(&cfg)
	assert.NoError(t, err)

	assert.Equal(t, "/etc/gagent/gagent.hcl", cfg.ConfigFile, "Expected default config file path")
	assert.Equal(t, "WARN", cfg.LogLevel, "Expected default log level")
	assert.Equal(t, "setup", cfg.Mode, "Expected default mode")
	assert.Equal(t, 0, cfg.MonitorPort, "Expected default monitor port")
}
