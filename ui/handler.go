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

package ui

import (
	"github.com/infinitbyte/framework/core/api"
	"github.com/infinitbyte/framework/core/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type UI struct {
	api.Handler
}

func (h UI) RedirectHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url := h.Get(r, "url", "")
	http.Redirect(w, r, util.UrlDecode(url), 302)
	return
}
