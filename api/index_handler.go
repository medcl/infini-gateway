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
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/env"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/util"
	"github.com/julienschmidt/httprouter"
	"github.com/medcl/elasticsearch-proxy/config"
	"net/http"
	"time"
)

// IndexAction returns cluster health information
func (handler *API) IndexAction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	data := map[string]interface{}{}
	data["name"] = global.Env().SystemConfig.NodeConfig.Name

	version := map[string]interface{}{}
	version["number"] = config.Version
	version["build_commit"] = config.LastCommitLog
	version["build_date"] = config.BuildDate

	data["version"] = version
	data["tagline"] = "You Know, for Proxy"
	data["uptime"] = time.Since(env.GetStartTime()).String()

	handler.WriteJSON(w, &data, http.StatusOK)
}

func (handler *API) ProxyAction(w http.ResponseWriter, req *http.Request) {

	//handler.WriteJSON(w,util.MapStr{
	//	"upstream":handler.Config.Upstream,
	//},200)

	body, err := handler.GetRawBody(req)
	if err != nil {
		log.Error(err)
	}

	cfg := handler.Config.Upstream[0].Elasticsearch

	url := fmt.Sprintf("%s%s", cfg.Endpoint, req.URL)

	request := util.NewPostRequest(url, body)
	request.Method = req.Method
	request.SetBasicAuth(cfg.Username, cfg.Password)
	response, err := util.ExecuteRequest(request)
	if err != nil {
		log.Error(err)
	}
	if global.Env().IsDebug {
		log.Debug(url)
		log.Debug(req.Method)
		log.Debug(string(body))
		log.Debug(util.ToJson(req.URL, true))
		log.Debug("search response: ", string(body), ",", string(response.Body))
	}

	handler.Write(w, response.Body)
}
