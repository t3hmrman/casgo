package cas

import (
	"os"
)

var CONFIG_ENV_OVERRIDE_MAP map[string]string = map[string]string{
	"Host":               "CASGO_HOST",
	"Port":               "CASGO_PORT",
	"DBHost":             "CASGO_DBHOST",
	"DBName":             "CASGO_DBNAME",
	"CookieSecret":       "CASGO_SECRET",
	"TemplatesDirectory": "CASGO_TEMPLATES",
	"CompanyName":        "CASGO_COMPNAME",
	"DefaultAuthMethod":  "CASGO_DEFAULT_AUTH",
}

var CONFIG_DEFAULTS map[string]string = map[string]string{
	"Host":               "0.0.0.0",
	"Port":               "9090",
	"DBHost":             "localhost:28015",
	"DBName":             "casgo",
	"CookieSecret":       "secret-casgo-secret",
	"TemplatesDirectory": "templates/",
	"CompanyName":        "companyABC",
	"DefaultAuthMethod":  "password",
}

func NewCASServerConfig(userOverrides map[string]string) (map[string]string, error) {
	// Set default config values
	serverConfig := make(map[string]string)
	for k, v := range CONFIG_DEFAULTS {
		serverConfig[k] = v
	}

	// Override defaults with passed in map
	for k, v := range userOverrides {
		serverConfig[k] = v
	}

	return serverConfig, nil
}

func (c *CAS) overrideConfigWithEnv() {
	for configKey, envVarName := range CONFIG_ENV_OVERRIDE_MAP {
		if envValue := os.Getenv(envVarName); len(envValue) > 0 {
			c.Config[configKey] = envValue
		}
	}
}
