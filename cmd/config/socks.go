package config

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type SocksConfig struct {
	mu      sync.RWMutex
	Records map[string]map[string]string
}

func LoadSocksConfig() (*SocksConfig, error) {
	conf := &SocksConfig{Records: make(map[string]map[string]string)}
	data, err := os.ReadFile(SocksProxyPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return conf, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, &conf.Records)
	return conf, err
}

func (c *SocksConfig) Has(provider, region string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.Records[provider][region]
	return ok
}

func (c *SocksConfig) Set(provider, region string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.Records[provider]
	if !ok {
		c.Records[provider] = make(map[string]string)
	}
	c.Records[provider][region] = ""
}

func (c *SocksConfig) Delete(provider, region string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.Records[provider], region)
}

func (c *SocksConfig) Save() error {
	return save(c.Records, SocksProxyPath)
}

func (c *SocksConfig) ToDoubleArray() [][]string {
	data := [][]string{}
	for provider, rmap := range c.Records {
		for region := range rmap {
			data = append(data, []string{provider, region})
		}
	}
	return data
}
