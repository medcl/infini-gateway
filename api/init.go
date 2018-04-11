package api

import (
	"github.com/infinitbyte/framework/core/api"
)

// API namespace
type API struct {
	api.Handler
}

// InitAPI init apis
func InitAPI() {

	apis := API{}

	//Index
	api.HandleAPIMethod(api.GET, "/", apis.IndexAction)
	api.HandleAPIMethod(api.GET, "/favicon.ico", apis.FaviconAction)

	//Stats APIs
	api.HandleAPIMethod(api.GET, "/_proxy/stats", apis.StatsAction)
	api.HandleAPIMethod(api.POST, "/_proxy/queue/resume", apis.QueueResumeAction)
	api.HandleAPIMethod(api.GET, "/_proxy/queue/stats", apis.QueueStatsAction)
	//api.HandleAPIMethod(api.GET, "/_proxy/requests/", apis.GetRequestsAction)
	api.HandleAPIMethod(api.POST, "/_proxy/request/redo", apis.RedoRequestsAction)

	// Handle proxy
	api.HandleAPIFunc("/", apis.ProxyAction)
}
