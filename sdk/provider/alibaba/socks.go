package alibaba

import (
	fcopen "github.com/alibabacloud-go/fc-open-20210406/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/shimmeris/SCFProxy/function"
	"github.com/shimmeris/SCFProxy/sdk"
)

func (p *Provider) DeploySocksProxy(opts *sdk.FunctionOpts) error {
	if err := p.createService(opts.Namespace); err != nil {
		return err
	}
	return p.createSocksFunction(opts.Namespace, opts.FunctionName)
}

func (p *Provider) ClearSocksProxy(opts *sdk.FunctionOpts) error {
	return p.deleteFunction(opts.Namespace, opts.FunctionName)
}

func (p *Provider) createSocksFunction(serviceName, functionName string) error {
	h := &fcopen.CreateFunctionHeaders{}
	r := &fcopen.CreateFunctionRequest{
		FunctionName: tea.String(functionName),
		Runtime:      tea.String("go1"),
		Handler:      tea.String("main"),
		Timeout:      tea.Int32(900),
		MemorySize:   tea.Int32(128),
		Code: &fcopen.Code{
			ZipFile: tea.String(function.AlibabaSocksCodeZip),
		},
	}

	_, err := p.fclient.CreateFunctionWithOptions(tea.String(serviceName), r, h, p.runtime)
	if err != nil {
		if err, ok := err.(*tea.SDKError); !ok || *err.StatusCode != 409 {
			return err
		}
	}
	return nil
}

func (p *Provider) InvokeFunction(opts *sdk.FunctionOpts, message string) error {
	h := &fcopen.InvokeFunctionHeaders{XFcInvocationType: tea.String("Async")}
	r := &fcopen.InvokeFunctionRequest{Body: []byte(message)}

	_, err := p.fclient.InvokeFunctionWithOptions(tea.String(opts.Namespace), tea.String(opts.FunctionName), r, h, p.runtime)
	return err
}
