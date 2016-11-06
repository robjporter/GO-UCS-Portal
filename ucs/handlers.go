package ucs

import (
    "fmt"
    "net/http"
    "html/template"
    "encoding/json"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    ucsData["PageTitle"]                = "UCS - Home"
    ucsData["ActiveUsername"]           = "Guest"

    testTemplate, err := template.ParseFiles(themePath+"base.tmpl", themePath+"layout/open.tmpl", themePath+"dashboard.tmpl", themePath+"layout/close.tmpl")
    if err != nil { panic(err) }

    w.Header().Set("Content-Type", "text/html")
    err = testTemplate.Execute(w, ucsData)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError)}
}

func DataHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-type", "application/json")

    mutex.RLock()
    defer mutex.RUnlock()
	json_bytes, _ := json.Marshal(ucsData)
	fmt.Fprintf(w, "%s\n", json_bytes)
}
