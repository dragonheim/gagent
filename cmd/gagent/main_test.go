package main_test

import (
	"io/ioutil"
	"os"
	"testing"

	main "github.com/dragonheim/gagent/cmd/gagent"
	gstructs "github.com/dragonheim/gagent/internal/gstructs"
)

// This function will create a temporary config file for testing purposes
func createTestConfigFile() (string, error) {
	tmpfile, err := ioutil.TempFile("", "test_config_*.hcl")
	if err != nil {
		return "", err
	}

	content := []byte(`mode = "setup"
listen_addr = "0.0.0.0"
monitor_port = 8888
client_port = 35572
router_port = 35570
worker_port = 35571
`)
	if _, err := tmpfile.Write(content); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func TestMain(t *testing.T) {
	t.Run("Test setup mode with temp config file", func(t *testing.T) {
		tmpConfig, err := createTestConfigFile()
		if err != nil {
			t.Fatalf("Failed to create temp config file: %v", err)
		}
		defer os.Remove(tmpConfig)

		config := gstructs.GagentConfig{
			File: tmpConfig,
			Mode: "setup",
		}

		// Run the main function with the temporary config
		main.Run(config)

		// Check if the config has been set up correctly
		expectedConfig := gstructs.GagentConfig{
			Mode:        "setup",
			ListenAddr:  "0.0.0.0",
			MonitorPort: 8888,
			ClientPort:  35572,
			RouterPort:  35570,
			WorkerPort:  35571,
		}

		if config.Mode != expectedConfig.Mode ||
			config.ListenAddr != expectedConfig.ListenAddr ||
			config.MonitorPort != expectedConfig.MonitorPort ||
			config.ClientPort != expectedConfig.ClientPort ||
			config.RouterPort != expectedConfig.RouterPort ||
			config.WorkerPort != expectedConfig.WorkerPort {
			t.Fatalf("Expected config %+v, got %+v", expectedConfig, config)
		}
	})
}
