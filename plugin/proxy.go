/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package elasticsearch_proxy

import (
	. "github.com/infinitbyte/framework/core/config"
	"github.com/medcl/elasticsearch-proxy/api"
	"github.com/medcl/elasticsearch-proxy/config"
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

	//register UI
	if proxyConfig.UIEnabled {
		ui.InitUI()
	}
	api.InitAPI(proxyConfig)
}

func (module ProxyPlugin) Stop() error {
	return nil
}
