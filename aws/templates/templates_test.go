package templates_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rosenhouse/proctor/aws/templates"
)

var _ = Describe("DefaultTemplate", func() {
	It("should match the fixture", func() {
		fixture, err := ioutil.ReadFile("fixtures/simple-stack.json")
		Expect(err).NotTo(HaveOccurred())
		Expect(templates.DefaultTemplate).To(MatchJSON(fixture))
	})
})
