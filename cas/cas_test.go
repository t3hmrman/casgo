package cas

import (
	"golang.org/x/net/html"
	"testing"
	"net/http/httptest"
	"net/http"
	"log"
	"fmt"
)

func TestLoginEndpoint(t *testing.T){
	if (true == false) {
		t.Error("useless test")
	}
}

// Login page tests
func TestLoginPage(t *testing.T) {
	if testing.Short() { t.Skip("Skipping integration test (in short mode).") }

	testHandler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {		fmt.Fprintln(w, "Hello, client")
	})
	
	loginServer := httptest.NewServer(testHandler)
	defer loginServer.Close()

	res, err := http.Get(loginServer.URL)
	if err != nil { log.Fatal(err) }

	doc, err := html.Parse(res.Body)
	if err != nil { log.Fatal(err) }

	if doc == nil {
		t.Error("Doc is nil!")
	}
}
