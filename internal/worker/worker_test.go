package worker_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"strings"
	"sync"

	gs "github.com/dragonheim/gagent/internal/gstructs"
	"github.com/dragonheim/gagent/internal/worker"
)

func TestWorkerMain(t *testing.T) {
	config := gs.GagentConfig{
		Name:       "test-config",
		Mode:       "worker",
		UUID:       "test-uuid",
		ListenAddr: "127.0.0.1",
		ClientPort: 1234,
		RouterPort: 5678,
		WorkerPort: 9012,
		Routers: []*gs.RouterDetails{
			{
				RouterName: "test-router",
				RouterID:   "test-router-id",
				RouterAddr: "127.0.0.1",
				WorkerPort: 9012,
			},
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	capturedOutput := captureOutput(func() {
		worker.Main(wg, config)
	})

	expectedOutput := `Starting worker`
	if !strings.Contains(capturedOutput, expectedOutput) {
		t.Errorf("Expected output to contain '%s', got '%s'", expectedOutput, capturedOutput)
	}

	wg.Wait()
}

func captureOutput(f func()) string {
	original := log.Writer()
	r, w, _ := os.Pipe()
	log.SetOutput(w)

	f()

	w.Close()
	log.SetOutput(original)

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}
