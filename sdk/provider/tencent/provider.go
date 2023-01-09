package tencent

import (
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
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
