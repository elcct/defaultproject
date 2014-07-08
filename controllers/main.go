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

// Home page route
func (controller *MainController) Index(c web.C, r *http.Request) (string, int) {	
	t := controller.GetTemplate(c)

	widgets := helpers.Parse(t, "home", nil)

	// With that kind of flags template can "figure out" what route is being rendered
	c.Env["IsIndex"] = true

	c.Env["Title"] = "Default Project - free Go website project template"
	c.Env["Content"] = template.HTML(widgets)

	return helpers.Parse(t, "main", c.Env), http.StatusOK
}

// Sign in route
func (controller *MainController) SignIn(c web.C, r *http.Request) (string, int) {
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

// Sign up route
func (controller *MainController) SignUp(c web.C, r *http.Request) (string, int) {
	t := controller.GetTemplate(c)
	session := controller.GetSession(c)
	
	// With that kind of flags template can "figure out" what route is being rendered
	c.Env["IsSignUp"] = true

	c.Env["Flash"] = session.Flashes("auth")	

	var widgets = controller.Parse(t, "auth/signup", c.Env)

	c.Env["Title"] = "Default Project - Sign Up"
	c.Env["Content"] = template.HTML(widgets)

	return controller.Parse(t, "main", c.Env), http.StatusOK
}


// Sign Up form submit route. Registers new user or shows Sign Up route with appropriate messages set in session
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

// This route logs user out
func (controller *MainController) Logout(c web.C, r *http.Request) (string, int) {	
	session := controller.GetSession(c)

	session.Values["User"] = nil

	return "/", http.StatusSeeOther
}
