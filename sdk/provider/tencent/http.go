package tencent

import (
	"encoding/json"
	"time"

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

func (p *Provider) DeployHttpProxy(opts *sdk.HttpProxyOpts) (*sdk.DeployHttpProxyResult, error) {
	if !opts.OnlyTrigger {
		if err := p.createHttpFunction(opts.FunctionName); err != nil {
			if err, ok := err.(*errors.TencentCloudSDKError); !ok || err.Code != scf.RESOURCEINUSE_FUNCTION {
				return nil, err
			}
		}
	}

	var api string
	var err error
	// tencent returns async. retry 3 times
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		api, err = p.createHttpTrigger(opts.FunctionName, opts.TriggerName)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	result := &sdk.DeployHttpProxyResult{
		Provider: p.Name(),
		Region:   p.region,
		API:      api,
	}
	return result, nil

}

func (p *Provider) ClearHttpProxy(opts *sdk.HttpProxyOpts) error {
	return p.clearFunctionProxy(opts.FunctionName, opts.TriggerName, "apigw", opts.OnlyTrigger)
}

func (p *Provider) createHttpFunction(functionName string) error {
	r := scf.NewCreateFunctionRequest()
	r.FunctionName = common.StringPtr(functionName)
	r.Code = &scf.Code{ZipFile: common.StringPtr(function.TencentHttpCodeZip)}
	r.Handler = common.StringPtr("index.handler")
	r.MemorySize = common.Int64Ptr(128)
	r.Timeout = common.Int64Ptr(30)
	r.Runtime = common.StringPtr("Python3.6")

	_, err := p.fclient.CreateFunction(r)
	return err
}

func (p *Provider) createHttpTrigger(functionName, triggerName string) (string, error) {
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
