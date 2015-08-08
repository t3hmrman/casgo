package integration_test

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

func TestCasgoEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CasGo integration Suite")
}

var _ = BeforeSuite(func() {
	// Start PhantomJS for integration tests
	StartPhantomJS()

	// Setup CAS server & DB
	testCASConfig, err := cas.NewCASServerConfig("")
	testCASConfig["companyName"] = "Casgo Testing Company"
	testCASConfig["dbName"] = "casgo_test"
	testCASConfig["templatesDirectory"] = "../templates"
	if err != nil {
		log.Fatalf("Failed to generate cas server config, err: %v", err)
	}

	testCASServer, err = cas.NewCASServer(testCASConfig)
	if err != nil {
		log.Fatalf("Failed to generate setup cas server, err: %v", err)
	}
	testCASServer.SetupDb()

	// Setup http test server
	testHTTPServer = httptest.NewServer(testCASServer.ServeMux)

	// Load database fixtures
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
