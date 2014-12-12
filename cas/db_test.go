package cas_test

import (
	//. "github.com/t3hmrman/casgo/cas"

	//"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cas DB", func() {

	Describe("DbExists", func() {
		It("should return whether the database exists or not", func() {
			exists, casErr := testCASServer.Db.DbExists()
			Expect(casErr).To(BeNil())
			Expect(exists).To(Equal(true))
		})
	})

})
