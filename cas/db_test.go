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

})
