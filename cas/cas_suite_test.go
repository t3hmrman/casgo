package cas_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/t3hmrman/casgo/cas"

	"net/http/httptest"
	"testing"
	"log"
)

// Testing globals for HTTP tests
var testHTTPServer *httptest.Server
var testCASConfig map[string]string
var testCASServer *cas.CAS

func TestCas(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cas Suite")
}

var _ = BeforeSuite(func() {
	// Setup CAS server & DB
	testCASConfig, _ = cas.NewCASServerConfig(map[string]string{
		"companyName":        "Casgo Testing Company",
		"dbName":             "casgo_test",
		"templatesDirectory": "../templates",
	})
	testCASServer, _ = cas.NewCASServer(testCASConfig)
	testCASServer.SetupDb()

	// Setup http test server
	testHTTPServer = httptest.NewServer(testCASServer.ServeMux)
	log.Printf("Started testing HTTP server @ %s", testHTTPServer.URL)
})

var _ = AfterSuite(func() {
	testHTTPServer.Close()
	testCASServer.TeardownDb()
})
