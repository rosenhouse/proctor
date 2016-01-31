package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Endpoints struct {
	Route53        string
	EC2            string
	S3             string
	Cloudformation string
}

type Client struct {
	EC2              EC2Client
	S3               S3Client
	Route53          Route53Client
	Cloudformation   CloudformationClient
	IAM              IAMClient
	cachedBucketName string
}

type IAMClient interface {
	GetUser(input *iam.GetUserInput) (*iam.GetUserOutput, error)
}

type S3Client interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

type EC2Client interface {
	CreateKeyPair(*ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error)
	DeleteKeyPair(*ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error)
	DescribeKeyPairs(*ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error)
	DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
}

type Route53Client interface {
	ChangeResourceRecordSets(input *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error)
	ListResourceRecordSets(input *route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error)
}

type CloudformationClient interface {
	CreateStack(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DeleteStack(*cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error)
	DescribeStacks(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
}

type AWSError struct {
	Method string
	Err    error
}

func (e *AWSError) Error() string {
	return fmt.Sprintf("%s: %s", e.Method, e.Err)
}

type Config struct {
	AccessKey         string
	SecretKey         string
	RegionName        string
	EndpointOverrides map[string]string
}

func (c *Config) getEndpoint(serviceName string) (*aws.Config, error) {
	if c.EndpointOverrides == nil {
		return &aws.Config{}, nil
	}
	endpointOverride, ok := c.EndpointOverrides[serviceName]
	if !ok || endpointOverride == "" {
		return nil, fmt.Errorf("EndpointOverrides set, but missing required service %q", serviceName)
	}
	return &aws.Config{Endpoint: aws.String(endpointOverride)}, nil
}

func New(config Config) (*Client, error) {
	credentials := credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, "")
	sdkConfig := &aws.Config{
		Credentials: credentials,
		Region:      aws.String(config.RegionName),
	}

	session := session.New(sdkConfig)

	route53EndpointConfig, err := config.getEndpoint("route53")
	if err != nil {
		return nil, err
	}

	ec2EndpointConfig, err := config.getEndpoint("ec2")
	if err != nil {
		return nil, err
	}

	s3EndpointConfig, err := config.getEndpoint("s3")
	if err != nil {
		return nil, err
	}

	cloudformationEndpointConfig, err := config.getEndpoint("cloudformation")
	if err != nil {
		return nil, err
	}

	iamEndpointConfig, err := config.getEndpoint("iam")
	if err != nil {
		return nil, err
	}

	return &Client{
		EC2:            ec2.New(session, ec2EndpointConfig),
		S3:             s3.New(session, s3EndpointConfig),
		Route53:        route53.New(session, route53EndpointConfig),
		Cloudformation: cloudformation.New(session, cloudformationEndpointConfig),
		IAM:            iam.New(session, iamEndpointConfig),
	}, nil
}

func toStringPointers(strings ...string) []*string {
	var output []*string
	for _, s := range strings {
		output = append(output, aws.String(s))
	}
	return output
}
