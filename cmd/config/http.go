package config

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type HttpRecord struct {
	Api string
}

type HttpConfig struct {
	mu      sync.RWMutex
	Records map[string]map[string]*HttpRecord
}

func LoadHttpConfig() (*HttpConfig, error) {
	conf := &HttpConfig{Records: make(map[string]map[string]*HttpRecord)}
	data, err := os.ReadFile(HttpProxyPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return conf, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, &conf.Records)
	return conf, err
}

func (c *HttpConfig) Get(provider, region string) (*HttpRecord, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	record, ok := c.Records[provider][region]
	return record, ok
}

func (c *HttpConfig) Set(provider, region string, record *HttpRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.Records[provider]
	if !ok {
		c.Records[provider] = make(map[string]*HttpRecord)
	}
	c.Records[provider][region] = record
}

func (c *HttpConfig) Delete(provider, region string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Records[provider], region)
}

func (c *HttpConfig) Save() error {
	return save(c.Records, HttpProxyPath)
}

func (c *HttpConfig) AvailableApis() []string {
	var apis []string
	for _, rmap := range c.Records {
		for _, record := range rmap {
			r, ok := interface{}(record).(*HttpRecord)
			if !ok {
				return apis
			}
			if r.Api != "" {
				apis = append(apis, r.Api)
			}
		}
	}
	return apis
}

func (c *HttpConfig) ToDoubleArray() [][]string {
	data := [][]string{}
	for provider, rmap := range c.Records {
		for region, record := range rmap {
			data = append(data, []string{provider, region, record.Api})
		}
	}
	return data
}
