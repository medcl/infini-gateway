package pipelines

import (
	"github.com/infinitbyte/framework/core/pipeline"
)

type LoggingJoint struct {
}

func (joint LoggingJoint) Name() string {
	return "logging"
}

func (joint LoggingJoint) Process(c *pipeline.Context) error {

	return nil
}
