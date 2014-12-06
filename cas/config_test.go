package cas

/**
* Tests for Casgo configuration
 */

import (
	"os"
	"path/filepath"
	"testing"
)

// Ensure creating a test config produces non-nil
func TestConfigCreation(t *testing.T) {
	config, err := NewCASServerConfig(nil)
	if config == nil || err != nil {
		t.Error("Config with no options should be not fail, and should produce a non-nil pointer")
	}
}

// Ensure defaults are valid
func TestConfigDefaults(t *testing.T) {
	config, err := NewCASServerConfig(nil)
	if err != nil {
		t.Errorf("Error creating config %s:", err)
	}

	// Get the absolute path of the config default directory
	expectedTemplatesDirectory, err := filepath.Abs(CONFIG_DEFAULTS["templatesDirectory"])
	if err != nil {
		t.Errorf("Failed to retrieve absolute file path of templatesDirectory in test: %v", err)
	}
	CONFIG_DEFAULTS["templatesDirectory"] = expectedTemplatesDirectory

	// Ensure all configurations are at defaults (or expected modified defaults)
	for k, v := range config {
		if config[k] != CONFIG_DEFAULTS[k] {
			t.Errorf("Configuration key [%s] expected to be default (%s), was [%s]", k, CONFIG_DEFAULTS[k], v)
		}
	}
}

// Test User overriding in configuration setup
func TestUserOverrideHost(t *testing.T) {

	host := "fake-host-string"
	config, err := NewCASServerConfig(map[string]string{"host": host})
	if err != nil {
		t.Error("CASServerConfig creation failed")
	}

	if config["host"] != host {
		t.Errorf("Expected config[host] to be [%s], saw [%s]", host, config["host"])
	}

}

// Ensure override with env variables properly overrides
func TestEnvOverrideHost(t *testing.T) {

	previousEnvValue := os.Getenv("CASGO_HOST")
	host := "TESTHOST"
	err := os.Setenv("CASGO_HOST", host)
	if err != nil {
		t.Error("os.Setenv failed...")
	}

	config, err := NewCASServerConfig(nil)
	if err != nil {
		t.Error("CASServerConfig creation failed")
	}
	config = overrideConfigWithEnv(config)

	if config["host"] != host {
		t.Error("host property not properly ENV overriden.")
	}

	os.Setenv("CASGO_HOST", previousEnvValue)
}
