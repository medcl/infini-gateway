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

package api

import (
	"github.com/infinitbyte/framework/core/api"
	"github.com/medcl/elasticsearch-proxy/config"
)

// API namespace
type API struct {
	api.Handler
	Config config.ProxyConfig
}

// InitAPI register apis
func InitAPI(config config.ProxyConfig) {

	apis := API{Config: config}

	//Index
	api.HandleAPIMethod(api.GET, "/_proxy/", apis.IndexAction)

	//Stats APIs
	api.HandleAPIMethod(api.GET, "/_proxy/stats", apis.StatsAction)
	api.HandleAPIMethod(api.GET, "/_proxy/queue/stats", apis.QueueStatsAction)

	// Handle proxy
	api.HandleAPIFunc("/", apis.ProxyAction)

}
