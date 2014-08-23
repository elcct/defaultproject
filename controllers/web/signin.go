package web

import (
	"net/http"

	"github.com/elcct/defaultproject/helpers"
	"github.com/zenazn/goji/web"
	"html/template"
)

// Sign in route
func (controller *Controller) SignIn(c web.C, r *http.Request) (string, int) {
	t := controller.GetTemplate(c)
	session := controller.GetSession(c)

	// With that kind of flags template can "figure out" what route is being rendered
	c.Env["IsSignIn"] = true

	c.Env["Flash"] = session.Flashes("auth")
	var widgets = controller.Parse(t, "auth/signin", c.Env)

	c.Env["Title"] = "Default Project - Sign In"
	c.Env["Content"] = template.HTML(widgets)

	return controller.Parse(t, "main", c.Env), http.StatusOK
}

// Sign In form submit route. Logs user in or set appropriate message in session if login was not succesful
func (controller *Controller) SignInPost(c web.C, r *http.Request) (string, int) {
	email, password := r.FormValue("email"), r.FormValue("password")

	session := controller.GetSession(c)
	database := controller.GetDatabase(c)

	user, err := helpers.Login(database, email, password)

	if err != nil {
		session.AddFlash("Invalid Email or Password", "auth")
		return controller.SignIn(c, r)
	}

	session.Values["User"] = user.ID

	return "/", http.StatusSeeOther
}
