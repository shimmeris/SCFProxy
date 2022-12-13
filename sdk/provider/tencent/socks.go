package tencent

import (
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/shimmeris/SCFProxy/fileutil"
	"github.com/shimmeris/SCFProxy/function"
	"github.com/shimmeris/SCFProxy/sdk"
)

func (p *Provider) DeploySocksProxy(opts *sdk.SocksProxyOpts) error {
	if !opts.OnlyTrigger {
		if err := p.createSocksFunction(opts.FunctionName); err != nil {
			if err, ok := err.(*errors.TencentCloudSDKError); !ok || err.Code != scf.RESOURCEINUSE_FUNCTION {
				return err
			}
		}
		time.Sleep(10 * time.Second)
	}

	message := opts.DumpBase64Message()

	var err error
	// tencent returns async. retry 3 times
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		err = p.createSocksTrigger(opts.FunctionName, opts.TriggerName, message)
		if err == nil {
			break
		}
	}
	return err
}

func (p *Provider) ClearSocksProxy(opts *sdk.SocksProxyOpts) error {
	return p.clearFunctionProxy(opts.FunctionName, opts.TriggerName, "timer", opts.OnlyTrigger)
}

func (p *Provider) createSocksFunction(functionName string) error {
	r := scf.NewCreateFunctionRequest()
	r.FunctionName = common.StringPtr(functionName)
	r.Handler = common.StringPtr("main")
	r.Runtime = common.StringPtr("Go1")
	r.Code = &scf.Code{ZipFile: common.StringPtr(fileutil.CreateZipBase64("main", function.TencentSocksCode))}
	r.Timeout = common.Int64Ptr(900)
	r.MemorySize = common.Int64Ptr(128)

	_, err := p.fclient.CreateFunction(r)
	return err

}

func (p *Provider) createSocksTrigger(functionName, triggerName, message string) error {
	r := scf.NewCreateTriggerRequest()
	r.FunctionName = common.StringPtr(functionName)
	r.TriggerName = common.StringPtr(triggerName)
	r.Type = common.StringPtr("timer")
	r.TriggerDesc = common.StringPtr("0 */1 * * * * *") // every 1 min
	r.CustomArgument = common.StringPtr(message)

	_, err := p.fclient.CreateTrigger(r)
	return err
}
