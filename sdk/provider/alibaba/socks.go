package alibaba

import (
	"fmt"

	fcopen "github.com/alibabacloud-go/fc-open-20210406/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/shimmeris/SCFProxy/function"
	"github.com/shimmeris/SCFProxy/sdk"
)

func (p *Provider) DeploySocksProxy(opts *sdk.SocksProxyOpts) error {
	if err := p.createService(ServiceName); err != nil {
		if err, ok := err.(*tea.SDKError); !ok || *err.StatusCode != 409 {
			return err
		}
	}

	if err := p.createSocksFunction(ServiceName, opts.FunctionName); err != nil {
		if err, ok := err.(*tea.SDKError); !ok || *err.StatusCode != 409 {
			return err
		}
	}
	// TODO: 创建触发器出错时，函数应该删除
	return p.createSocksTrigger(ServiceName, opts.FunctionName, opts.TriggerName, opts.DumpBase64Message())

}

func (p *Provider) ClearSocksProxy(opts *sdk.SocksProxyOpts) error {
	return p.clearProxy(ServiceName, opts.FunctionName, opts.TriggerName, opts.OnlyTrigger)
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
	return err
}

func (p *Provider) createSocksTrigger(serviceName, functionName, triggerName, message string) error {
	h := &fcopen.CreateTriggerHeaders{}
	r := &fcopen.CreateTriggerRequest{
		TriggerType:   tea.String("timer"),
		TriggerName:   tea.String(triggerName),
		TriggerConfig: tea.String(fmt.Sprintf("{\"enable\": true,  \"cronExpression\": \"@every 1m\",  \"payload\": \"%s\"}", message)),
	}

	_, err := p.fclient.CreateTriggerWithOptions(tea.String(serviceName), tea.String(functionName), r, h, p.runtime)
	return err
}
