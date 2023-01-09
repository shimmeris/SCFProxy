package aws

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/shimmeris/SCFProxy/function"
	"github.com/shimmeris/SCFProxy/sdk"
)

func (p *Provider) DeployHttpProxy(opts *sdk.FunctionOpts) (string, error) {
	if !opts.OnlyTrigger {
		if err := p.createHttpFunction(opts.FunctionName); err != nil {
			return "", err
		}

		if err := p.addPermission(opts.FunctionName); err != nil {
			return "", err
		}
	}

	api, err := p.createHttpTrigger(opts.FunctionName)
	if err != nil {
		return "", err
	}
	return api, nil
}

func (p *Provider) ClearHttpProxy(opts *sdk.FunctionOpts) error {
	if err := p.deleteHttpTrigger(opts.FunctionName); err != nil {
		return err
	}

	if opts.OnlyTrigger {
		return nil
	}

	return p.deleteFunction(opts.FunctionName)
}

func (p *Provider) createHttpFunction(functionName string) error {
	input := &lambda.CreateFunctionInput{
		FunctionName:  aws.String(functionName),
		Code:          &types.FunctionCode{ZipFile: []byte(function.AwsHttpCodeZip)},
		Handler:       aws.String("index.handler"),
		MemorySize:    aws.Int32(128),
		Architectures: []types.Architecture{types.ArchitectureArm64},
		Timeout:       aws.Int32(10),
		Runtime:       types.RuntimePython39,
		PackageType:   types.PackageTypeZip,
		Role:          p.roleArn,
	}

	_, err := p.fclient.CreateFunction(p.ctx, input)
	return err
}

func (p *Provider) createHttpTrigger(functionName string) (string, error) {
	input := &lambda.CreateFunctionUrlConfigInput{
		FunctionName: aws.String(functionName),
		AuthType:     types.FunctionUrlAuthTypeNone,
	}

	output, err := p.fclient.CreateFunctionUrlConfig(p.ctx, input)
	if err != nil {
		return "", err
	}
	return *output.FunctionUrl, nil
}

func (p *Provider) addPermission(functionName string) error {
	input := &lambda.AddPermissionInput{
		FunctionName:        aws.String(functionName),
		Action:              aws.String("lambda:InvokeFunctionUrl"),
		FunctionUrlAuthType: types.FunctionUrlAuthTypeNone,
		Principal:           aws.String("*"),
		StatementId:         aws.String("FunctionURLAllowPublicAccess"),
	}

	_, err := p.fclient.AddPermission(p.ctx, input)
	return err
}

func (p *Provider) deleteFunction(functionName string) error {
	input := &lambda.DeleteFunctionInput{
		FunctionName: aws.String(functionName),
	}

	_, err := p.fclient.DeleteFunction(p.ctx, input)
	if err != nil {
		var rnf *types.ResourceNotFoundException
		if errors.As(err, &rnf) {
			return nil
		}
	}
	return err
}

func (p *Provider) deleteHttpTrigger(functionName string) error {
	input := &lambda.DeleteFunctionUrlConfigInput{
		FunctionName: aws.String(functionName),
	}

	_, err := p.fclient.DeleteFunctionUrlConfig(p.ctx, input)
	if err != nil {
		var rnf *types.ResourceNotFoundException
		if errors.As(err, &rnf) {
			return nil
		}
	}
	return err
}
