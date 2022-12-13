package config

import (
	"os"

	"github.com/pelletier/go-toml"
)

const ProviderConfigContent = `[alibaba]
AccessKeyId = ""
AccessKeySecret = "" 
AccountId = ""

[tencent]
# Named SecretId in tencent
AccessKeyId = ""
# Named SecretKey in tencent
AccessKeySecret = "" 

[huawei]
AccessKeyId = ""
AccessKeySecret = "" 
`

type Credential struct {
	AccessKeyId     string
	AccessKeySecret string
	AccountId       string
}

func (c Credential) isSet() bool {
	return c.AccessKeyId != "" && c.AccessKeySecret != ""
}

type ProviderConfig struct {
	Alibaba *Credential
	Tencent *Credential
}

func LoadProviderConfig(path string) (*ProviderConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &ProviderConfig{}
	if err := toml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *ProviderConfig) ProviderCredentialByName(provider string) *Credential {
	switch provider {
	case "alibaba":
		return c.Alibaba
	case "tencent":
		return c.Tencent
	default:
		return nil
	}
}

func (c *ProviderConfig) IsSet(provider string) bool {
	cred := c.ProviderCredentialByName(provider)
	if cred == nil {
		return false
	}
	return cred.isSet()
}
