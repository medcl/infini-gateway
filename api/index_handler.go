package api

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/api/router"
	"github.com/infinitbyte/framework/core/elastic"
	"github.com/infinitbyte/framework/core/env"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/queue"
	"github.com/infinitbyte/framework/core/util"
	"github.com/medcl/elasticsearch-proxy/config"
	"github.com/medcl/elasticsearch-proxy/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// IndexAction returns cluster health information
func (handler *API) IndexAction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	upstream := handler.GetHeader(req, "UPSTREAM", "auto")
	if upstream != "auto" {
		log.Debug("parameter upstream: ", upstream)

		cfg := config.GetUpstreamConfig(upstream)
		if cfg.Enabled && cfg.Writeable {

			response, err := handler.executeHttpRequest(elastic.GetConfig(cfg.Elasticsearch), req.URL.String(), req.Method, nil)
			if err != nil {
				log.Error(err)

				handler.WriteJSON(w, util.MapStr{
					"error": err.Error(),
				}, 500)
				return
			}
			w.Header().Add("upstream", cfg.Name)
			w.WriteHeader(response.StatusCode)
			w.Write(response.Body)
			return
		} else {
			handler.WriteJSON(w, util.MapStr{
				"error": "upstram is not exist nor readable",
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
		cfg := elastic.GetConfig(v.Elasticsearch)
		if v.Enabled {
			m[v.Name] = util.MapStr{
				"endpoint":        cfg.Endpoint,
				"queue":           v.QueueName,
				"max_queue_depth": v.MaxQueueDepth,
				"readable":        v.Readable,
				"writeable":       v.Writeable,
				"timeout":         v.Timeout,
			}
		}
	}
	data["upstream"] = m

	handler.WriteJSON(w, &data, http.StatusOK)
}

func (handler *API) executeHttpRequest(cfg elastic.ElasticsearchConfig, url, method string, body []byte) (*util.Result, error) {
	url = fmt.Sprintf("%s%s", cfg.Endpoint, url)
	request := util.NewPostRequest(url, body)
	request.Method = method
	if cfg.BasicAuth != nil {
		request.SetBasicAuth(cfg.BasicAuth.Username, cfg.BasicAuth.Password)
	}
	return util.ExecuteRequest(request)
}

func getHash(keyPrefix string, req *http.Request, body []byte) string {
	return fmt.Sprintf("%s-%s", keyPrefix, util.MD5digest(fmt.Sprintf("%s-%s", req.URL, string(body))))
}

func (handler *API) handleRead(w http.ResponseWriter, req *http.Request, body []byte) {

	hash := getHash(handler.cacheConfig.KeyPrefix, req, body)

	if handler.cacheConfig.CacheEnabled {
		cache, _ := handler.redis.Get(hash).Result()
		if cache != "" {
			if global.Env().IsDebug {
				log.Trace("hit cache: ", req.URL, ",", cache)
			}
			w.Header().Add("upstream", "cache")
			w.WriteHeader(200)
			w.Write([]byte(cache))
			return
		}
	}

	upstream := handler.GetHeader(req, "UPSTREAM", "auto")
	if upstream != "auto" {
		log.Debug("parameter upstream: ", upstream)

		cfg := config.GetUpstreamConfig(upstream)
		if cfg.Enabled && cfg.Readable {

			response, err := handler.executeHttpRequest(elastic.GetConfig(cfg.Elasticsearch), req.URL.String(), req.Method, body)
			if err != nil {
				log.Error(err)

				request := model.Request{}
				request.Url = req.URL.String()
				request.Upstream = cfg.Name
				request.Method = req.Method
				request.Body = string(body)
				request.Message = err.Error()
				model.CreateRequest(&request)

				handler.WriteJSON(w, util.MapStr{
					"error": err.Error(),
				}, 500)
				return
			}
			w.Header().Add("upstream", cfg.Name)
			w.WriteHeader(response.StatusCode)
			w.Write(response.Body)

			if handler.cacheConfig.CacheEnabled {
				handler.redis.Set(hash, string(response.Body), *handler.cacheConfig.GetTTLDuration()).Err()
				if global.Env().IsDebug {
					log.Debug("update cache: ", hash)
				}
			}

			return
		} else {
			handler.WriteJSON(w, util.MapStr{
				"error": "upstram is not exist nor readable",
			}, 500)
			return
		}
	}

	ups := config.GetUpstreamConfigs()
	for _, v := range ups {
		if v.Enabled && v.Readable {

			response, err := handler.executeHttpRequest(elastic.GetConfig(v.Elasticsearch), req.URL.String(), req.Method, body)

			if err != nil {
				log.Error(err)

				//v.Readable = false
				//log.Error("upstream: ",v.Name, ", disabled")

				request := model.Request{}
				request.Url = req.URL.String()
				request.Upstream = v.Name
				request.Method = req.Method
				request.Body = string(body)
				request.Message = err.Error()
				model.CreateRequest(&request)

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

			if handler.cacheConfig.CacheEnabled {
				handler.redis.Set(hash, string(response.Body), *handler.cacheConfig.GetTTLDuration()).Err()
				if global.Env().IsDebug {
					log.Debug("update cache: ", hash)
				}
			}

			return
		}
	}

	handler.WriteJSON(w, util.MapStr{
		"error": noUpstreamMsg,
	}, 500)

}

// POST should not used to serve as search/read/ requests
func (handler *API) handleWrite(w http.ResponseWriter, req *http.Request, body []byte) {
	url := fmt.Sprintf("%s", req.URL)

	//TODO add HEADER to support read through and write through

	//indexing/bulk
	//_bulk
	//_delete_by_query?
	//_update_by_query?
	//_reindex?
	//xxx?refresh=yy
	if util.ContainsAnyInArray(url, config.GetProxyConfig().PassthroughPatterns) ||
		(handler.GetParameter(req, "refresh") != "") {
		handler.handleRead(w, req, body)
		return
	}

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

			context := pipeline.Context{}
			context.Set(config.Upstream, v.Name)
			context.Set(config.Url, url)
			context.Set(config.Method, req.Method)
			context.Set(config.Body, string(body))

			queue.Push(v.SafeGetQueueName(), util.ToJSONBytes(context))
			response[v.Name] = "success"
		}
	}

	code := 200
	if !ack {
		code = 500
	}

	handler.WriteJSON(w, util.MapStr{
		"acknowledge": ack,
		"_upstream":   response,
	}, code)
}

var noUpstreamMsg = "no upstream available"

func (handler *API) ProxyAction(w http.ResponseWriter, req *http.Request) {

	handler.WriteJSONHeader(w)

	body, err := handler.GetRawBody(req)
	if err != nil {
		log.Error(err)
		handler.WriteJSON(w, util.MapStr{
			"error": err.Error(),
		}, 500)
	}

	if global.Env().IsDebug {
		log.Debug(req.URL)
		log.Debug(req.Method)
		log.Debug(string(body))
		log.Debug("request: ", string(body))
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
		methodNotAllowed := fmt.Sprintf("method %s is not supported", req.Method)
		request := model.Request{}
		request.Url = req.URL.String()
		request.Method = req.Method
		request.Body = string(body)
		request.Message = err.Error()
		model.CreateRequest(&request)

		handler.WriteJSON(w, util.MapStr{
			"error": methodNotAllowed,
		}, 500)
		return
	}

}

func (handler *API) GetRequestsAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	fr := handler.GetParameter(req, "from")
	si := handler.GetParameter(req, "size")
	upstream := handler.GetParameter(req, "upstream")
	status := handler.GetIntOrDefault(req, "status", -1)

	from, err := strconv.Atoi(fr)
	if err != nil {
		from = 0
	}
	size, err := strconv.Atoi(si)
	if err != nil {
		size = 10
	}

	total, tasks, err := model.GetRequestList(from, size, upstream, status)
	if err != nil {
		handler.WriteJSON(w, util.MapStr{
			"error": err.Error(),
		}, 500)
	} else {
		handler.WriteJSONListResult(w, total, tasks, http.StatusOK)
	}
}

