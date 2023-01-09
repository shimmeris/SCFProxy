package alibaba

import (
	fcopen "github.com/alibabacloud-go/fc-open-20210406/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/shimmeris/SCFProxy/function"
	"github.com/shimmeris/SCFProxy/sdk"
)

func (p *Provider) DeployHttpProxy(opts *sdk.FunctionOpts) (string, error) {
	if err := p.createService(opts.Namespace); err != nil {
		return "", err
	}

	if err := p.createHttpFunction(opts.Namespace, opts.FunctionName); err != nil {
		return "", err
	}

	api, err := p.createHttpTrigger(opts.Namespace, opts.FunctionName, opts.TriggerName)
	if err != nil {
		return "", err
	}
	return api, err
}

func (p *Provider) ClearHttpProxy(opts *sdk.FunctionOpts) error {
	if err := p.deleteTrigger(opts.Namespace, opts.FunctionName, opts.TriggerName); err != nil {
		return err
	}

	if opts.OnlyTrigger {
		return nil
	}

	return p.deleteFunction(opts.Namespace, opts.FunctionName)
}

func (p *Provider) createService(serviceName string) error {
	h := &fcopen.CreateServiceHeaders{}
	r := &fcopen.CreateServiceRequest{ServiceName: tea.String(serviceName)}
	_, err := p.fclient.CreateServiceWithOptions(r, h, p.runtime)
	if err != nil {
		if err, ok := err.(*tea.SDKError); !ok || *err.StatusCode != 409 {
			return err
		}
	}
	return nil
}

func (p *Provider) createHttpFunction(serviceName, functionName string) error {
	h := &fcopen.CreateFunctionHeaders{}
	r := &fcopen.CreateFunctionRequest{
		FunctionName: tea.String(functionName),
		Runtime:      tea.String("python3.9"),
		Handler:      tea.String("index.handler"),
		Timeout:      tea.Int32(10),
		MemorySize:   tea.Int32(128),
		Code: &fcopen.Code{
			ZipFile: tea.String(function.AlibabaHttpCodeZip),
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

func (p *Provider) deleteService(serviceName string) error {
	h := &fcopen.DeleteServiceHeaders{}
	_, err := p.fclient.DeleteServiceWithOptions(tea.String(serviceName), h, p.runtime)
	if err != nil {
		if err, ok := err.(*tea.SDKError); !ok || *err.StatusCode != 404 {
			return err
		}
	}
	return nil
}

func (p *Provider) deleteFunction(serviceName, functionName string) error {
	h := &fcopen.DeleteFunctionHeaders{}
	_, err := p.fclient.DeleteFunctionWithOptions(tea.String(serviceName), tea.String(functionName), h, p.runtime)
	if err != nil {
		if err, ok := err.(*tea.SDKError); !ok || *err.StatusCode != 404 {
			return err
		}
	}
	return nil
}

func (p *Provider) deleteTrigger(serviceName, functionName, triggerName string) error {
	deleteTriggerHeaders := &fcopen.DeleteTriggerHeaders{}
	_, err := p.fclient.DeleteTriggerWithOptions(
		tea.String(serviceName),
		tea.String(functionName),
		tea.String(triggerName),
		deleteTriggerHeaders,
		p.runtime,
	)
	if err != nil {
		if err, ok := err.(*tea.SDKError); !ok || *err.StatusCode != 404 {
			return err
		}
	}
	return nil
}
