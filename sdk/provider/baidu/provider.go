package baidu

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/cfc"
)

type Provider struct {
	region  string
	fClient *cfc.Client
}

func New(accessKeyId, accessKeySecret, region string) (*Provider, error) {

	// 用户指定的Endpoint
	endpoint := fmt.Sprintf("https://cfc.%s.baidubce.com", region)

	// 初始化一个cfcClient
	fClient, err := cfc.NewClient(accessKeyId, accessKeySecret, endpoint)
	if err != nil {
		return nil, err
	}
	provider := &Provider{
		region:  region,
		fClient: fClient,
	}
	return provider, nil
}

func (p *Provider) Name() string {
	return "baidu"
}

func (p *Provider) Region() string {
	return p.region
}

func Regions() []string {
	return []string{
		"bj",
		"gz",
		"sz",
	}
}
