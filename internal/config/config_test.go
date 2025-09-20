package config

import (
	"io"
	"os"
	"testing"
)

func TestLoadConfigJSON(t *testing.T) {
	// Create a temporary config.json file
	jsonContent := `{
		"ListenAddr": ":9090",
		"TLSCertFile": "test_cert.pem",
		"TLSKeyFile": "test_key.pem",
		"BalancerType": "least_conn",
		"Backends": [
			{"URL": "http://test_backend1:3000", "Weight": 1, "MaxConnections": 50},
			{"URL": "http://test_backend2:3000", "Weight": 2, "MaxConnections": 100}
		]
	}`
	tempFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.WriteString(jsonContent)
	tempFile.Close()

	// Set the file path and load the config
	t.Logf("Temporary config.json path: %s", tempFile.Name())
	// Copy the temporary file to the current directory
	input, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer input.Close()
	output, err := os.Create("config.json")
	if err != nil {
		t.Fatalf("Failed to create config.json: %v", err)
	}
	defer output.Close()
	if _, err := io.Copy(output, input); err != nil {
		t.Fatalf("Failed to copy temp file: %v", err)
	}
	cfg := LoadConfig()

	// Validate the loaded config
	if cfg.ListenAddr != ":9090" {
		t.Errorf("Expected :9090, got %s", cfg.ListenAddr)
	}
	if cfg.TLSCertFile != "test_cert.pem" {
		t.Errorf("Expected test_cert.pem, got %s", cfg.TLSCertFile)
	}
	if cfg.BalancerType != "least_conn" {
		t.Errorf("Expected least_conn, got %s", cfg.BalancerType)
	}
	if len(cfg.Backends) != 2 {
		t.Errorf("Expected 2 backends, got %d", len(cfg.Backends))
	}
}
