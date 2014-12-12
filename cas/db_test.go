package cas_test

import (
	//. "github.com/t3hmrman/casgo/cas"

	//"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cas DB adapter", func() {

	Describe("DbExists function", func() {
		It("should return whether the database exists or not", func() {
			exists, casErr := testCASServer.Db.DbExists()
			Expect(casErr).To(BeNil())
			Expect(exists).To(Equal(true))
		})
	})

	Describe("GetDbName function", func() {
		It("should return the name of the server", func() {
			actual, expected := testCASServer.Db.GetDbName(), testCASServer.Config["dbName"]
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("GetUsersTableName function", func() {
		It("should return the default if not set differently", func() {
			actual, expected := testCASServer.Db.GetUsersTableName(), "users"
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("GetServicesTableName function", func() {
		It("should return the default if not set differently", func() {
			actual, expected := testCASServer.Db.GetServicesTableName(), "services"
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("GetTicketsTableName function", func() {
		It("should return the default if not set differently", func() {
			actual, expected := testCASServer.Db.GetTicketsTableName(), "tickets"
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("LoadJSONFixture function", func() {
		It("should not error when loading JSON into the database", func() {
			err := testCASServer.Db.LoadJSONFixture(
				testCASServer.Db.GetDbName(),
				testCASServer.Db.GetServicesTableName(),
				"fixtures/services.json",
			)
			Expect(err).To(BeNil())
		})

		It("Should increase the number of items in the given table", func() {
			err := testCASServer.Db.TeardownTable("services")
			Expect(err).To(BeNil())

			// Attempt to find a service in the fixture should fail
			service, err := testCASServer.Db.FindServiceByUrl("localhost:9090/validateCASLogin")
			Expect(err).ToNot(BeNil())
			Expect(service).To(BeNil())

			// Load the fixture
			err = testCASServer.Db.LoadJSONFixture(
				testCASServer.Db.GetDbName(),
				testCASServer.Db.GetServicesTableName(),
				"fixtures/services.json")
			Expect(err).To(BeNil())

			// Attempting to find a serive in the fixture shoudl pass now
			service, err = testCASServer.Db.FindServiceByUrl("localhost:9090/validateCASLogin")
			Expect(err).To(BeNil())
			Expect(service).ToNot(BeNil())

		})
	})

})
