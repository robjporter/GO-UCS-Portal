package aci

import (
    "fmt"
    "strconv"
    "net/http"
    "math/rand"
    "html/template"
    "encoding/json"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    aciData["PageTitle"]                = "ACI - Home"
    aciData["ActiveUsername"]           = "Guest"

    fmt.Println(themePath)

    testTemplate, err := template.ParseFiles(themePath+"base2.tmpl", themePath+"layout/open.tmpl", themePath+"dashboard.tmpl", themePath+"layout/close.tmpl")
    if err != nil { panic(err) }

    w.Header().Set("Content-Type", "text/html")
    err = testTemplate.Execute(w, aciData)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError)}
}

func DataHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-type", "application/json")
    aciData["PageTitle"]                = "ACI - Dashboard"

    mutex.Lock()
    aciData["New"] = strconv.Itoa(rand.Intn(1000))
    aciData["Tasks"] = strconv.Itoa(rand.Intn(1000))
    aciData["New2"] = strconv.Itoa(rand.Intn(1000))
    aciData["Orders"] = strconv.Itoa(rand.Intn(1000))
    aciData["Issues"] = strconv.Itoa(rand.Intn(1000))
    mutex.Unlock()

    mutex.RLock()
    defer mutex.RUnlock()
	json_bytes, _ := json.Marshal(aciData)
	fmt.Fprintf(w, "%s\n", json_bytes)
}
