package main

import (
	"flag"
	"github.com/t3hmrman/casgo/cas"
	"log"
)

func main() {

	// Flag handling
	var configFilePath = flag.String("config", "", "Read configuration from specified (JSON) file")
	flag.Parse()

	// Create new CAS Server config with default values
	config, err := cas.NewCASServerConfig(*configFilePath)
	if err != nil {
		log.Fatalf("Failed to create new CAS server configuration, err: %v", err)
	}

	// Create CAS Server (registers appropriate handlers to http)
	casServer, err := cas.NewCASServer(config)
	if err != nil {
		log.Fatal("Failed to create new CAS Server instance...", err)
	}

	// Start the CAS Server
	log.Printf("Starting CasGo on port %s...\n", casServer.Config["port"])
	casServer.Start()
}
