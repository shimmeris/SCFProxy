package tencent

import (
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
)

type Provider struct {
	region  string
	fclient *scf.Client
	gclient *apigateway.Client
}

func New(secretId, secretKey, region string) (*Provider, error) {
	credential := common.NewCredential(secretId, secretKey)

	fcpf := profile.NewClientProfile()
	fcpf.HttpProfile.Endpoint = "scf.tencentcloudapi.com"
	fclient, err := scf.NewClient(credential, region, fcpf)
	if err != nil {
		return nil, err
	}

	gcpf := profile.NewClientProfile()
	gcpf.HttpProfile.Endpoint = "apigateway.tencentcloudapi.com"
	gclient, err := apigateway.NewClient(credential, region, gcpf)
	if err != nil {
		return nil, err
	}

	provider := &Provider{
		region:  region,
		fclient: fclient,
		gclient: gclient,
	}
	return provider, nil
}

func (p *Provider) Name() string {
	return "tencent"
}

func (p *Provider) Region() string {
	return p.region
}

func (p *Provider) clearFunctionProxy(functionName, triggerName, triggerType string, onlyTrigger bool) error {
	if err := p.deleteTrigger(functionName, triggerName, triggerType); err != nil {
		if err, ok := err.(*errors.TencentCloudSDKError); !ok || (err.Code != scf.RESOURCENOTFOUND && err.Code != scf.RESOURCENOTFOUND_FUNCTION) {
			return err
		}
	}

	if onlyTrigger {
		return nil
	}

	if err := p.deleteFunction(functionName); err != nil {
		if err, ok := err.(*errors.TencentCloudSDKError); !ok || err.Code != scf.RESOURCENOTFOUND_FUNCTION {
			return err
		}
	}
	return nil
}

func (p *Provider) deleteFunction(functionName string) error {
	r := scf.NewDeleteFunctionRequest()
	r.FunctionName = common.StringPtr(functionName)

	_, err := p.fclient.DeleteFunction(r)
	return err
}

func (p *Provider) deleteTrigger(functionName, triggerName, triggerType string) error {
	r := scf.NewDeleteTriggerRequest()
	r.FunctionName = common.StringPtr(functionName)
	r.TriggerName = common.StringPtr(triggerName)
	r.Type = common.StringPtr(triggerType) // timer

	_, err := p.fclient.DeleteTrigger(r)
	return err
}

func Regions() []string {
	// 腾讯云大陆外地区部署延迟巨大，暂不进行部署
	return []string{
		"ap-beijing",
		"ap-chengdu",
		"ap-guangzhou",
		"ap-shanghai",
		"ap-nanjing",
		//"ap-hongkong",
		//"ap-mumbai",
		//"ap-singapore",
		//"ap-bangkok",
		//"ap-seoul",
		//"ap-tokyo",
		//"eu-frankfurt",
		//"eu-moscow",
		//"na-ashburn",
		//"na-toronto",
		//"na-siliconvalley",
	}
}
