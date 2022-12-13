package cmd

import (
	"github.com/shimmeris/SCFProxy/cmd/config"
	"github.com/shimmeris/SCFProxy/sdk"
	"github.com/shimmeris/SCFProxy/sdk/provider/alibaba"
	"github.com/shimmeris/SCFProxy/sdk/provider/tencent"
)

const (
	HTTPFunctionName = "scf_http"
	HTTPTriggerName  = "http_trigger"

	SocksFunctionName = "scf_socks"
	SocksTriggerName  = "socks_trigger"
)

var (
	allProviders     = []string{"alibaba", "tencent"}
	httpProviders    = []string{"alibaba", "tencent"}
	socksProviders   = []string{"alibaba", "tencent"}
	reverseProviders = []string{"tencent"}
)

func listProviders(module string) []string {
	switch module {
	case "http":
		return httpProviders
	case "socks":
		return socksProviders
	case "reverse":
		return reverseProviders
	default:
		return allProviders
	}
}

func listRegions(provider string) []string {
	switch provider {
	case "alibaba":
		return alibaba.Regions()
	case "tencent":
		return tencent.Regions()
	default:
		return nil
	}
}

func createProvider(name, region string, config *config.ProviderConfig) (sdk.Provider, error) {
	c := config.ProviderCredentialByName(name)
	ak := c.AccessKeyId
	sk := c.AccessKeySecret
	switch name {
	case "alibaba":
		accountId := c.AccountId
		return alibaba.New(ak, sk, accountId, region)
	//case "huawei":
	//	return huawei.New(ak, sk, region), nil
	case "tencent":
		return tencent.New(ak, sk, region)
	default:
		return nil, nil
	}
}
