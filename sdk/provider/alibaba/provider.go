package alibaba

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	fcopen "github.com/alibabacloud-go/fc-open-20210406/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Provider struct {
	region  string
	fclient *fcopen.Client
	runtime *util.RuntimeOptions
}

func New(accessKeyId, accessKeySecret, accountId, region string) (*Provider, error) {
	auth := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(region),
		Endpoint:        tea.String(fmt.Sprintf("%s.%s.fc.aliyuncs.com", accountId, region)),
	}

	fclient, err := fcopen.NewClient(auth)
	if err != nil {
		return nil, err
	}

	provider := &Provider{
		region:  region,
		fclient: fclient,
		runtime: &util.RuntimeOptions{},
	}
	return provider, nil
}

func (p *Provider) Name() string {
	return "alibaba"
}

func (p *Provider) Region() string {
	return p.region
}

func (p *Provider) clearProxy(serviceName, functionName, triggerName string, onlyTrigger bool) error {
	if err := p.deleteTrigger(serviceName, functionName, triggerName); err != nil {
		if err, ok := err.(*tea.SDKError); !ok || *err.StatusCode != 404 {
			return err
		}
	}

	if onlyTrigger {
		return nil
	}

	if err := p.deleteFunction(serviceName, functionName); err != nil {
		if err, ok := err.(*tea.SDKError); !ok || *err.StatusCode != 404 {
			return err
		}
	}
	return nil
}

func (p *Provider) deleteService(serviceName string) error {
	h := &fcopen.DeleteServiceHeaders{}
	_, err := p.fclient.DeleteServiceWithOptions(tea.String(serviceName), h, p.runtime)
	return err
}

func (p *Provider) deleteFunction(serviceName, functionName string) error {
	h := &fcopen.DeleteFunctionHeaders{}
	_, err := p.fclient.DeleteFunctionWithOptions(tea.String(serviceName), tea.String(functionName), h, p.runtime)
	return err
}

func (p *Provider) deleteTrigger(serviceName, functionName, triggerName string) error {
	deleteTriggerHeaders := &fcopen.DeleteTriggerHeaders{}
	_, err := p.fclient.DeleteTriggerWithOptions(
		tea.String(serviceName),
		tea.String(functionName),
		tea.String(triggerName),
		deleteTriggerHeaders,
		p.runtime,
	)
	return err
}

func Regions() []string {
	return []string{
		"ap-northeast-1",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-southeast-3",
		"ap-southeast-5",
		"cn-beijing",
		"cn-chengdu",
		"cn-hangzhou",
		"cn-hongkong",
		"cn-huhehaote",
		"cn-qingdao",
		"cn-shanghai",
		"cn-shenzhen",
		"cn-zhangjiakou",
		"eu-central-1",
		"eu-west-1",
		"us-east-1",
		"us-west-1",
	}
}
