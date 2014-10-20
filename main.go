package main

import (
	"net/http"
	"log"

	"github.com/t3hmrman/casgo/cas/v3"
	"github.com/unrolled/render"
)

// Render object
var r = render.New(render.Options{})

func handleIndex(w http.ResponseWriter, req *http.Request) {
	r.HTML(w, http.StatusOK, "login", map[string]string{"companyName": "CompanyABC"})
}

func main() {
	http.HandleFunc("/login", cas.HandleLogin)
	http.HandleFunc("/", handleIndex)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server Failed:", err)
	}
}
