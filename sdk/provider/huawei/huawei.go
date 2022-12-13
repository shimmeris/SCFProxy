package huawei

// Temporarily can not use Huawei Cloud as a proxy, because its sdk has problems, need to wait for its repair

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	functiongraph "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/functiongraph/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/functiongraph/v2/model"
	functionregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/functiongraph/v2/region"

	"github.com/shimmeris/SCFProxy/sdk/provider/huawei/sign"
)

type Provider struct {
	fclient *functiongraph.FunctionGraphClient
	region  string
	signer  *sign.Signer
}

func New(ak, sk, region string) *Provider {
	auth := basic.NewCredentialsBuilder().WithAk(ak).WithSk(sk).Build()
	hcClient := functiongraph.FunctionGraphClientBuilder().
		WithRegion(functionregion.ValueOf(region)).
		WithCredential(auth).
		Build()

	fclient := functiongraph.NewFunctionGraphClient(hcClient)
	signer := &sign.Signer{Key: ak, Secret: sk}

	provider := &Provider{fclient: fclient, region: region, signer: signer}
	return provider

}

func (p *Provider) Name() string { return "huawei" }

func (p *Provider) Region() string {
	return p.region
}

func (p *Provider) deleteFunction(functionUrn string) {
	request := &model.DeleteFunctionRequest{}
	request.FunctionUrn = functionUrn
	response, err := p.fclient.DeleteFunction(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
}

func (p *Provider) deleteTrigger(functionUrn, triggerId string) {

	request := &model.DeleteFunctionTriggerRequest{}
	request.FunctionUrn = functionUrn
	request.TriggerTypeCode = model.GetDeleteFunctionTriggerRequestTriggerTypeCodeEnum().APIG
	request.TriggerId = triggerId
	response, err := p.fclient.DeleteFunctionTrigger(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
}

func Regions() []string {
	// accquired from https://github.com/huaweicloud/huaweicloud-sdk-go-v3/blob/v0.1.5/services/functiongraph/v2/region/region.go
	return []string{
		"cn-north-4",
		"cn-north-1",
		"cn-east-2",
		"cn-east-3",
		"cn-south-1",
		"cn-southwest-2",
		"ap-southeast-2",
		"ap-southeast-1",
		"ap-southeast-3",
		"af-south-1",
	}
}
