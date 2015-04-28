package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/t3hmrman/casgo/cas"
	"testing"
)

// Testing globals for HTTP tests
var testCASConfig map[string]string
var testCASServer *cas.CAS

func TestCas(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CasGo DB Adapter Suite")
}

var _ = BeforeSuite(func() {
	// Setup CAS server & DB
	testCASConfig, err := cas.NewCASServerConfig(map[string]string{
		"companyName":        "Casgo Testing Company",
		"dbName":             "casgo_test",
		"templatesDirectory": "../templates",
	})
	Expect(err).To(BeNil())

	testCASServer, err = cas.NewCASServer(testCASConfig)
	Expect(err).To(BeNil())

	err = testCASServer.SetupDb()
	Expect(err).To(BeNil())

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
