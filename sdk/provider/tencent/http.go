package tencent

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/shimmeris/SCFProxy/function"
	"github.com/shimmeris/SCFProxy/sdk"
)

type ApiExtractor struct {
	Service struct {
		SubDomain string `json:"subDomain"`
	} `json:"service"`
}

func (p *Provider) DeployHttpProxy(opts *sdk.FunctionOpts) (string, error) {
	if err := p.createNamespace(opts.Namespace); err != nil {
		return "", err
	}

	if !opts.OnlyTrigger {
		if err := p.createHttpFunction(opts.Namespace, opts.FunctionName); err != nil {
			return "", err
		}
	}

	var api string
	var err error
	// tencent returns async. retry 3 times
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		api, err = p.createHttpTrigger(opts.Namespace, opts.FunctionName, opts.TriggerName)
		if err == nil {
			break
		}
		logrus.Errorf("Failed creating http proxy function in tencent.%s, retry after 10 sec", p.region)
	}
	if err != nil {
		return "", err
	}

	return api, nil

}

func (p *Provider) ClearHttpProxy(opts *sdk.FunctionOpts) error {
	if err := p.deleteTrigger(opts.Namespace, opts.FunctionName, opts.TriggerName, "apigw"); err != nil {
		return err
	}

	if opts.OnlyTrigger {
		return nil
	}

	return p.deleteFunction(opts.Namespace, opts.FunctionName)

}

func (p *Provider) createNamespace(namespace string) error {
	r := scf.NewCreateNamespaceRequest()
	r.Namespace = common.StringPtr(namespace)

	_, err := p.fclient.CreateNamespace(r)
	if err != nil {
		if err, ok := err.(*errors.TencentCloudSDKError); !ok || err.Code != scf.RESOURCEINUSE_NAMESPACE {
			return err
		}
	}
	return nil
}

func (p *Provider) createHttpFunction(namespace, functionName string) error {
	r := scf.NewCreateFunctionRequest()
	r.Namespace = common.StringPtr(namespace)
	r.FunctionName = common.StringPtr(functionName)
	r.Code = &scf.Code{ZipFile: common.StringPtr(function.TencentHttpCodeZip)}
	r.Handler = common.StringPtr("index.handler")
	r.MemorySize = common.Int64Ptr(128)
	r.Timeout = common.Int64Ptr(30)
	r.Runtime = common.StringPtr("Python3.6")

	_, err := p.fclient.CreateFunction(r)
	if err != nil {
		if err, ok := err.(*errors.TencentCloudSDKError); !ok || err.Code != scf.RESOURCEINUSE_FUNCTION {
			return err
		}
	}
	return nil
}

func (p *Provider) createHttpTrigger(namespace, functionName, triggerName string) (string, error) {
	r := scf.NewCreateTriggerRequest()
	r.FunctionName = common.StringPtr(functionName)
	r.TriggerName = common.StringPtr(triggerName)
	r.Type = common.StringPtr("apigw")
	r.TriggerDesc = common.StringPtr(`{
				"api":{
					"authRequired":"FALSE",
					"requestConfig":{
						"method":"POST"
					},
					"isIntegratedResponse":"TRUE"
				},
				"service":{
					"serviceName":"SCF_API_SERVICE"
				},
				"release":{
					"environmentName":"release"
				}
			}`)
	r.Namespace = common.StringPtr(namespace)

	response, err := p.fclient.CreateTrigger(r)
	if err != nil {
		return "", err
	}

	extractor := &ApiExtractor{}
	desc := *response.Response.TriggerInfo.TriggerDesc
	if err := json.Unmarshal([]byte(desc), extractor); err != nil {
		return "", err
	}

	api := extractor.Service.SubDomain
	return api, nil
}

func (p *Provider) deleteNamespace(namespace string) error {
	r := scf.NewDeleteNamespaceRequest()
	r.Namespace = common.StringPtr(namespace)

	_, err := p.fclient.DeleteNamespace(r)
	if err != nil {
		if err, ok := err.(*errors.TencentCloudSDKError); !ok || err.Code != scf.RESOURCENOTFOUND_NAMESPACE {
			return err
		}
	}
	return nil
}

func (p *Provider) deleteFunction(namespace, functionName string) error {
	r := scf.NewDeleteFunctionRequest()
	r.Namespace = common.StringPtr(namespace)
	r.FunctionName = common.StringPtr(functionName)

	_, err := p.fclient.DeleteFunction(r)
	if err != nil {
		if err, ok := err.(*errors.TencentCloudSDKError); !ok || (err.Code != scf.RESOURCENOTFOUND_NAMESPACE && err.Code != scf.RESOURCENOTFOUND_FUNCTION) {
			return err
		}
	}
	return nil
}

func (p *Provider) deleteTrigger(namespace, functionName, triggerName, triggerType string) error {
	r := scf.NewDeleteTriggerRequest()
	r.Namespace = common.StringPtr(namespace)
	r.FunctionName = common.StringPtr(functionName)
	r.TriggerName = common.StringPtr(triggerName)
	r.Type = common.StringPtr(triggerType)

	_, err := p.fclient.DeleteTrigger(r)
	if err != nil {
		if err, ok := err.(*errors.TencentCloudSDKError); !ok || (err.Code != scf.RESOURCENOTFOUND && err.Code != scf.RESOURCENOTFOUND_FUNCTION) {
			return err
		}
	}
	return nil
}
