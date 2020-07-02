package pipelines

import (
	"fmt"
	log "github.com/cihub/seelog"
	"infini.sh/framework/core/elastic"
	"infini.sh/framework/core/errors"
	"infini.sh/framework/core/filter"
	"infini.sh/framework/core/global"
	"infini.sh/framework/core/pipeline"
	"infini.sh/framework/core/queue"
	"infini.sh/framework/core/util"
	"infini.sh/proxy/config"
	"infini.sh/proxy/model"
)

type IndexJoint struct {
}

func (joint IndexJoint) Name() string {
	return "index"
}

func (joint IndexJoint) Process(c *pipeline.Context) error {

	upstream := c.MustGetString(config.Upstream)

	cfg := config.GetUpstreamConfig(upstream)

	esConfig := elastic.GetConfig(cfg.Elasticsearch)

	url := fmt.Sprintf("%s%s", esConfig.Endpoint, c.MustGetString(config.Url))

	method := c.MustGetString(config.Method)
	request := util.NewRequest(method, url)

	body, ok := c.GetString(config.Body)

	if ok && body != "" {
		request.Body = []byte(body)
	}

	if esConfig.BasicAuth != nil {
		request.SetBasicAuth(esConfig.BasicAuth.Username, esConfig.BasicAuth.Password)
	}

	request.SetContentType(util.ContentTypeJson)

	response, err := util.ExecuteRequest(request)

	if err != nil {
		log.Error(err)
		joint.handleError(c, err)
		return nil
	}

	if global.Env().IsDebug {
		log.Debug(upstream)
		log.Debug(url)
		log.Debug(method)
		log.Debug(body)
		log.Debug("response: ", body, ",", string(response.Body))
	}

	c.Set(config.ResponseSize, response.Size)
	c.Set(config.ResponseStatusCode, response.StatusCode)
	c.Set(config.Response, response.Body)

	if response.StatusCode >= 400 {
		err := errors.Errorf("response:%s, %v, %s ", body, response.StatusCode, string(response.Body))
		log.Error(err)
		joint.handleError(c, err)
		return nil
	}

	return nil
}

func (joint IndexJoint) handleError(c *pipeline.Context, err error) {

	//TODO move to standard error pipeline process
	// handle error
	// stop ingestion, record the current request and error message
	// mark this upstream as inactive,
	// waiting for manual active, and manually redo the request

	if c.Has(config.Upstream) {
		upstream := c.MustGetString(config.Upstream)
		filter.Add(config.InactiveUpstream, []byte(upstream))
		config.UpdateUpstreamWriteableStatus(upstream, false)
		queue.PauseRead(upstream)
	}

	c.Set(config.Message, err.Error())

	//save msg, TODO remove below, use logging joint to process the save process
	request := model.Request{}
	request.Status = model.Created
	request.Url = c.MustGetString(config.Url)
	request.Upstream = c.MustGetString(config.Upstream)
	request.Method = c.MustGetString(config.Method)
	request.Body = c.GetStringOrDefault(config.Body, "")
	request.Message = c.GetStringOrDefault(config.Message, "")
	if c.Has(config.ResponseStatusCode) {
		request.ResponseStatusCode = c.MustGetInt(config.ResponseStatusCode)
	}
	request.ResponseSize = c.GetInt64OrDefault(config.ResponseSize, 0)
	request.Response = c.GetStringOrDefault(config.Response, "")
	model.CreateRequest(&request)

}
