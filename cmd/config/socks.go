package config

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"
)

type SocksRecord struct {
	Host string
	Port int
	Key  string
}

type SocksConfig struct {
	mu      sync.RWMutex
	Records map[string]map[string]*SocksRecord
}

func LoadSocksConfig() (*SocksConfig, error) {
	conf := &SocksConfig{Records: make(map[string]map[string]*SocksRecord)}
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

func (c *SocksConfig) Get(provider, region string) (*SocksRecord, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	record, ok := c.Records[provider][region]
	return record, ok
}

func (c *SocksConfig) Set(provider, region string, record *SocksRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.Records[provider]
	if !ok {
		c.Records[provider] = make(map[string]*SocksRecord)
	}
	c.Records[provider][region] = record
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
		for region, record := range rmap {
			data = append(data, []string{provider, region, record.Host, strconv.Itoa(record.Port), record.Key})
		}
	}
	return data
}
