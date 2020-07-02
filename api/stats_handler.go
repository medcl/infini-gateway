package api

import (
	"infini.sh/framework/core/api/router"
	"infini.sh/framework/core/stats"
	"net/http"
)

func getMapValue(mapData map[string]int, key string, defaultValue int32) int {
	data := mapData[key]
	return data
}

// StatsAction return stats information
func (handler *API) StatsAction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	m := stats.StatsAll()
	handler.WriteJSONHeader(w)
	handler.Write(w, *m)
}

func (handler *API) FaviconAction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Write([]byte("."))
}
