package cas

/**
 * CAS protocol V3 API
 */

import (
	"net/http"
	"github.com/unrolled/render"
)

// CAS server interface
type CASServer interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleLogout(w http.ResponseWriter, r *http.Request)
	HandleValidate(w http.ResponseWriter, r *http.Request)
	HandleServiceValidate(w http.ResponseWriter, r *http.Request)
	HandleProxyValidate(w http.ResponseWriter, r *http.Request)
	HandleProxy(w http.ResponseWriter, r *http.Request)
}

// CAS Server
type CAS struct {
	config *CASServerConfig
	render *render.Render
}

func New(config *CASServerConfig) *CAS {
	r := render.New(render.Options{Directory: config.TemplatesDirectory})
	c := &CAS{config, r}
	return c
}

// (Optional) Handles Index route
func (c *CAS) HandleIndex(w http.ResponseWriter, req *http.Request) {
	c.render.HTML(w, http.StatusOK, "index", map[string]string{"companyName": "CompanyABC"})
}

// Credential acceptor endpoint (requestor is Handled in main)
func (c *CAS) HandleLogin(w http.ResponseWriter, req *http.Request) {
	// Show login page if credentials are not provided, attempt login otherwise
	if true {
		c.render.HTML(w, http.StatusOK, "login", nil)
	} else {
		c.render.HTML(w, http.StatusOK, "login", nil)
	}
}

// Endpoint for destroying CAS sessions (logging out)
func (c *CAS) HandleLogout(w http.ResponseWriter, r *http.Request) {

}

// Endpoint for validating service tickets
func (c *CAS) HandleValidate(w http.ResponseWriter, r *http.Request) {

}

// Endpoint for validating service tickets (CAS 2.0)
func (c *CAS) HandleServiceValidate(w http.ResponseWriter, r *http.Request) {

}

// Endpoint for validating proxy tickets (CAS 2.0)
func (c *CAS) HandleProxyValidate(w http.ResponseWriter, r *http.Request) {

}

// Endpoint for handling proxy tickets (CAS 2.0)
func (c *CAS) HandleProxy(w http.ResponseWriter, r *http.Request) {

}
