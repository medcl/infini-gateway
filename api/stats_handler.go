/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"github.com/infinitbyte/framework/core/stats"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func getMapValue(mapData map[string]int, key string, defaultValue int32) int {
	data := mapData[key]
	return data
}

// StatsAction return stats information
func (handler API) StatsAction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	m := stats.StatsAll()
	handler.WriteJSONHeader(w)
	handler.Write(w, *m)
}
