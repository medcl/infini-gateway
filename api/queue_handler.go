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
	"github.com/infinitbyte/framework/core/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// QueueStatsAction return queue stats information
func (handler API) QueueStatsAction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	data := map[string]int64{}
	//data["check"] = queue.Depth(config.CheckChannel)
	handler.WriteJSON(w, util.MapStr{
		"depth": data,
	}, 200)
}
