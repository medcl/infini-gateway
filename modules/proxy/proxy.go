package proxy

import (
	. "infini.sh/framework/core/config"
	"infini.sh/framework/core/pipeline"
	"infini.sh/proxy/api"
	"infini.sh/proxy/config"
	"infini.sh/proxy/pipelines"
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
	}

	return nil
}

func (module ProxyPlugin) Stop() error {
	return nil
}
