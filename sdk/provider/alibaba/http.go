package alibaba

import (
	fcopen "github.com/alibabacloud-go/fc-open-20210406/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"

	"github.com/shimmeris/SCFProxy/function"
	"github.com/shimmeris/SCFProxy/sdk"
)

const ServiceName = "scf"

func (p *Provider) DeployHttpProxy(opts *sdk.HttpProxyOpts) (*sdk.DeployHttpProxyResult, error) {
	if err := p.createService(ServiceName); err != nil {
		if err, ok := err.(*tea.SDKError); ok && *err.StatusCode == 409 {
			logrus.Info("Service name already exists, will use existing")
		} else {
			return nil, err
		}
	}

	if err := p.createHttpFunction(ServiceName, opts.FunctionName); err != nil {
		if err, ok := err.(*tea.SDKError); ok && *err.StatusCode == 409 {
			logrus.Info("function already exists, will use existing")
		} else {
			return nil, err
		}
	}

	api, err := p.createHttpTrigger(ServiceName, opts.FunctionName, opts.TriggerName)
	if err != nil {
		return nil, err
	}
	//TODO: 创建触发器出错时，函数应该删除
	return &sdk.DeployHttpProxyResult{Provider: p.Name(), Region: p.region, API: api}, err
}

func (p *Provider) ClearHttpProxy(opts *sdk.HttpProxyOpts) error {
	return p.clearProxy(ServiceName, opts.FunctionName, opts.TriggerName, opts.OnlyTrigger)
}

func (p *Provider) createService(serviceName string) error {
	h := &fcopen.CreateServiceHeaders{}
	r := &fcopen.CreateServiceRequest{ServiceName: tea.String(serviceName)}
	_, err := p.fclient.CreateServiceWithOptions(r, h, p.runtime)
	return err
}

func (p *Provider) createHttpFunction(serviceName, functionName string) error {
	h := &fcopen.CreateFunctionHeaders{}
	r := &fcopen.CreateFunctionRequest{
		FunctionName: tea.String(functionName),
		Runtime:      tea.String("python3.9"),
		Handler:      tea.String("index.handler"),
		Timeout:      tea.Int32(30),
		MemorySize:   tea.Int32(128),
		Code: &fcopen.Code{
			ZipFile: tea.String(function.AlibabaHttpCodeZip),
		},
	}

	_, err := p.fclient.CreateFunctionWithOptions(tea.String(serviceName), r, h, p.runtime)
	return err
}

func (p *Provider) createHttpTrigger(serviceName, functionName, triggerName string) (string, error) {
	h := &fcopen.CreateTriggerHeaders{}
	r := &fcopen.CreateTriggerRequest{
		TriggerType:   tea.String("http"),
		TriggerName:   tea.String(triggerName),
		TriggerConfig: tea.String("{\"authType\": \"anonymous\", \"methods\": [\"POST\"]}"),
	}

	res, err := p.fclient.CreateTriggerWithOptions(tea.String(serviceName), tea.String(functionName), r, h, p.runtime)
	if err != nil {
		return "", err
	}
	return *res.Body.UrlInternet, nil
}
