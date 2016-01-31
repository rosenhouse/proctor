package aws_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rosenhouse/proctor/aws"
)

var _ = Describe("IAM", func() {
	Describe("GetAccountNumber", func() {
		It("should return the account number of the current IAM user", func() {
			accountNumber, err := awsClient.GetAccountNumber()
			Expect(err).NotTo(HaveOccurred())
			Expect(accountNumber).To(HaveLen(12))
		})
	})

	Describe("ParseARN", func() {
		It("should parse basic ARNs", func() {
			arnFormat0 := "arn:partition:service:region:account-id:resource"
			Expect(awsClient.ParseARN(arnFormat0)).To(Equal(aws.ARN{
				Partition: "partition",
				Service:   "service",
				Region:    "region",
				AccountID: "account-id",
				Resource:  "resource",
			}))
		})

		It("should group the resourcetype and resource together when they are colon-separated", func() {
			arnFormat1 := "arn:partition:service:region:account-id:resourcetype:resource"
			Expect(awsClient.ParseARN(arnFormat1)).To(Equal(aws.ARN{
				Partition: "partition",
				Service:   "service",
				Region:    "region",
				AccountID: "account-id",
				Resource:  "resourcetype:resource",
			}))
		})

		It("should group the resourcetype and resource together when they are slash separated", func() {
			arnFormat2 := "arn:partition:service:region:account-id:resourcetype/resource"
			Expect(awsClient.ParseARN(arnFormat2)).To(Equal(aws.ARN{
				Partition: "partition",
				Service:   "service",
				Region:    "region",
				AccountID: "account-id",
				Resource:  "resourcetype/resource",
			}))
		})

		It("should handle resources with arbitrary number of slashes", func() {
			sampleARN := "arn:aws:iam::123456789012:server-certificate/division_abc/subdivision_xyz/ProdServerCert"
			Expect(awsClient.ParseARN(sampleARN)).To(Equal(aws.ARN{
				Partition: "aws",
				Service:   "iam",
				Region:    "",
				AccountID: "123456789012",
				Resource:  "server-certificate/division_abc/subdivision_xyz/ProdServerCert",
			}))
		})

		Context("when the input string is malformed", func() {
			It("should return an error", func() {
				malformedARN := "arn:partition:service:region:account-id"
				_, err := awsClient.ParseARN(malformedARN)
				Expect(err).To(MatchError(fmt.Sprintf("malformed ARN string %q", malformedARN)))
			})
		})
	})
})
