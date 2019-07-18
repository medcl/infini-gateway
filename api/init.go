package api

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/api"
	"github.com/infinitbyte/framework/core/env"
	"src/github.com/go-redis/redis"
)

// API namespace
type API struct {
	api.Handler
	redis       *redis.Client
	cacheConfig CacheConfig
}

type CacheConfig struct {
	CacheEnabled bool `config:"enabled"`
}

// InitAPI init apis
func InitAPI() {

	cacheConfig := CacheConfig{}

	env.ParseConfig("cache", &cacheConfig)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Error(err)
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
