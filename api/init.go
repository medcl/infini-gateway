package api

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/go-redis/redis"
	"infini.sh/framework/core/api"
	"infini.sh/framework/core/env"
	"time"
)

// API namespace
type API struct {
	api.Handler
	cacheHandler *CacheHandler
	cacheConfig  CacheConfig
}

type CacheHandler struct {
	config *CacheConfig
	client *redis.Client
}

func (handler CacheHandler) Init() {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", handler.config.RedisConfig.Host, handler.config.RedisConfig.Port),
		Password: handler.config.RedisConfig.Password,
		DB:       handler.config.RedisConfig.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Error("cache server is not ready: ", err)
		panic(err)
	}
}

func (handler CacheHandler) Get(key string) ([]byte, error) {

	return nil, nil
}

func (handler CacheHandler) Set(string, []byte, int64) (bool, error) {
	return false, nil
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
	duration     int64
	RedisConfig  RedisConfig `config:"redis"`
}

func (config CacheConfig) GetTTLMilliseconds() int64 {
	if config.duration > 0 {
		return config.duration
	}

	if config.TTL != "" {
		dur, err := time.ParseDuration(config.TTL)
		if err != nil {
			dur, _ = time.ParseDuration("10s")
		}
		config.duration = dur.Milliseconds()
	}
	return config.duration
}

// InitAPI init apis
func InitAPI() {

	cacheConfig := CacheConfig{
		KeyPrefix:   "proxy_",
		CacheEnabled: false,
		RedisConfig: RedisConfig{Host: "localhost", Port: "6379", Password: "", DB: 0},
	}



	env.ParseConfig("cache", &cacheConfig)

	cacheHandler := CacheHandler{config: &cacheConfig}

	if cacheConfig.CacheEnabled{
		cacheHandler.Init()
	}

	apis := API{cacheHandler: &cacheHandler, cacheConfig: cacheConfig}

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
