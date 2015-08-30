package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/t3hmrman/casgo/cas"
	"net/http/httptest"
	"testing"
)

// Testing globals for HTTP tests
var testHTTPServer *httptest.Server
var testCASConfig map[string]string
var testCASServer *cas.CAS

func TestCasgoAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CasGo API Suite")
}

var _ = BeforeSuite(func() {
	// Setup CAS server & DB
	testCASConfig, _ = cas.NewCASServerConfig("")
	testCASConfig["companyName"] = "Casgo Testing Company"
	testCASConfig["dbName"] = "casgo_test"
	testCASConfig["templatesDirectory"] = "../templates"

	testCASServer, _ = cas.NewCASServer(testCASConfig)
	testCASServer.SetupDb()

	// Setup http test server
	testHTTPServer = httptest.NewTLSServer(testCASServer.ServeMux)

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
	testCASServer.Db.LoadJSONFixture(
		testCASServer.Db.GetDbName(),
		testCASServer.Db.GetApiKeysTableName(),
		"../../fixtures/api_keys.json",
	)

})

var _ = AfterSuite(func() {
	testHTTPServer.Close()
	testCASServer.TeardownDb()
})
