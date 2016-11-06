package demo

import (
    "net/http"
    "html/template"
	"github.com/gorilla/mux"
)

var siteDemoTemplate string

func init() {
	siteDemoTemplate = "demo/templates/default/"
}

func getDemoPathPrefix() string {
    return "/demostatic/"
}

func getDemoPath() string {
    return siteDemoTemplate
}

func GetRoutes() map[string]http.HandlerFunc {
    routes := make(map[string]http.HandlerFunc)
    routes["/demo"] = HomeHandler
    routes["/demo/tabs"] = TabHandler
	routes["/demo/ui"] = UiHandler
	routes["/demo/charts"] = ChartHandler
	routes["/demo/tables"] = TableHandler
	routes["/demo/forms"] = FormsHandler
	routes["/demo/login"] = LoginHandler
	routes["/demo/registration"] = RegistrationHandler
	routes["/demo/blank"] = BlankHandler
    return routes
}

func AddRoutes(r *mux.Router) {
    r.HandleFunc("/demo",HomeHandler)
    r.HandleFunc("/demo/tabs",TabHandler)
    r.HandleFunc("/demo/ui",UiHandler)
    r.HandleFunc("/demo/charts",ChartHandler)
    r.HandleFunc("/demo/tables",TableHandler)
    r.HandleFunc("/demo/forms",FormsHandler)
    r.HandleFunc("/demo/login",LoginHandler)
    r.HandleFunc("/demo/registration",RegistrationHandler)
    r.HandleFunc("/demo/blank",BlankHandler)
    r.PathPrefix(getDemoPathPrefix()).Handler(http.StripPrefix(getDemoPathPrefix(), http.FileServer(http.Dir(getDemoPath()))))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["PageTitle"] = "Template example: map"
	m["Message"] = "Hello"
	m["User"] = "World"
	m["LastGenerated"] = "30 May 2014"
    m["SiteTitle"] = "UCS Portal"

	testTemplate, err := template.ParseFiles(siteDemoTemplate+"base.tmpl", siteDemoTemplate+"layout/open.tmpl", siteDemoTemplate+"dashboard.tmpl", siteDemoTemplate+"layout/close.tmpl")
	if err != nil {
    	panic(err)
  	}

	w.Header().Set("Content-Type", "text/html")
  	err = testTemplate.Execute(w, m)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	}
}

func TabHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["PageTitle"] = "Template example: map"
	m["Message"] = "Hello"
	m["User"] = "World"
	m["LastGenerated"] = "30 May 2014"
    m["SiteTitle"] = "UCS Portal"

	testTemplate, err := template.ParseFiles(siteDemoTemplate+"base.tmpl", siteDemoTemplate+"layout/open.tmpl", siteDemoTemplate+"tabpanel.tmpl", siteDemoTemplate+"layout/close.tmpl")
	if err != nil {
    	panic(err)
  	}

	w.Header().Set("Content-Type", "text/html")
  	err = testTemplate.Execute(w, m)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	}
}

func UiHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["PageTitle"] = "Template example: map"
	m["Message"] = "Hello"
	m["User"] = "World"
	m["LastGenerated"] = "30 May 2014"
    m["SiteTitle"] = "UCS Portal"

	testTemplate, err := template.ParseFiles(siteDemoTemplate+"base.tmpl", siteDemoTemplate+"layout/open.tmpl", siteDemoTemplate+"uielements.tmpl", siteDemoTemplate+"layout/close.tmpl")
	if err != nil {
    	panic(err)
  	}

	w.Header().Set("Content-Type", "text/html")
  	err = testTemplate.Execute(w, m)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	}
}

func ChartHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["PageTitle"] = "Template example: map"
	m["Message"] = "Hello"
	m["User"] = "World"
	m["LastGenerated"] = "30 May 2014"
    m["SiteTitle"] = "UCS Portal"

	testTemplate, err := template.ParseFiles(siteDemoTemplate+"base.tmpl", siteDemoTemplate+"layout/open.tmpl", siteDemoTemplate+"charts.tmpl", siteDemoTemplate+"layout/close.tmpl")
	if err != nil {
    	panic(err)
  	}

	w.Header().Set("Content-Type", "text/html")
  	err = testTemplate.Execute(w, m)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	}
}

func TableHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["PageTitle"] = "Template example: map"
	m["Message"] = "Hello"
	m["User"] = "World"
	m["LastGenerated"] = "30 May 2014"
    m["SiteTitle"] = "UCS Portal"

	testTemplate, err := template.ParseFiles(siteDemoTemplate+"base.tmpl", siteDemoTemplate+"layout/open.tmpl", siteDemoTemplate+"tables.tmpl", siteDemoTemplate+"layout/close.tmpl")
	if err != nil {
    	panic(err)
  	}

	w.Header().Set("Content-Type", "text/html")
  	err = testTemplate.Execute(w, m)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	}
}

func FormsHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["PageTitle"] = "Template example: map"
	m["Message"] = "Hello"
	m["User"] = "World"
	m["LastGenerated"] = "30 May 2014"
    m["SiteTitle"] = "UCS Portal"

	testTemplate, err := template.ParseFiles(siteDemoTemplate+"base.tmpl", siteDemoTemplate+"layout/open.tmpl", siteDemoTemplate+"forms.tmpl", siteDemoTemplate+"layout/close.tmpl")
	if err != nil {
    	panic(err)
  	}

	w.Header().Set("Content-Type", "text/html")
  	err = testTemplate.Execute(w, m)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["PageTitle"] = "Template example: map"
	m["Message"] = "Hello"
	m["User"] = "World"
	m["LastGenerated"] = "30 May 2014"
    m["SiteTitle"] = "UCS Portal"

	testTemplate, err := template.ParseFiles(siteDemoTemplate+"base.tmpl", siteDemoTemplate+"layout/open.tmpl", siteDemoTemplate+"login.tmpl", siteDemoTemplate+"layout/close.tmpl")
	if err != nil {
    	panic(err)
  	}

	w.Header().Set("Content-Type", "text/html")
  	err = testTemplate.Execute(w, m)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	}
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["PageTitle"] = "Template example: map"
	m["Message"] = "Hello"
	m["User"] = "World"
	m["LastGenerated"] = "30 May 2014"
    m["SiteTitle"] = "UCS Portal"

	testTemplate, err := template.ParseFiles(siteDemoTemplate+"base.tmpl", siteDemoTemplate+"layout/open.tmpl", siteDemoTemplate+"registration.tmpl", siteDemoTemplate+"layout/close.tmpl")
	if err != nil {
    	panic(err)
  	}

	w.Header().Set("Content-Type", "text/html")
  	err = testTemplate.Execute(w, m)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	}
}

func BlankHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["PageTitle"] = "Template example: map"
	m["Message"] = "Hello"
	m["User"] = "World"
	m["LastGenerated"] = "30 May 2014"
    m["SiteTitle"] = "UCS Portal"

	testTemplate, err := template.ParseFiles(siteDemoTemplate+"base.tmpl", siteDemoTemplate+"layout/open.tmpl", siteDemoTemplate+"blank.tmpl", siteDemoTemplate+"layout/close.tmpl")
	if err != nil {
    	panic(err)
  	}

	w.Header().Set("Content-Type", "text/html")
  	err = testTemplate.Execute(w, m)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	}
}
