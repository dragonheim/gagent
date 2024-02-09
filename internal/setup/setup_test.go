package setup_test

import (
	bytes "bytes"
	io "io"
	log "log"
	os "os"
	testing "testing"

	strings "strings"
	sync "sync"

	gs "github.com/dragonheim/gagent/internal/gstructs"
	setup "github.com/dragonheim/gagent/internal/setup"
)

func TestSetupMain(t *testing.T) {
	config := gs.GagentConfig{
		Name:       "test-config",
		Mode:       "client",
		UUID:       "test-uuid",
		ListenAddr: "127.0.0.1",
		ClientPort: 1234,
		RouterPort: 5678,
		WorkerPort: 9012,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	capturedOutput := captureOutput(func() {
		setup.Main(wg, config)
	})

	expectedOutput := `Configuration file created`
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
