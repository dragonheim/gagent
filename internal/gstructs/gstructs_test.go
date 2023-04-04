package gstructs_test

import (
	"testing"

	"github.com/dragonheim/gagent/internal/gstructs"
)

func TestGagentConfig(t *testing.T) {
	config := gstructs.GagentConfig{
		Name:        "test-config",
		Mode:        "client",
		UUID:        "test-uuid",
		ListenAddr:  "127.0.0.1",
		ChainDBPath: "/tmp/chaindb",
		MonitorPort: 8888,
		ClientPort:  1234,
		RouterPort:  5678,
		WorkerPort:  9012,
		Clients: []*gstructs.ClientDetails{
			{
				ClientName: "test-client",
				ClientID:   "client-id",
			},
		},
		Routers: []*gstructs.RouterDetails{
			{
				RouterName: "test-router",
				RouterID:   "router-id",
				RouterAddr: "192.168.1.1",
				RouterTags: []string{"tag1", "tag2"},
				ClientPort: 1234,
				RouterPort: 5678,
				WorkerPort: 9012,
			},
		},
		Workers: []*gstructs.WorkerDetails{
			{
				WorkerName: "test-worker",
				WorkerID:   "worker-id",
				WorkerTags: []string{"tag3", "tag4"},
			},
		},
		Version: "1.0.0",
		File:    "config.hcl",
		Agent:   "agent.gagent",
		CMode:   true,
	}

	if config.Name != "test-config" {
		t.Errorf("Expected config name to be 'test-config', got %s", config.Name)
	}
	if config.Mode != "client" {
		t.Errorf("Expected config mode to be 'client', got %s", config.Mode)
	}
	// TODO: add more assertions for other config fields
}

func TestAgentDetails(t *testing.T) {
	agent := gstructs.AgentDetails{
		Status: 1,
		Client: "test-client",
		Shasum: "123456789abcdef",
		Hints:  []string{"tag1", "tag2", "tag3"},
		Script: []byte("sample script content"),
		Answer: []byte("sample answer content"),
	}

	if agent.Status != 1 {
		t.Errorf("Expected agent status to be 1, got %d", agent.Status)
	}
	if agent.Client != "test-client" {
		t.Errorf("Expected agent client to be 'test-client', got %s", agent.Client)
	}
	// TODO: add more assertions for other agent fields
}
