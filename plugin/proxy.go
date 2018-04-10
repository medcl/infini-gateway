package elasticsearch_proxy

import (
	. "github.com/infinitbyte/framework/core/config"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/medcl/elasticsearch-proxy/api"
	"github.com/medcl/elasticsearch-proxy/config"
	"github.com/medcl/elasticsearch-proxy/pipelines"
	"github.com/medcl/elasticsearch-proxy/ui"
)

type ProxyPlugin struct {
}

func (this ProxyPlugin) Name() string {
	return "Proxy"
}

var (
	proxyConfig = config.ProxyConfig{}
)

func (module ProxyPlugin) Start(cfg *Config) {

	cfg.Unpack(&proxyConfig)

	config.SetUpstream(proxyConfig.Upstream)

	//register UI
	if proxyConfig.UIEnabled {
		ui.InitUI()
	}

	api.InitAPI()

	//register pipeline joints
	pipeline.RegisterPipeJoint(pipelines.IndexJoint{})
	pipeline.RegisterPipeJoint(pipelines.LoggingJoint{})
}

func (module ProxyPlugin) Stop() error {
	return nil
}
