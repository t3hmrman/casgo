package main

import (
	r "github.com/dancannon/gorethink"
	"github.com/t3hmrman/casgo/cas"
	"log"
	"net/http"
)

func main() {

	// Create new CAS Server config with default values
	config, err := cas.NewCASServerConfig()

	// Database setup
	dbSession, err := r.Connect(r.ConnectOpts{
		Address:  config.DBHost,
		Database: config.DBName,
	})

	if err != nil {
		log.Fatalln(err.Error())
	} else {
		config.RDBSession = dbSession
	}

	// Create CAS Server
	cas := cas.NewCASServer(config)

	// Setup handlers
	http.HandleFunc("/login", cas.HandleLogin)
	http.HandleFunc("/logout", cas.HandleLogout)
	http.HandleFunc("/register", cas.HandleRegister)
	http.HandleFunc("/", cas.HandleIndex)

	log.Printf("Starting CasGo on port %s...\n", config.Port)
	log.Fatal(http.ListenAndServe(config.GetAddr(), nil))
}
