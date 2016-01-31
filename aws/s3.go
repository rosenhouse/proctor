package aws

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (c *Client) StoreObject(name string, dataBytes []byte,
	downloadFileName string, contentType string) error {
	bucketName, err := c.GetBucketName()
	if err != nil {
		return err
	}
	if err = c.EnsureBucketExists(bucketName); err != nil {
		return err
	}
	_, err = c.S3.PutObject(&s3.PutObjectInput{
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(dataBytes),
		Bucket:               aws.String(bucketName),
		ContentDisposition:   aws.String(fmt.Sprintf("attachment; filename=%s;", downloadFileName)),
		ContentType:          aws.String(contentType),
		Key:                  aws.String(name),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

func (c *Client) DeleteObject(name string) error {
	bucketName, err := c.GetBucketName()
	if err != nil {
		return err
	}
	_, err = c.S3.DeleteObject(&s3.DeleteObjectInput{
		Key:    aws.String(name),
		Bucket: aws.String(bucketName),
	})
	return err
}

func (c *Client) URLForObject(name string) (string, error) {
	bucketName, err := c.GetBucketName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, name), nil
}

func (c *Client) GetBucketName() (string, error) {
	if c.cachedBucketName == "" {
		accountNumber, err := c.GetAccountNumber()
		if err != nil {
			return "", err
		}
		c.cachedBucketName = fmt.Sprintf("bosh101-proctor-%s", accountNumber)
	}

	return c.cachedBucketName, nil
}

func (c *Client) EnsureBucketExists(bucketName string) error {
	_, err := c.S3.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if isAlreadyOwnedByMeError(err) {
		return nil
	}
	return err
}

// returns true if and only if the provided err is
// non-nil and is the S3 "BucketAlreadyOwnedByYou" error
// which occurs when trying to create a bucket that this account
// already owns
func isAlreadyOwnedByMeError(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		return awsErr.Code() == "BucketAlreadyOwnedByYou"
	}
	return false
}
