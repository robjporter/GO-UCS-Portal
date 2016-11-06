package main

// UPDATE CODE SYNTAX
// gotags -tag-relative=true -R=true -sort=true -f="tags" -fields=+l .

import (
    "fmt"
    "net/http"
	"github.com/gorilla/mux"

	"./demo"
	"./ucs"
    "./aci"
)

var siteTemplate    string
var portWeListenOn  string

func init() {
    portWeListenOn =       "8080"
    ucs.SetDemoStatus(false)
	ucs.SetCDNStatus(false)
    ucs.SetTheme("binary")

    aci.SetDemoStatus(false)
    aci.SetCDNStatus(false)
    aci.SetTheme("flat")

    err := ucs.LoadConfig("config/config.json")
    if err != nil {
        panic(err)
    }
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	demo.AddRoutes(r)
    ucs.AddRoutes(r)
    aci.AddRoutes(r)
	http.Handle("/", r)

    fmt.Println("Server started and listening on port :"+portWeListenOn)
    
    ucs.Start()
    aci.Start()
    err := http.ListenAndServe(":"+portWeListenOn, nil) // set listen port
    if err != nil {
        fmt.Println(err)
    }
}
