package cas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var CONFIG_ENV_OVERRIDE_MAP map[string]string = map[string]string{
	"host":               "CASGO_HOST",
	"port":               "CASGO_PORT",
	"dbHost":             "CASGO_DBHOST",
	"dbName":             "CASGO_DBNAME",
	"cookieSecret":       "CASGO_SECRET",
	"templatesDirectory": "CASGO_TEMPLATES",
	"companyName":        "CASGO_COMPNAME",
	"authMethod":         "CASGO_DEFAULT_AUTH",
	"logLevel":           "CASGO_LOG_LVL",
}

var CONFIG_DEFAULTS map[string]string = map[string]string{
	"host":               "0.0.0.0",
	"port":               "9090",
	"dbHost":             "localhost:28015",
	"dbName":             "casgo",
	"cookieSecret":       "secret-casgo-secret",
	"templatesDirectory": "templates/",
	"companyName":        "companyABC",
	"authMethod":         "password",
	"logLevel":           "WARN",
}

// Create default casgo configuration, with user overrides if any
func NewCASServerConfig(configFilePath string) (map[string]string, error) {
	// Set default config values
	serverConfig := make(map[string]string)
	for k, v := range CONFIG_DEFAULTS {
		serverConfig[k] = v
	}

	// Read config from file, if non-empty config file path is passed
	if len(configFilePath) > 0 {
		// Read file
		buf, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Failed to read given configuration file [%s] %v", configFilePath, err)
		}

		if len(buf) == 0 {
			log.Printf("[WARNING] Loaded configuration file was empty")
		}

		// Load JSON
		fileConfig := make(map[string]string)
		err = json.Unmarshal(buf, &fileConfig)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Failed to unmarshal JSON in configuration file [%s], %v", configFilePath, err)
		}

		// Copy values from loaded file config to actual server config
		for k, v := range fileConfig {
			serverConfig[k] = v
		}

	}

	// Override config with what is stored in env
	serverConfig = overrideConfigWithEnv(serverConfig)

	// Update filepath with absolute path
	absDirPath, err := filepath.Abs(serverConfig["templatesDirectory"])
	if err != nil {
		log.Printf("[WARNING] Failed to resolve absolute path for templatesDirectory %s", serverConfig["templatesDirectory"])
	} else {
		serverConfig["templatesDirectory"] = absDirPath
	}

	return serverConfig, nil
}

// Override a configuration hash with values provided by ENV
func overrideConfigWithEnv(config map[string]string) map[string]string {
	for configKey, envVarName := range CONFIG_ENV_OVERRIDE_MAP {
		if envValue := os.Getenv(envVarName); len(envValue) > 0 {
			config[configKey] = envValue
		}
	}
	return config
}
