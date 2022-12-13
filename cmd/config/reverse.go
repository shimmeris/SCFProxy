package config

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type ReverseRecord struct {
	Provider  string
	Region    string
	Origin    string
	Api       string
	ServiceId string
	ApiId     string
	PluginId  string
	Ips       []string
}

type ReverseConfig struct {
	mu      sync.Mutex
	Records []*ReverseRecord
}

func LoadReverseConfig() (*ReverseConfig, error) {
	config := &ReverseConfig{}
	data, err := os.ReadFile(ReverseProxyPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return config, nil
		}
		return nil, err
	}

	if err = json.Unmarshal(data, &config.Records); err != nil {
		return nil, err
	}
	return config, nil
}

func (r *ReverseConfig) Add(record *ReverseRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Records = append(r.Records, record)
}

func (r *ReverseConfig) Save() error {
	return save(r.Records, ReverseProxyPath)
}

func (r *ReverseConfig) ToDoubleArray() [][]string {
	data := [][]string{}
	for _, r := range r.Records {
		data = append(data, []string{r.Provider, r.Region, r.Origin, r.Api})
	}
	return data
}
