package tencent

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/shimmeris/SCFProxy/sdk"
)

func (p *Provider) DeployReverseProxy(opts *sdk.ReverseProxyOpts) (*sdk.DeployReverseProxyResult, error) {
	serviceId, serviceDomain, err := p.createService(opts.Origin)
	if err != nil {
		return nil, err
	}

	apiId, err := p.createApi(serviceId, opts.Origin)
	if err != nil {
		return nil, err
	}

	if err = p.releaseService(serviceId); err != nil {
		return nil, err
	}

	result := &sdk.DeployReverseProxyResult{
		Provider:      p.Name(),
		Region:        p.region,
		ServiceId:     serviceId,
		ApiId:         apiId,
		ServiceDomain: serviceDomain,
		Origin:        opts.Origin,
	}

	if len(opts.Ips) == 0 {
		return result, nil
	}

	pluginId, err := p.createIPControlPlugin(opts.Ips)
	if err != nil {
		logrus.Errorf("create IPControl plugin failed for %s in %s.%s", opts.Origin, p.Name(), p.region)
		return result, nil
	}

	if err := p.attachPlugin(serviceId, apiId, pluginId); err != nil {
		logrus.Errorf("attach IPControl plugin failed for %s in %s.%s ", opts.Origin, p.Name(), p.region)
		return result, nil
	}
	result.PluginId = pluginId
	return result, nil
}

func (p *Provider) ClearReverseProxy(opts *sdk.ReverseProxyOpts) error {
	if opts.PluginId != "" {
		if err := p.deletePlugin(opts.PluginId); err != nil {
			if err, ok := err.(*errors.TencentCloudSDKError); !ok || err.Code != apigateway.RESOURCENOTFOUND_INVALIDPLUGIN {
				return err
			}
		}
	}

	if err := p.deleteApi(opts.ServiceId, opts.ApiId); err != nil {
		if err, ok := err.(*errors.TencentCloudSDKError); ok && err.Code == apigateway.RESOURCENOTFOUND_INVALIDSERVICE {
			return nil
		}
		return err
	}

	if err := p.unreleaseService(opts.ServiceId); err != nil {
		return err
	}

	return p.deleteService(opts.ServiceId)

}

func (p *Provider) createService(name string) (string, string, error) {
	r := apigateway.NewCreateServiceRequest()
	r.Protocol = common.StringPtr("http&https")
	r.ServiceName = common.StringPtr(name)

	resp, err := p.gclient.CreateService(r)
	if err != nil {
		return "", "", err
	}
	serviceId, serviceDomain := *resp.Response.ServiceId, *resp.Response.OuterSubDomain
	return serviceId, serviceDomain, nil
}

func (p *Provider) createApi(serviceId, origin string) (string, error) {
	protocol := "HTTP"
	method := "ANY"
	u, _ := url.Parse(origin)
	if u.Scheme == "ws" || u.Scheme == "wss" {
		protocol = "WEBSOCKET"
		method = "GET"
	}

	r := apigateway.NewCreateApiRequest()
	r.ApiName = common.StringPtr(strings.ToLower(protocol))
	r.AuthType = common.StringPtr("NONE")
	r.ResponseType = common.StringPtr("BINARY")
	r.ServiceType = common.StringPtr(protocol)
	r.ServiceTimeout = common.Int64Ptr(900)
	r.Protocol = common.StringPtr(protocol)
	r.ServiceId = common.StringPtr(serviceId)
	r.RequestConfig = &apigateway.ApiRequestConfig{
		Path:   common.StringPtr("/"),
		Method: common.StringPtr(method),
	}
	r.ServiceConfig = &apigateway.ServiceConfig{
		Url:    common.StringPtr(origin),
		Path:   common.StringPtr("/"),
		Method: common.StringPtr(method),
	}

	resp, err := p.gclient.CreateApi(r)
	if err != nil {
		return "", err
	}
	apiId := *resp.Response.Result.ApiId

	return apiId, nil
}

func (p *Provider) releaseService(serviceId string) error {
	r := apigateway.NewReleaseServiceRequest()
	r.ServiceId = common.StringPtr(serviceId)
	r.ReleaseDesc = common.StringPtr("")
	r.EnvironmentName = common.StringPtr("release")

	_, err := p.gclient.ReleaseService(r)
	return err
}

func (p *Provider) unreleaseService(serviceId string) error {
	r := apigateway.NewUnReleaseServiceRequest()
	r.ServiceId = common.StringPtr(serviceId)
	r.EnvironmentName = common.StringPtr("release")

	_, err := p.gclient.UnReleaseService(r)
	return err
}

func (p *Provider) deleteService(serviceId string) error {
	r := apigateway.NewDeleteServiceRequest()
	r.ServiceId = common.StringPtr(serviceId)

	_, err := p.gclient.DeleteService(r)
	return err
}

func (p *Provider) deleteApi(serviceId, apiId string) error {
	r := apigateway.NewDeleteApiRequest()
	r.ServiceId = common.StringPtr(serviceId)
	r.ApiId = common.StringPtr(apiId)

	_, err := p.gclient.DeleteApi(r)
	return err
}

func (p *Provider) createIPControlPlugin(ips []string) (string, error) {
	r := apigateway.NewCreatePluginRequest()
	r.PluginName = common.StringPtr("")
	r.PluginType = common.StringPtr("IPControl")
	r.Description = common.StringPtr(strings.Join(ips, "\n"))
	r.PluginData = common.StringPtr(generatePluginData(ips))

	// 返回的resp是一个CreatePluginResponse的实例，与请求对象对应
	resp, err := p.gclient.CreatePlugin(r)
	if err != nil {
		return "", err
	}

	return *resp.Response.Result.PluginId, nil
}

func (p *Provider) attachPlugin(serviceId, apiId, pluginId string) error {
	r := apigateway.NewAttachPluginRequest()
	r.PluginId = common.StringPtr(pluginId)
	r.ServiceId = common.StringPtr(serviceId)
	r.EnvironmentName = common.StringPtr("release")
	r.ApiIds = common.StringPtrs([]string{apiId})

	_, err := p.gclient.AttachPlugin(r)
	return err
}

func (p *Provider) deletePlugin(pluginId string) error {
	r := apigateway.NewDeletePluginRequest()

	r.PluginId = common.StringPtr(pluginId)

	// 返回的resp是一个DeletePluginResponse的实例，与请求对象对应
	_, err := p.gclient.DeletePlugin(r)
	return err
}

func generatePluginData(ips []string) string {
	ip := strings.Join(ips, "\\n")
	return fmt.Sprintf("{\"type\": \"white_list\", \"blocks\": \"%s\"}", ip)
}
