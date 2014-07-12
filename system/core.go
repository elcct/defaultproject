package system

import (	
	"encoding/json"
	"encoding/gob"
	"github.com/golang/glog"
	"net/http"
	
	"reflect"

	"github.com/zenazn/goji/web"

	"github.com/gorilla/sessions"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"


	"html/template"

	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Application struct {
	Configuration *Configuration
	Template      *template.Template
	Store         *sessions.CookieStore
	DBSession     *mgo.Session
}

func (application *Application) Init(filename *string) {
	gob.Register(bson.ObjectId(""))

	data, err := ioutil.ReadFile(*filename)

	if err != nil {
		glog.Fatalf("Can't read configuration file: %s", err)
		panic(err)
	}

	application.Configuration = &Configuration{}

	err = json.Unmarshal(data, &application.Configuration)

	if err != nil {
		glog.Fatalf("Can't parse configuration file: %s", err)
		panic(err)
	}

	application.Store = sessions.NewCookieStore([]byte(application.Configuration.Secret))
}

func (application *Application) LoadTemplates() error {
	var templates []string

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			templates = append(templates, path)
		}
		return nil
	}

	err := filepath.Walk(application.Configuration.TemplatePath, fn)

	if err != nil {
		return err
	}

	application.Template = template.Must(template.ParseFiles(templates...))
	return nil
}


func (application *Application) ConnectToDatabase() {
	var err error
	application.DBSession, err = mgo.Dial(application.Configuration.Database.Hosts)

	if err != nil {
		glog.Fatalf("Can't connect to the database: %v", err)
		panic(err)
	}
}

func (application *Application) Close() {
	glog.Info("Bye!")
	application.DBSession.Close()
}

func (application *Application) Route(controller interface{}, route string) interface{} {	
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {		
		c.Env["Content-Type"] = "text/html"

		methodValue := reflect.ValueOf(controller).MethodByName(route)
		methodInterface := methodValue.Interface()
		method := methodInterface.(func(c web.C, r *http.Request) (string, int))

		body, code := method(c, r)

		if session, exists := c.Env["Session"]; exists {		
			err := session.(*sessions.Session).Save(r, w)
			if err != nil {
				glog.Errorf("Can't save session: %v", err)				
			}
		}

		switch code {
			case http.StatusOK:
				if _, exists := c.Env["Content-Type"]; exists {
					w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
				}
				io.WriteString(w, body)
			case http.StatusSeeOther, http.StatusFound:
				http.Redirect(w, r, body, code)
			default:
				w.WriteHeader(code)
				io.WriteString(w, body)				
		}
	}
	return fn
}
