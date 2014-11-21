package cas

import (
    "os"
    r "github.com/dancannon/gorethink"
)

// CAS server configuration object
type CASServerConfig struct {
    Host string
    Port string
    DBHost string
    DBName string
    CookieSecret string
    TemplatesDirectory string
    CompanyName string
    RDBSession *r.Session
    DefaultAuthMethod string
}

func NewCASServerConfig() (*CASServerConfig, error)  {
	// Default values
	config := &CASServerConfig{
		Host: "0.0.0.0",
		Port: "9090",
		DBHost: "localhost:28015",
		DBName: "casgo",
		CookieSecret: "my-super-secret-casgo-secret",
		TemplatesDirectory: "templates/",
		CompanyName: "companyABC",
	}

	// ENV overrides
	config.OverrideWithEnvVariables()

	return config, nil
}

func (c *CASServerConfig) GetAddr() string {
    return c.Host + ":" + c.Port
}

func (c *CASServerConfig) OverrideWithEnvVariables() {
    // Environment overrides
    if v := os.Getenv("CASGO_HOST"); len(v) > 0 {
        c.Host = os.Getenv("CASGO_HOST")
    }
    if v := os.Getenv("CASGO_PORT"); len(v) > 0 {
        c.Port = os.Getenv("CASGO_PORT")
    }
    if v := os.Getenv("CASGO_DBHOST"); len(v) > 0 {
        c.DBHost = os.Getenv("CASGO_DBHOST")
    }
    if v := os.Getenv("CASGO_DBNAME"); len(v) > 0 {
        c.DBName = os.Getenv("CASGO_DBNAME")
    }
    if v := os.Getenv("CASGO_SECRET"); len(v) > 0 {
        c.CookieSecret = os.Getenv("CASGO_SECRET")
    }
    if v := os.Getenv("CASGO_TEMPLATES"); len(v) > 0 {
        c.TemplatesDirectory = os.Getenv("CASGO_TEMPLATES")
    }
    if v := os.Getenv("CASGO_COMPNAME"); len(v) > 0 {
        c.CompanyName = os.Getenv("CASGO_COMPNAME")
    }
    if v := os.Getenv("CASGO_DEFAULT_AUTH"); len(v) > 0 {
        c.DefaultAuthMethod = os.Getenv("CASGO_DEFAULT_AUTH")
    }
}
