package pipelines

import (
	"infini.sh/framework/core/pipeline"
	"infini.sh/proxy/config"
	"infini.sh/proxy/model"
)

type LoggingJoint struct {
}

func (joint LoggingJoint) Name() string {
	return "logging"
}

func (joint LoggingJoint) Process(c *pipeline.Context) error {

	request := model.Request{}
	request.Url = c.MustGetString(config.Url)
	request.Upstream = c.GetStringOrDefault(config.Upstream, "")
	request.Method = c.MustGetString(config.Method)
	request.Body = c.GetStringOrDefault(config.Body, "")
	request.Message = c.GetStringOrDefault(config.Message, "")

	if c.Has(config.ResponseStatusCode) {
		request.ResponseStatusCode = c.MustGetInt(config.ResponseStatusCode)
	}
	request.Response = c.GetStringOrDefault(config.Response, "")
	request.ResponseSize = c.GetInt64OrDefault(config.ResponseSize, 0)
	model.CreateRequest(&request)
	return nil
}
