package main

import (
	"github.com/t3hmrman/casgo/cas"
	"log"
	"net/http"
)

func main() {

	// Create new CAS Server config with default values
	config, err := cas.NewCASServerConfig(nil)
	if err != nil {
		log.Fatal("Failed to create new CAS Server config...")
	}

	// Create CAS Server
	cas := cas.NewCASServer(config)

	// Setup handlers
	http.HandleFunc("/login", cas.HandleLogin)
	http.HandleFunc("/logout", cas.HandleLogout)
	http.HandleFunc("/register", cas.HandleRegister)
	http.HandleFunc("/", cas.HandleIndex)

	log.Printf("Starting CasGo on port %s...\n", cas.Config["Port"])
	log.Fatal(http.ListenAndServe(cas.GetAddr(), nil))
}