//curl  -XPOST http://localhost:2900/_proxy/request/redo -d'{"ids":["bb6t4cqaukihf1ag10q0","bb6t4daaukihf1ag10r0"]}'
//{
//"acknowledge": true,
//"result": {
//"bb6t4cqaukihf1ag10q0": "{\"_index\":\"myindex\",\"_type\":\"doc\",\"_id\":\"1\",\"_version\":17,\"result\":\"updated\",\"_shards\":{\"total\":2,\"successful\":1,\"failed\":0},\"_seq_no\":16,\"_primary_term\":2}",
//"bb6t4daaukihf1ag10r0": "{\"_index\":\"myindex\",\"_type\":\"doc\",\"_id\":\"1\",\"_version\":18,\"result\":\"updated\",\"_shards\":{\"total\":2,\"successful\":1,\"failed\":0},\"_seq_no\":17,\"_primary_term\":2}"
//}
//}
func (handler *API) RedoRequestsAction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	//TODO check status, add `force` parameter to force execute the replay
	json, err := handler.GetJSON(req)
	if err != nil {
		log.Error(err)
		handler.WriteJSON(w, util.MapStr{
			"error": err.Error(),
		}, 500)
		return
	}

	ids, err := json.ArrayOfStrings("ids")
	if err != nil {
		log.Error(err)
		handler.WriteJSON(w, util.MapStr{
			"error": err.Error(),
		}, 500)
		return
	}
	ack := true
	msg := util.MapStr{}
	for _, id := range ids {
		request, err := model.GetRequest(id)
		if err != nil {
			log.Error(err)
			ack = false
			msg[id] = err.Error()
			continue
		}

		//replay request
		cfg := config.GetUpstreamConfig(request.Upstream)
		result, err := handler.executeHttpRequest(elastic.GetConfig(cfg.Elasticsearch), request.Url, request.Method, []byte(request.Body))

		//update request status
		request.Status = model.ReplayedSuccess
		request.Updated = time.Now()
		request.Response = string(result.Body)
		request.ResponseSize = int64(result.Size)
		request.ResponseStatusCode = result.StatusCode
		msg[id] = request.Response

		if err != nil {
			request.Status = model.ReplayedFailure
			request.Message = err.Error()
			ack = false
			msg[id] = err.Error()
		}

		model.UpdateRequest(&request)
	}

	handler.WriteJSON(w, util.MapStr{
		"acknowledge": ack,
		"result":      msg,
	}, 500)
}
