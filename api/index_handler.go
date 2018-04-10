package api

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/api/router"
	"github.com/infinitbyte/framework/core/env"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/index"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/queue"
	"github.com/infinitbyte/framework/core/util"
	"github.com/medcl/elasticsearch-proxy/config"
	"net/http"
	"strings"
	"time"
)

// IndexAction returns cluster health information
func (handler *API) IndexAction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	upstream := handler.GetHeader(req, "UPSTREAM", "auto")
	if upstream != "auto" {
		log.Debug("parameter upstream: ", upstream)

		cfg := config.GetUpstreamConfig(upstream)
		if cfg.Enabled && cfg.Active {

			response, err := handler.executeHttpRequest(cfg.Elasticsearch, req, nil)
			if err != nil {
				handler.WriteJSON(w, util.MapStr{
					"error": err,
				}, 500)
				return
			}
			w.Header().Add("upstream", cfg.Name)
			w.WriteHeader(response.StatusCode)
			w.Write(response.Body)
			return
		} else {
			handler.WriteJSON(w, util.MapStr{
				"error": "upstram is not exist nor active",
			}, 500)
			return
		}
	}

	data := map[string]interface{}{}
	data["name"] = global.Env().SystemConfig.NodeConfig.Name

	version := map[string]interface{}{}
	version["number"] = util.TrimSpaces(config.Version)
	version["build_commit"] = util.TrimSpaces(config.LastCommitLog)
	version["build_date"] = strings.TrimSpace(config.BuildDate)

	data["version"] = version
	data["tagline"] = "You Know, for Proxy"
	data["uptime"] = time.Since(env.GetStartTime()).String()

	ups := config.GetUpstreamConfigs()
	m := util.MapStr{}
	for _, v := range ups {
		if v.Enabled {
			m[v.Name] = v.Elasticsearch.Endpoint
		}
	}
	data["upstream"] = m

	handler.WriteJSON(w, &data, http.StatusOK)
}

func (handler *API) executeHttpRequest(cfg index.ElasticsearchConfig, req *http.Request, body []byte) (*util.Result, error) {
	url := fmt.Sprintf("%s%s", cfg.Endpoint, req.URL)
	request := util.NewPostRequest(url, body)
	request.Method = req.Method
	request.SetBasicAuth(cfg.Username, cfg.Password)
	return util.ExecuteRequest(request)
}

func (handler *API) handleRead(w http.ResponseWriter, req *http.Request, body []byte) {
	upstream := handler.GetHeader(req, "UPSTREAM", "auto")
	if upstream != "auto" {
		log.Debug("parameter upstream: ", upstream)

		cfg := config.GetUpstreamConfig(upstream)
		if cfg.Enabled && cfg.Active {

			response, err := handler.executeHttpRequest(cfg.Elasticsearch, req, body)
			if err != nil {
				handler.WriteJSON(w, util.MapStr{
					"error": err,
				}, 500)
				return
			}
			w.Header().Add("upstream", cfg.Name)
			w.WriteHeader(response.StatusCode)
			w.Write(response.Body)
			return
		} else {
			handler.WriteJSON(w, util.MapStr{
				"error": "upstram is not exist nor active",
			}, 500)
			return
		}
	}

	ups := config.GetUpstreamConfigs()
	for _, v := range ups {
		if v.Enabled && v.Active {

			cfg := v.Elasticsearch
			response, err := handler.executeHttpRequest(cfg, req, body)

			if err != nil {
				log.Error(err)
				continue
			}

			if global.Env().IsDebug {
				log.Debug(req.URL)
				log.Debug(req.Method)
				log.Debug(string(body))
				log.Debug("search response: ", string(body), ",", string(response.Body))
			}

			w.Header().Add("upstream", v.Name)
			w.WriteHeader(response.StatusCode)
			w.Write(response.Body)

			return
		}
	}

	handler.WriteJSON(w, util.MapStr{
		"error": noUpstreamMsg,
	}, 500)

}

func (handler *API) handleWrite(w http.ResponseWriter, req *http.Request, body []byte) {
	response := map[string]string{}
	ack := true
	ups := config.GetUpstreamConfigs()
	for _, v := range ups {
		if v.Enabled {

			if v.MaxQueueDepth > 0 {
				depth := queue.Depth(v.QueueName)
				if depth >= v.MaxQueueDepth {
					response[v.Name] = "reach to maximum queue depth, the message was rejected"
					ack = false
					continue
				}
			}

			url := fmt.Sprintf("%s", req.URL)
			context := pipeline.Context{}
			context.Set(config.Upstream, v.Name)
			context.Set(config.Url, url)
			context.Set(config.Method, req.Method)
			context.Set(config.Body, string(body))

			queue.Push(v.SafeGetQueueName(), util.ToJSONBytes(context))
			response[v.Name] = "success"
		}
	}

	handler.WriteJSON(w, util.MapStr{
		"acknowledge": ack,
		"_upstream":   response,
	}, 200)
}

var noUpstreamMsg = "no upstream available"

func (handler *API) ProxyAction(w http.ResponseWriter, req *http.Request) {

	body, err := handler.GetRawBody(req)
	if err != nil {
		handler.WriteJSON(w, util.MapStr{
			"error": err,
		}, 500)
	}

	ups := config.GetUpstreamConfigs()
	if len(ups) == 0 {
		handler.WriteJSON(w, util.MapStr{
			"error": noUpstreamMsg,
		}, 500)
		return
	}

	switch req.Method {
	case "GET":
		handler.handleRead(w, req, body)
		break
	case "POST":
		handler.handleWrite(w, req, body)
		break
	case "PUT":
		handler.handleWrite(w, req, body)
		break
	case "DELETE":
		handler.handleWrite(w, req, body)
		break
	default:
		handler.WriteJSON(w, util.MapStr{
			"error": fmt.Sprintf("method %s is not supported", req.Method),
		}, 200)
		return
	}

}
