package cas

import "os"

// CAS server configuration object
type CASServerConfig struct {
	Host string
	Port string
	DBHost string
	DBName string
	TemplatesDirectory string
	CompanyName string
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
		c.Port = os.Getenv("CASGO_DBHOST")
	}
	if v := os.Getenv("CASGO_DBNAME"); len(v) > 0 {
		c.Port = os.Getenv("CASGO_DBNAME")
	}
	if v := os.Getenv("CASGO_TEMPLATES"); len(v) > 0 {
		c.TemplatesDirectory = os.Getenv("CASGO_TEMPLATES")
	}
	if v := os.Getenv("CASGO_COMPNAME"); len(v) > 0 {
		c.CompanyName = os.Getenv("CASGO_COMPNAME")
	}
}
