package sdk

type FunctionOpts struct {
	Namespace    string
	FunctionName string
	TriggerName  string
	OnlyTrigger  bool
}

type ReverseProxyOpts struct {
	Origin    string
	ServiceId string
	ApiId     string
	PluginId  string
	Ips       []string
}
