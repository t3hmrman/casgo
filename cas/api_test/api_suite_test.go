package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/t3hmrman/casgo/cas"
	"testing"
)

var testCASConfig map[string]string
var testCASServer *cas.CAS

func TestCas(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CasGo API Suite")
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
	testCASServer.TeardownDb()
})
