package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/iam"
)

func (c *Client) GetAccountNumber() (string, error) {
	out, err := c.IAM.GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", err
	}
	arnString := *out.User.Arn

	arn, err := c.ParseARN(arnString)
	if err != nil {
		return "", err
	}

	return arn.AccountID, nil
}

// ARN represents an Amazon Resource Name
// http://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
type ARN struct {
	Partition string
	Service   string
	Region    string
	AccountID string
	Resource  string
}

// ParseARN parses an ARN string into its component fields
func (c *Client) ParseARN(arnString string) (ARN, error) {
	const numExpectedParts = 6
	parts := strings.SplitN(arnString, ":", numExpectedParts)
	if len(parts) < numExpectedParts {
		return ARN{}, fmt.Errorf("malformed ARN string %q", arnString)
	}
	return ARN{
		Partition: parts[1],
		Service:   parts[2],
		Region:    parts[3],
		AccountID: parts[4],
		Resource:  parts[5],
	}, nil
}
