package sdk

type DeployHttpProxyResult struct {
	API      string
	Region   string
	Provider string
}

type DeployReverseProxyResult struct {
	ServiceId     string
	ApiId         string
	PluginId      string
	ServiceDomain string
	Origin        string
	Region        string
	Provider      string
	Protocol      string
}
