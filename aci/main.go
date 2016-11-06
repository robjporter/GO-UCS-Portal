package aci

import (
    "fmt"
    "sync"
    "time"
    "strconv"
    "net/http"
    //"github.com/boltdb/bolt"
    "github.com/gorilla/mux"
    "github.com/jasonlvhit/gocron"
)

var demo bool
var cdn bool
var counter int
var themePath string
var themeName string
var mutex = &sync.RWMutex{}
var aciData map[string]string

func init() {
    demo = false
    cdn = false
    aciData = make(map[string]string)
    aciData["SiteTitle"] = "ACI"
    aciData["SiteURL"] = "/aci/"
    counter = 0
    fillBlankData()
}

func Start() {
    go func() {getAllData()}()
    gocron.Every(60).Seconds().Do(getAllData)
    go func() {<- gocron.Start()}()
}

func AddRoutes(r *mux.Router) {
    r.HandleFunc("/aci",DashboardHandler)
    r.HandleFunc("/aci/data",DataHandler)
    r.PathPrefix("/acistatic/").Handler(http.StripPrefix("/acistatic/", http.FileServer(http.Dir(themePath))))
}

func SetTheme(theme string) {
    themeName = theme
    themePath = "aci/templates/"+theme+"/"
}

func SetDemoStatus(status bool) {
    demo = status
}

func SetCDNStatus(status bool) {
    cdn = status
}

func getAllData() {
    t := time.Now()
    mutex.Lock()
    aciData["LastRefreshDateTime"] = t.Format("Mon Jan _2 2006 15:04:05")
    aciData["NewOrders"] = "0"
    aciData["UnreadMessages"] = "0"
    aciData["Notifications"] = "0"
    mutex.Unlock()
    counter += 1
    logTime("")
    getDataLogin()
    ret := getDataLogout()
    logTime(ret)
}

func logTime(param string) {
    t := time.Now()
    tim := t.Format("Mon Jan _2 2006 15:04:05")
    if param == "" {
        fmt.Println("ACI Inventory collection " + strconv.Itoa(counter) + " started at:   ",tim)
    } else if param == "unknown" {
        fmt.Println("ACI: A potential error has been encountered and some data has not been updated.")
    } else if param == "success" {
        fmt.Println("ACI Inventory collection " + strconv.Itoa(counter) + " completed at: ",tim)
    }
}

func getDataLogin() {
}

func getDataLogout() string {
    return "success"
}

func fillBlankData() {

}

func main() {

}
