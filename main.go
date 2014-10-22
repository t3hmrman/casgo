package main

import (
	"github.com/t3hmrman/casgo/cas"
	"log"
	"net/http"
)

func main() {

	// Configuration
	config := &cas.CASServerConfig{"0.0.0.0", "8080", "templates/", "companyABC"}
	config.OverrideWithEnvVariables()

	// Create CAS Server
	cas := cas.New(config)

	// Setup handlers
	http.HandleFunc("/login", cas.HandleLogin)
	http.HandleFunc("/", cas.HandleIndex)

	log.Printf("Starting CasGo on port %s...\n", config.Port)
	err := http.ListenAndServe(config.GetAddr(), nil)
	if err != nil {
		log.Fatal("CasGo server startup failed: ", err)
	}
}
