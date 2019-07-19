package api

import (
	log "github.com/cihub/seelog"
	"github.com/golang/go/src/pkg/fmt"
	"github.com/infinitbyte/framework/core/api"
	"github.com/infinitbyte/framework/core/env"
	"src/github.com/go-redis/redis"
	"time"
)

// API namespace
type API struct {
	api.Handler
	redis       *redis.Client
	cacheConfig CacheConfig
}

type RedisConfig struct {
	Host     string `config:"host"`
	Port     string `config:"port"`
	Password string `config:"password"`
	DB       int    `config:"db"`
}

type CacheConfig struct {
	CacheEnabled bool   `config:"enabled"`
	KeyPrefix    string `config:"key_prefix"`
	TTL          string `config:"ttl"`
	duration     *time.Duration
	RedisConfig  RedisConfig `config:"redis"`
}

func (config CacheConfig) GetTTLDuration() *time.Duration {
	if config.duration != nil {
		return config.duration
	}

	if config.TTL != "" {
		dur, err := time.ParseDuration(config.TTL)
		if err != nil {
			dur, _ = time.ParseDuration("10s")
		}
		config.duration = &dur
	}
	return config.duration
}

// InitAPI init apis
func InitAPI() {

	cacheConfig := CacheConfig{
		KeyPrefix:   "proxy_",
		RedisConfig: RedisConfig{Host: "localhost", Port: "6379", Password: "", DB: 0},
	}

	env.ParseConfig("cache", &cacheConfig)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cacheConfig.RedisConfig.Host, cacheConfig.RedisConfig.Port),
		Password: cacheConfig.RedisConfig.Password,
		DB:       cacheConfig.RedisConfig.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Error("cache server is not ready: ", err)
		panic(err)
	}

	apis := API{redis: client, cacheConfig: cacheConfig}

	//Index
	api.HandleAPIMethod(api.HEAD, "/", apis.IndexAction)
	api.HandleAPIMethod(api.GET, "/", apis.IndexAction)
	api.HandleAPIMethod(api.GET, "/favicon.ico", apis.FaviconAction)

	//Stats APIs
	api.HandleAPIMethod(api.GET, "/_proxy/stats", apis.StatsAction)
	api.HandleAPIMethod(api.POST, "/_proxy/queue/resume", apis.QueueResumeAction)
	api.HandleAPIMethod(api.GET, "/_proxy/queue/stats", apis.QueueStatsAction)
	api.HandleAPIMethod(api.GET, "/_proxy/requests/", apis.GetRequestsAction)
	api.HandleAPIMethod(api.POST, "/_proxy/request/redo", apis.RedoRequestsAction)

	// Handle proxy
	api.HandleAPIFunc("/", apis.ProxyAction)
}
