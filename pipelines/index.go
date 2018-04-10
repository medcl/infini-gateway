package pipelines

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/queue"
	"github.com/infinitbyte/framework/core/util"
	"github.com/medcl/elasticsearch-proxy/config"
)

type IndexJoint struct {
}

func (joint IndexJoint) Name() string {
	return "index"
}

func (joint IndexJoint) Process(c *pipeline.Context) error {

	upstream := c.MustGetString(config.Upstream)

	cfg := config.GetUpstreamConfig(upstream)

	url := fmt.Sprintf("%s%s", cfg.Elasticsearch.Endpoint, c.MustGetString(config.Url))

	method := c.MustGetString(config.Method)
	request := util.NewRequest(method, url)

	body, ok := c.GetString(config.Body)

	if ok {
		request.Body = []byte(body)
	}

	request.SetBasicAuth(cfg.Elasticsearch.Username, cfg.Elasticsearch.Password)
	response, err := util.ExecuteRequest(request)

	if err != nil {
		log.Error(err)
		joint.handleError(c)
		return nil
	}

	if global.Env().IsDebug {
		log.Debug(upstream)
		log.Debug(url)
		log.Debug(method)
		log.Debug(body)
		log.Debug("response: ", body, ",", string(response.Body))
	}

	if response.StatusCode >= 400 {
		log.Error("response: ", body, ",", response.StatusCode, ",", string(response.Body))
		joint.handleError(c)
		return nil
	}

	c.Set(config.ResponseSize, response.Size)
	c.Set(config.ResponseStatusCode, response.StatusCode)
	c.Set(config.Response, response.Body)

	return nil
}

func (joint IndexJoint) handleError(c *pipeline.Context) {

	//TODO move to standard error pipeline process
	// handle error
	// stop ingestion, record the current request and error message
	// mark this upstream as inactive,
	// waiting for manual active, and manually redo the request

	upstream := c.MustGetString(config.Upstream)
	config.UpdateUpstreamConfigStatus(upstream, false)
	queue.PauseRead(upstream)
}
