package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"

	"github.com/shimmeris/SCFProxy/sdk"
)

func (p *Provider) DeployReverse(opts *sdk.ReverseProxyOpts) error {
	apiId, err := p.createApi()
	if err != nil {
		return err
	}

	rootResourceId, err := p.getRootResourceId(apiId)
	if err != nil {
		return err
	}

	resourceId, err := p.createResource(apiId, rootResourceId)
	if err != nil {
		return err
	}

	if err = p.putMethod(apiId, rootResourceId); err != nil {
		return err
	}

	if err = p.putIntegration(apiId, rootResourceId, aws.String(opts.Origin)); err != nil {
		return err
	}

	if err = p.putMethod(apiId, resourceId); err != nil {
		return err
	}

	if err = p.putIntegration(apiId, resourceId, aws.String(fmt.Sprintf("%s/{proxy}", opts.Origin))); err != nil {
		return err
	}

	if err = p.createDeployment(apiId); err != nil {
		return err
	}

	return nil

}

func (p *Provider) createApi() (*string, error) {
	input := &apigateway.CreateRestApiInput{
		Name: aws.String(""),
		EndpointConfiguration: &types.EndpointConfiguration{
			Types: []types.EndpointType{types.EndpointTypeRegional},
		},
	}

	output, err := p.gclient.CreateRestApi(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return output.Id, nil
}

func (p *Provider) getRootResourceId(apiId *string) (*string, error) {
	input := &apigateway.GetResourcesInput{
		RestApiId: apiId,
	}

	output, err := p.gclient.GetResources(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return output.Items[0].Id, nil
}

func (p *Provider) createResource(apiId, resourceId *string) (*string, error) {
	input := &apigateway.CreateResourceInput{
		RestApiId: apiId,
		ParentId:  resourceId,
		PathPart:  aws.String("{proxy+}"),
	}

	output, err := p.gclient.CreateResource(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return output.Id, nil
}

func (p *Provider) putMethod(apiId, resourceId *string) error {
	input := &apigateway.PutMethodInput{
		RestApiId:         apiId,
		ResourceId:        resourceId,
		HttpMethod:        aws.String("ANY"),
		AuthorizationType: aws.String("NONE"),
		RequestParameters: map[string]bool{
			"method.request.path.proxy":                  true,
			"method.request.header.X-My-X-Forwarded-For": true,
		},
	}

	_, err := p.gclient.PutMethod(context.TODO(), input)
	return err
}

func (p *Provider) putIntegration(apiId, resourceId, uri *string) error {
	input := &apigateway.PutIntegrationInput{
		Uri:                   uri,
		RestApiId:             apiId,
		ResourceId:            resourceId,
		HttpMethod:            aws.String("ANY"),
		IntegrationHttpMethod: aws.String("ANY"),
		Type:                  types.IntegrationTypeHttpProxy,
		ConnectionType:        types.ConnectionTypeInternet,
		RequestParameters: map[string]string{
			"integration.request.path.proxy":             "method.request.path.proxy",
			"integration.request.header.X-Forwarded-For": "method.request.header.X-My-X-Forwarded-For",
		},
	}

	_, err := p.gclient.PutIntegration(context.TODO(), input)
	return err
}

func (p *Provider) createDeployment(apiId *string) error {
	input := &apigateway.CreateDeploymentInput{
		RestApiId: apiId,
		StageName: aws.String("proxy"),
	}

	_, err := p.gclient.CreateDeployment(context.TODO(), input)
	return err
}

func (p *Provider) deleteApi(apiId *string) error {
	input := &apigateway.DeleteRestApiInput{
		RestApiId: apiId,
	}

	_, err := p.gclient.DeleteRestApi(context.TODO(), input)
	return err
}
