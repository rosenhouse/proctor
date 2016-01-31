package aws_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IAM", func() {
	Describe("GetAccountNumber", func() {
		It("should return the account number of the current IAM user", func() {
			accountNumber, err := awsClient.GetAccountNumber()
			Expect(err).NotTo(HaveOccurred())
			Expect(accountNumber).To(HaveLen(12))
		})
	})
})
