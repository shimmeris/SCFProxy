package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/shimmeris/SCFProxy/function"
	"github.com/shimmeris/SCFProxy/sdk"
)

func (p *Provider) DeploySocksProxy(opts *sdk.FunctionOpts) error {
	return p.createSocksFunction(opts.FunctionName)
}

func (p *Provider) ClearSocksProxy(opts *sdk.FunctionOpts) error {
	return p.deleteFunction(opts.FunctionName)
}

func (p *Provider) createSocksFunction(functionName string) error {
	input := &lambda.CreateFunctionInput{
		FunctionName:  aws.String(functionName),
		Code:          &types.FunctionCode{ZipFile: function.AwsSocksCodeZip},
		Handler:       aws.String("main"),
		Runtime:       types.RuntimeGo1x,
		MemorySize:    aws.Int32(128),
		Architectures: []types.Architecture{types.ArchitectureX8664},
		Timeout:       aws.Int32(900),
		PackageType:   types.PackageTypeZip,
		Role:          p.roleArn,
	}

	_, err := p.fclient.CreateFunction(p.ctx, input)
	return err
}

func (p *Provider) InvokeFunction(opts *sdk.FunctionOpts, message string) error {
	input := &lambda.InvokeInput{
		FunctionName:   aws.String(opts.FunctionName),
		InvocationType: types.InvocationTypeEvent,
		Payload:        []byte(message),
	}

	_, err := p.fclient.Invoke(p.ctx, input)
	return err
}
