package main

import (
	"os"
	"net/http"
	"log"

	"github.com/t3hmrman/casgo/cas/v3"
	"github.com/unrolled/render"
)

// Configuration object
type CasServerConfig struct {
	host string
	port string
	companyName string
}

func (c *CasServerConfig) getAddr() string {
	return c.host + ":" + c.port
}

// Render object
var r = render.New(render.Options{})

func handleIndex(w http.ResponseWriter, req *http.Request) {
	r.HTML(w, http.StatusOK, "index", map[string]string{"companyName": "CompanyABC"})
}

func main() {

	// Configuration
	config := &CasServerConfig{"0.0.0.0", "8080", "companyABC"}

	// Environment overrides
	if v := os.Getenv("PORT"); len(v) > 0 { config.port = os.Getenv("PORT") }
	if v := os.Getenv("HOST"); len(v) > 0 { config.host = os.Getenv("HOST") }
	if v := os.Getenv("COMPANY_NAME"); len(v) > 0 { config.companyName = os.Getenv("COMPANY_NAME") }

	// Setup handlers
	http.HandleFunc("/login", cas.HandleLogin)
	http.HandleFunc("/", handleIndex)

	log.Printf("Starting CasGo on port %s...\n", config.port)
	err := http.ListenAndServe(config.getAddr(), nil)
	if err != nil {
		log.Fatal("CasGo server startup failed: ", err)
	}
}
