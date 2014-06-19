package system

import (	
	"net/http"
	"github.com/golang/glog"
	"github.com/zenazn/goji/web"
	"github.com/elcct/defaultproject/models"
	"github.com/gorilla/sessions"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// Makes sure templates are stored in the context
func (application *Application) ApplyTemplates(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c.Env["Template"] = application.Template
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Makes sure controllers can have access to session
func (application *Application) ApplySessions(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session, _ := application.Store.Get(r, "session")
		c.Env["Session"] = session
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Makes sure controllers can have access to the database
func (application *Application) ApplyDatabase(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {				
		session := application.DBSession.Clone()
		defer session.Close()
		c.Env["DBSession"] = session		
		c.Env["DBName"] = application.Configuration.Database.Database
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (application *Application) ApplyAuth(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session := c.Env["Session"].(*sessions.Session)		
		if userId, ok := session.Values["User"].(bson.ObjectId); ok {
			dbSession := c.Env["DBSession"].(*mgo.Session)
			database := dbSession.DB(c.Env["DBName"].(string))

			user := new(models.User)		
			err := database.C("users").Find(bson.M{"_id": userId}).One(&user)
			if err != nil {
				glog.Warningf("Auth error: %v", err)
				c.Env["User"] = nil				
			} else {
				c.Env["User"] = user
			}
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

