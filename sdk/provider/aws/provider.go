package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type Provider struct {
	fclient *lambda.Client
	gclient *apigateway.Client
	region  string
	roleArn *string
	ctx     context.Context
}

func New(accessKey, secretKey, region, roleArn string) (*Provider, error) {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))

	return &Provider{
		fclient: lambda.New(lambda.Options{
			Region:      region,
			Credentials: creds,
		}),
		gclient: apigateway.New(apigateway.Options{
			Region:      region,
			Credentials: creds,
		}),
		region:  region,
		roleArn: aws.String(roleArn),
		ctx:     context.TODO(),
	}, nil
}

func (p *Provider) Name() string {
	return "aws"
}

func (p *Provider) Region() string {
	return p.region
}

func Regions() []string {
	return []string{
		"us-east-2",
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"af-south-1",
		"ap-east-1",
		"ap-southeast-3",
		"ap-south-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"ca-central-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-south-1",
		"eu-west-3",
		"eu-north-1",
		"me-south-1",
		"me-central-1",
		"sa-east-1",
		"us-gov-east-1",
		"us-gov-west-1",
	}
}
