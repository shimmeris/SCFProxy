package sdk

type Provider interface {
	Name() string
	Region() string
}

type HttpProxyProvider interface {
	Provider
	DeployHttpProxy(*HttpProxyOpts) (*DeployHttpProxyResult, error)
	ClearHttpProxy(*HttpProxyOpts) error
}

type SocksProxyProvider interface {
	Provider
	DeploySocksProxy(*SocksProxyOpts) error
	ClearSocksProxy(*SocksProxyOpts) error
}

type ReverseProxyProvider interface {
	Provider
	DeployReverseProxy(*ReverseProxyOpts) (*DeployReverseProxyResult, error)
	ClearReverseProxy(*ReverseProxyOpts) error
}
