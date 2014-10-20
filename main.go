package main

import (
	"net/http"
	"log"

	"github.com/t3hmrman/casgo/cas/v3"
	"github.com/unrolled/render"
)

// Configuration object
type CasServerConfig struct {
	host string
	port string
}

// Render object
var r = render.New(render.Options{})

func handleIndex(w http.ResponseWriter, req *http.Request) {
	r.HTML(w, http.StatusOK, "index", map[string]string{"companyName": "CompanyABC"})
}

func main() {

	// Configuration
	config := &CasServerConfig{"0.0.0.0", ":8080"}

	// Setup handlers
	http.HandleFunc("/login", cas.HandleLogin)
	http.HandleFunc("/", handleIndex)

	log.Printf("Starting CasGo on port %s...\n", config.port)
	err := http.ListenAndServe(config.port, nil)
	if err != nil {
		log.Fatal("Server Failed:", err)
	}
}
