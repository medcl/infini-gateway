package config

import (
	"github.com/infinitbyte/framework/core/index"
	"github.com/infinitbyte/framework/core/pipeline"
	"sync"
)

type UpstreamConfig struct {
	Name          string                    `config:"name"`
	QueueName     string                    `config:"queue_name"`
	MaxQueueDepth int64                     `config:"max_queue_depth"`
	Enabled       bool                      `config:"enabled"`
	Active        bool                      `config:"active"`
	Timeout       string                    `config:"timeout"`
	Elasticsearch index.ElasticsearchConfig `config:"elasticsearch"`
}

func (v *UpstreamConfig) SafeGetQueueName() string {
	queueName := v.QueueName
	if queueName == "" {
		queueName = v.Name
	}
	return queueName
}

type ProxyConfig struct {
	UIEnabled bool
	Upstream  []UpstreamConfig `config:"upstream"`
	Algorithm string
}

const Url pipeline.ParaKey = "url"
const Method pipeline.ParaKey = "method"
const Body pipeline.ParaKey = "body"
const Upstream pipeline.ParaKey = "upstream"
const Response pipeline.ParaKey = "response"
const ResponseSize pipeline.ParaKey = "response_size"
const ResponseStatusCode pipeline.ParaKey = "response_code"
const Message pipeline.ParaKey = "message"

//Bucket
const InactiveUpstream = "inactive_upstream"

var upstreams map[string]UpstreamConfig = map[string]UpstreamConfig{}

var l sync.RWMutex

func GetUpstreamConfig(key string) UpstreamConfig {
	l.RLock()
	defer l.RUnlock()
	v := upstreams[key]
	return v
}

func GetUpstreamConfigs() map[string]UpstreamConfig {
	return upstreams
}

func UpdateUpstreamConfigStatus(key string, active bool) {
	l.Lock()
	defer l.Unlock()
	v := upstreams[key]
	v.Active = active
	upstreams[key] = v
}

func SetUpstream(ups []UpstreamConfig) {
	l.Lock()
	defer l.Unlock()
	for _, v := range ups {
		//default Active is true
		v.Active = true

		//TODO get upstream status from DB, override active field
		upstreams[v.Name] = v
	}
}
