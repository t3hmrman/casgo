package cas

import "os"

// CAS server configuration object
type CASServerConfig struct {
	Host string
	Port string
	TemplatesDirectory string
	CompanyName string
}

func (c *CASServerConfig) GetAddr() string {
	return c.Host + ":" + c.Port
}

func (c *CASServerConfig) OverrideWithEnvVariables() {
	// Environment overrides
	if v := os.Getenv("HOST"); len(v) > 0 {
		c.Host = os.Getenv("HOST")
	}
	if v := os.Getenv("PORT"); len(v) > 0 {
		c.Port = os.Getenv("PORT")
	}
	if v := os.Getenv("TEMPLATES_DIR"); len(v) > 0 {
		c.TemplatesDirectory = os.Getenv("TEMPLATES_DIR")
	}
	if v := os.Getenv("COMPANY_NAME"); len(v) > 0 {
		c.CompanyName = os.Getenv("COMPANY_NAME")
	}
}