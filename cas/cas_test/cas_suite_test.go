package cas_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"

	"github.com/t3hmrman/casgo/cas"
	"log"
	"net/http/httptest"
	"testing"
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
	// Start PhantomJS for integration tests
	StartPhantomJS()

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

	// Load database fixtures
	log.Printf("Loading database fixtures...")
	testCASServer.Db.LoadJSONFixture(
		testCASServer.Db.GetDbName(),
		testCASServer.Db.GetServicesTableName(),
		"../../fixtures/services.json",
	)
	testCASServer.Db.LoadJSONFixture(
		testCASServer.Db.GetDbName(),
		testCASServer.Db.GetUsersTableName(),
		"../../fixtures/users.json",
	)

})

var _ = AfterSuite(func() {
	testHTTPServer.Close()
	testCASServer.TeardownDb()
	StopWebdriver()
})
