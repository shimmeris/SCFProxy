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
