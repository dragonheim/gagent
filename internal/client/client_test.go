package client

import (
	bytes "bytes"
	errors "errors"
	io "io"
	log "log"
	os "os"
	sync "sync"
	testing "testing"

	gs "github.com/dragonheim/gagent/internal/gstructs"

	zmq "github.com/pebbe/zmq4"
)

type mockSocket struct {
	sendMessageError error
}

func (m *mockSocket) Close() error                                               { return nil }
func (m *mockSocket) Bind(endpoint string) error                                 { return nil }
func (m *mockSocket) Connect(endpoint string) error                              { return nil }
func (m *mockSocket) SetIdentity(identity string) error                          { return nil }
func (m *mockSocket) SendMessage(parts ...interface{}) (int, error)              { return 0, m.sendMessageError }
func (m *mockSocket) RecvMessage(flags zmq.Flag) ([]string, error)               { return nil, nil }
func (m *mockSocket) RecvMessageBytes(flags zmq.Flag) ([][]byte, error)          { return nil, nil }
func (m *mockSocket) RecvMessageString(flags zmq.Flag) ([]string, error)         { return nil, nil }
func (m *mockSocket) SetSubscribe(filter string) error                           { return nil }
func (m *mockSocket) SetUnsubscribe(filter string) error                         { return nil }
func (m *mockSocket) Send(msg string, flags zmq.Flag) (int, error)               { return 0, nil }
func (m *mockSocket) SendBytes(msg []byte, flags zmq.Flag) (int, error)          { return 0, nil }
func (m *mockSocket) SendFrame(msg []byte, flags zmq.Flag) (int, error)          { return 0, nil }
func (m *mockSocket) SendMultipart(parts [][]byte, flags zmq.Flag) (int, error)  { return 0, nil }
func (m *mockSocket) Recv(flags zmq.Flag) (string, error)                        { return "", nil }
func (m *mockSocket) RecvBytes(flags zmq.Flag) ([]byte, error)                   { return nil, nil }
func (m *mockSocket) RecvFrame(flags zmq.Flag) ([]byte, error)                   { return nil, nil }
func (m *mockSocket) RecvMultipart(flags zmq.Flag) ([][]byte, error)             { return nil, nil }
func (m *mockSocket) SetOption(option zmq.SocketOption, value interface{}) error { return nil }
func (m *mockSocket) GetOption(option zmq.SocketOption) (interface{}, error)     { return nil, nil }
func (m *mockSocket) Events() zmq.State                                          { return 0 }
func (m *mockSocket) String() string                                             { return "" }

func TestGetTagsFromHints(t *testing.T) {
	agent := gs.AgentDetails{
		Script: []byte(`*set GHINT[split "tag1,tag2,tag3",]`),
	}

	expectedHints := []string{"tag1", "tag2", "tag3"}
	hints := getTagsFromHints(agent)

	if !equalStringSlices(hints, expectedHints) {
		t.Errorf("Expected hints %v, but got %v", expectedHints, hints)
	}
}

func TestSendAgent(t *testing.T) {
	wg := &sync.WaitGroup{}
	config := gs.GagentConfig{
		UUID:       "test-uuid",
		ClientPort: 1234,
		Routers: map[string]gs.Router{
			"test-router": {
				RouterAddr: "127.0.0.1",
				ClientPort: 1234,
			},
		},
	}

	agent := gs.AgentDetails{
		Client: "test-client",
		Script: []byte(`*set GHINT[split "tag1,tag2,tag3",]`),
	}

	// Replace zmq.NewSocket with a function that returns a mock socket
	origNewSocket := newSocket
	defer func() { newSocket = origNewSocket }()
	newSocket = func(t zmq.Type) (zmq.Socket, error) {
		return &mockSocket{}, nil
	}

	wg.Add(1)
	go sendAgent(wg, config.UUID, "tcp://127.0.0.1:1234", agent)
	wg.Wait()

	// Test with an error in sending a message
	newSocket = func(t zmq.Type) (zmq.Socket, error) {
		return &mockSocket{sendMessageError: errors.New("send message error")}, nil
	}

	wg.Add(1)
	go sendAgent(wg, config.UUID, "tcp://127.0.0.1:1234", agent)
	wg.Wait()
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestMain(t *testing.T) {
	// Prepare a temporary agent file for testing
	tmpAgentFile, err := io.TempFile("", "agent")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpAgentFile.Name())

	content := []byte(`*set GHINT[split "tag1,tag2,tag3",]`)
	if _, err := tmpAgentFile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpAgentFile.Close(); err != nil {
		t.Fatal(err)
	}

	config := gs.GagentConfig{
		CMode:      true,
		UUID:       "test-uuid",
		ClientPort: 1234,
		Agent:      tmpAgentFile.Name(),
		Routers: map[string]gs.Router{
			"test-router": {
				RouterAddr: "127.0.0.1",
				ClientPort: 1234,
			},
		},
	}

	// Replace log output with a buffer to suppress output during testing
	origLogOutput := log.Writer()
	defer log.SetOutput(origLogOutput)
	log.SetOutput(&bytes.Buffer{})

	// Replace zmq.NewSocket with a function that returns a mock socket
	origNewSocket := newSocket
	defer func() { newSocket = origNewSocket }()
	newSocket = func(t zmq.Type) (zmq.Socket, error) {
		return &mockSocket{}, nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go Main(wg, config)
	wg.Wait()
}
