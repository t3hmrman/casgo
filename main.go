package main

import (
	"io"
	"net/http"
	"log"
)

func handleIndex(res http.ResponseWriter, req *http.Request, ) {
	io.WriteString(res, "Yep")
}

func main() {
	http.HandleFunc("/", handleIndex)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server Failed:", err)
	}
}
