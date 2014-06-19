package controllers

import (			
	"github.com/golang/glog"
	"net/http"

	"github.com/zenazn/goji/web"
	"html/template"
	"github.com/elcct/defaultproject/helpers"
	"github.com/elcct/defaultproject/system"
	"github.com/elcct/defaultproject/models"
	"time"
)

type MainController struct {
	system.Controller
}

func (controller *MainController) Index(c web.C, r *http.Request) (string, int) {	
	t := controller.GetTemplate(c)

	widgets := helpers.Parse(t, "home", nil)

	c.Env["Title"] = "Default Project"
	c.Env["Content"] = template.HTML(widgets)

	return helpers.Parse(t, "main", c.Env), http.StatusOK
}

func (controller *MainController) SignIn(c web.C, r *http.Request) (string, int) {
	t := controller.GetTemplate(c)
	session := controller.GetSession(c)
	
	w := struct { Flash []interface{} } { session.Flashes("auth") }	
	var widgets = controller.Parse(t, "auth/signin", w)

	c.Env["Title"] = "Default Project - Sign In"
	c.Env["Content"] = template.HTML(widgets)

	return controller.Parse(t, "main", c.Env), http.StatusOK
}

func (controller *MainController) SignInPost(c web.C, r *http.Request) (string, int) {
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

func (controller *MainController) SignUp(c web.C, r *http.Request) (string, int) {
	t := controller.GetTemplate(c)
	session := controller.GetSession(c)
	
	w := struct { Flash []interface{} } { session.Flashes("auth") }	
	var widgets = controller.Parse(t, "auth/signup", w)

	c.Env["Title"] = "Default Project - Sign Up"
	c.Env["Content"] = template.HTML(widgets)

	return controller.Parse(t, "main", c.Env), http.StatusOK
}

func (controller *MainController) SignUpPost(c web.C, r *http.Request) (string, int) {
	email, password := r.FormValue("email"), r.FormValue("password")

	session := controller.GetSession(c)
	database := controller.GetDatabase(c)

	user := models.GetUserByEmail(database, email)

	if user != nil {
		session.AddFlash("User exists", "auth")
		return controller.SignUp(c, r)
	}

	user = &models.User{
		Username: email,
		Email: email,
		Timestamp: time.Now(),	
	}
	user.HashPassword(password)

	if err := models.InsertUser(database, user); err != nil {
		session.AddFlash("Error whilst registering user.")
		glog.Errorf("Error whilst registering user: %v", err)
		return controller.SignUp(c, r)
	}

	session.Values["User"] = user.ID	

	return "/", http.StatusSeeOther
}

func (controller *MainController) Logout(c web.C, r *http.Request) (string, int) {	
	session := controller.GetSession(c)

	session.Values["User"] = nil

	return "/", http.StatusSeeOther
}
