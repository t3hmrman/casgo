package main

import (
    r "github.com/dancannon/gorethink"
    "github.com/t3hmrman/casgo/cas"
    "log"
    "net/http"
)

func main() {

    // Configuration
    config := &cas.CASServerConfig{
        "0.0.0.0",
        "8080",
        "localhost:28015",
        "casgo",
        "templates/",
        "companyABC",
    }
    config.OverrideWithEnvVariables()

    // Database setup
    var session = *r.Session
    session, err := r.Connect(r.ConnectOps{
        Address:  config.DBHost,
        Database: config.DBName,
    })

    if err != nil {
        log.Fatalln(err.Error())
    }

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
