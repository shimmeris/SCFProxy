package huawei

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/functiongraph/v2/model"

	"github.com/shimmeris/SCFProxy/function"
	"github.com/shimmeris/SCFProxy/sdk"
)

const GroupName = "scf"

func (p *Provider) DeployHttpProxy(opts *sdk.HttpProxyOpts) (*sdk.DeployHttpProxyResult, error) {
	if err := p.createGroup(GroupName); err != nil {
		return nil, err
	}

	functionUrn, err := p.createFunction(opts.FunctionName)
	if err != nil {
		return nil, err
	}

	triggerId, err := p.createHttpTrigger(functionUrn, opts.TriggerName)
	if err != nil {
		return nil, err
	}

	return &sdk.DeployHttpProxyResult{API: triggerId, Region: p.region, Provider: p.Name()}, nil

}

func (p *Provider) createGroup(groupName string) error {
	createGroupApi := fmt.Sprintf("https://apig.%s.myhuaweicloud.com/v1.0/apigw/api-groups", p.region)
	data := fmt.Sprintf("{\"name\": \"%s\"}", groupName)
	req, err := http.NewRequest("POST", createGroupApi, strings.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if err = p.signer.Sign(req); err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
	defer resp.Body.Close()
	return nil

}

func (p *Provider) createFunction(functionName string) (string, error) {
	r := &model.CreateFunctionRequest{}
	r.Body = &model.CreateFunctionRequestBody{
		Package:    "default",
		FuncName:   functionName,
		Handler:    "index.handler",
		Timeout:    30,
		MemorySize: 128,
		CodeType:   model.GetCreateFunctionRequestBodyCodeTypeEnum().ZIP,
		Runtime:    model.GetCreateFunctionRequestBodyRuntimeEnum().PYTHON3_9,
		FuncCode: &model.FuncCode{
			File: &function.HuaweiHttpCodeZip,
		},
	}

	res, err := p.fclient.CreateFunction(r)
	if err != nil {
		return "", err
	}
	return *res.FuncUrn, nil
}

func (p *Provider) createHttpTrigger(functionUrn, triggerName string) (string, error) {
	r := &model.CreateFunctionTriggerRequest{}
	r.FunctionUrn = functionUrn
	triggerStatus := model.GetCreateFunctionTriggerRequestBodyTriggerStatusEnum().ACTIVE
	r.Body = &model.CreateFunctionTriggerRequestBody{
		TriggerTypeCode: model.GetCreateFunctionTriggerRequestBodyTriggerTypeCodeEnum().APIG,
		TriggerStatus:   &triggerStatus,
		EventData: map[string]string{
			"func_info":    "{timeout=5000}",
			"name":         triggerName,
			"env_id":       "DEFAULT_ENVIRONMENT_RELEASE_ID",
			"env_name":     "RELEASE",
			"protocol":     "HTTPS",
			"auth":         "NONE",
			"group_id":     "",
			"sl_domain":    "",
			"match_mode":   "SWA",
			"req_method":   "ANY",
			"backend_type": "FUNCTION",
			"type":         "1", //TODO: `type` must be `int`, but the sdk only suppors `string`, wait for Huawei to fix
			"path":         "/http",
		},
	}

	response, err := p.fclient.CreateFunctionTrigger(r)
	if err != nil {
		return "", err
	}
	return *response.TriggerId, nil
}
