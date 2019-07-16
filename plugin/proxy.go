package plugin

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
	proxyConfig = config.ProxyConfig{
		PassthroughPatterns: []string{
			"_search", "_count", "_analyze", "_mget",
			"_doc", "_mtermvectors", "_msearch", "_search_shards", "_suggest",
			"_validate", "_explain", "_field_caps", "_rank_eval", "_aliases",
			"_open", "_close"},
	}
)

func (module ProxyPlugin) Setup(cfg *Config) {
	cfg.Unpack(&proxyConfig)

	config.SetProxyConfig(proxyConfig)

	api.InitAPI()

	//register pipeline joints
	pipeline.RegisterPipeJoint(pipelines.IndexJoint{})
	pipeline.RegisterPipeJoint(pipelines.LoggingJoint{})

}

func (module ProxyPlugin) Start() error {

	//register UI
	if proxyConfig.UIEnabled {
		ui.InitUI()
	}

	return nil
}

func (module ProxyPlugin) Stop() error {
	return nil
}
