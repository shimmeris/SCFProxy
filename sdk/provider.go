package sdk

type Provider interface {
	Name() string
	Region() string
}

type HttpProxyProvider interface {
	Provider
	DeployHttpProxy(*FunctionOpts) (string, error)
	ClearHttpProxy(*FunctionOpts) error
}

type SocksProxyProvider interface {
	Provider
	DeploySocksProxy(*FunctionOpts) error
	ClearSocksProxy(*FunctionOpts) error
	InvokeFunction(*FunctionOpts, string) error
}

type ReverseProxyProvider interface {
	Provider
	DeployReverseProxy(*ReverseProxyOpts) (*DeployReverseProxyResult, error)
	ClearReverseProxy(*ReverseProxyOpts) error
}
