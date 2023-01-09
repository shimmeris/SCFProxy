package tencent

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/shimmeris/SCFProxy/function"
	"github.com/shimmeris/SCFProxy/sdk"
)

func (p *Provider) DeploySocksProxy(opts *sdk.FunctionOpts) error {
	if err := p.createNamespace(opts.Namespace); err != nil {
		return err
	}

	return p.createSocksFunction(opts.Namespace, opts.FunctionName)
}

func (p *Provider) ClearSocksProxy(opts *sdk.FunctionOpts) error {
	return p.deleteFunction(opts.Namespace, opts.FunctionName)
}

func (p *Provider) createSocksFunction(namespace, functionName string) error {
	r := scf.NewCreateFunctionRequest()
	r.Namespace = common.StringPtr(namespace)
	r.FunctionName = common.StringPtr(functionName)
	r.Handler = common.StringPtr("main")
	r.Runtime = common.StringPtr("Go1")
	r.Code = &scf.Code{ZipFile: common.StringPtr(function.TencentSocksCodeZip)}
	r.Timeout = common.Int64Ptr(900)
	r.MemorySize = common.Int64Ptr(128)

	_, err := p.fclient.CreateFunction(r)
	if err != nil {
		if err, ok := err.(*errors.TencentCloudSDKError); !ok || err.Code != scf.RESOURCEINUSE_FUNCTION {
			return err
		}
	}
	return nil
}

func (p *Provider) InvokeFunction(opts *sdk.FunctionOpts, message string) error {
	r := scf.NewInvokeRequest()
	r.Namespace = common.StringPtr(opts.Namespace)
	r.FunctionName = common.StringPtr(opts.FunctionName)
	r.InvocationType = common.StringPtr("Event")
	r.ClientContext = common.StringPtr(message)

	_, err := p.fclient.Invoke(r)
	return err

}
