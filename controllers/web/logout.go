package web

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

// This route logs user out
func (controller *Controller) Logout(c web.C, r *http.Request) (string, int) {
	session := controller.GetSession(c)

	session.Values["User"] = nil

	return "/", http.StatusSeeOther
}
