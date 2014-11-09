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
        nil,
    }
    config.OverrideWithEnvVariables()

    // Database setup
    var session *r.Session
    session, err := r.Connect(r.ConnectOpts{
        Address:  config.DBHost,
        Database: config.DBName,
    })

    if err != nil {
        log.Fatalln(err.Error())
    } else {
        config.RDBSession = session;
    }

    // Create CAS Server
    cas := cas.New(config)

    // Setup handlers
    http.HandleFunc("/login", cas.HandleLogin)
    http.HandleFunc("/", cas.HandleIndex)

    log.Printf("Starting CasGo on port %s...\n", config.Port)
    log.Fatal(http.ListenAndServe(config.GetAddr(), nil))
}
