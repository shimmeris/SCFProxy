package sdk

import (
	"encoding/base64"
	"encoding/json"
)

type HttpProxyOpts struct {
	FunctionName string
	TriggerName  string
	OnlyTrigger  bool
}

type SocksProxyOpts struct {
	FunctionName string
	TriggerName  string
	OnlyTrigger  bool
	Key          string
	Addr         string
	Auth         string
}

func (s *SocksProxyOpts) DumpBase64Message() string {
	message := struct {
		Key  string
		Addr string
		Auth string
	}{
		s.Key,
		s.Addr,
		s.Auth,
	}

	b, _ := json.Marshal(message)
	return base64.StdEncoding.EncodeToString(b)
}

type ReverseProxyOpts struct {
	Origin    string
	ServiceId string
	ApiId     string
	PluginId  string
	Ips       []string
}
