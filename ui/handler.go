package ui

import (
	"github.com/infinitbyte/framework/core/api"
	"github.com/infinitbyte/framework/core/api/router"
	"github.com/infinitbyte/framework/core/util"
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
