package aws

//
//import (
//	"SCFProxy-go/sdk"
//	"github.com/aws/aws-sdk-go-v2/aws"
//	"github.com/aws/aws-sdk-go-v2/credentials"
//	"github.com/aws/aws-sdk-go-v2/service/lambda"
//)
//
//type FunctionConfig struct {
//}
//
//func NewFunctionConfig() {
//	lambda.CreateFunctionInput{
//		FunctionName: &sdk.HTTPFunctionName,
//
//
//	}
//}
//
//type Provider struct {
//	client  *lambda.Client
//	region  string
//	fconfig *FunctionConfig
//}
//
//func New(accessKey, secretKey, region string, fconfig *FunctionConfig) sdk.Provider {
//	opts := lambda.Options{
//		Region:      region,
//		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
//	}
//	client := lambda.New(opts)
//	return &Provider{
//		client:  client,
//		region:  region,
//		fconfig: fconfig,
//	}
//}
//
//func (p *Provider) Name() string {
//	return "aws"
//}
//
//func (p *Provider) Deploy() (*sdk.Result, error) {
//
//}
//
//func (p *Provider)
