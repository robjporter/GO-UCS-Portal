package ucs

import (
    "fmt"
    "math"
    "sync"
    "time"
    "strings"
    "strconv"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/jasonlvhit/gocron"
    "github.com/dustin/go-humanize"
)

var demo bool
var cdn bool
var themePath string
var themeName string
var ucs UCS
var commands map[string]string
var ucsData map[string]string
var data UCSData
var mutex = &sync.RWMutex{}
var compute []UCSCompute
var counter int

func init() {
    demo = false
    cdn = false
	commands = make(map[string]string)
    ucsData = make(map[string]string)
    ucsData["SiteTitle"] = "UCS Portal"
    counter = 0
    fillBlankData()
}

func Start() {
    go func() {getAllData()}()
    gocron.Every(60).Seconds().Do(getAllData)
    go func() {<- gocron.Start()}()
}

//////////////////////////////////////////////////////////////////////
////  SETUP
//////////////////////////////////////////////////////////////////////

func AddRoutes(r *mux.Router) {
    r.HandleFunc("/ucs",DashboardHandler)
    r.HandleFunc("/ucs/data",DataHandler)
    r.PathPrefix("/ucsstatic/").Handler(http.StripPrefix("/ucsstatic/", http.FileServer(http.Dir(themePath))))
}

func SetTheme(theme string) {
    themeName = theme
    themePath = "ucs/templates/"+theme+"/"
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
    ucsData["LastRefreshDateTime"] = t.Format("Mon Jan _2 2006 15:04:05")
    mutex.Unlock()
    counter += 1
    logTime("")
    getDataLogin()
    getAssignedServiceProfileCount()
    getUnassignedServiceProfileCount()
    getAllMacPoolCount()
    getAllFaultCount()
    getAllBladeCount()
    ret := getDataLogout()
    logTime(ret)
}

func logTime(param string) {
    t := time.Now()
    tim := t.Format("Mon Jan _2 2006 15:04:05")
    if param == "" {
        fmt.Println("UCSM Inventory collection " + strconv.Itoa(counter) + " started at:   ",tim)
    } else if param == "unknown" {
        fmt.Println("UCSM: A potential error has been encountered and some data has not been updated.")
    } else if param == "success" {
        fmt.Println("UCSM Inventory collection " + strconv.Itoa(counter) + " completed at: ",tim)
    }
}

func getDataLogin() {
    pass, err := decryptString(ucsData["UCS_PASSWORD"])
    if err == nil {
        response := login(ucsData["UCS_ADDRESS"], ucsData["UCS_USERNAME"], pass)
        ret := getBulkAttribute(response, "aaaLogin", []string{"outCookie","outVersion","outPriv"})

        mutex.Lock()
        ucsData["UCS_COOKIE"] = getValueOfAttribute(ret,"outCookie")
        ucsData["UCS_VERSION"] = getValueOfAttribute(ret,"outVersion")
        ucsData["UCS_PRIVILEDGE"] = getValueOfAttribute(ret,"outPriv")
        mutex.Unlock()
    }
}

func getDataLogout() string {
    response := logout(ucsData["UCS_ADDRESS"], ucsData["UCS_COOKIE"])
    return getAttribute(response,"aaaLogout","outStatus")
}

func getAssignedServiceProfileCount() {
    servers := getAssociatedServers(ucsData["UCS_ADDRESS"], ucsData["UCS_COOKIE"])

    mutex.Lock()
    ucsData["UCS_ASSIGNED_SP_COUNT"] = strconv.Itoa(getElementCount(servers,[]string{"configResolveClass","outConfigs"},"lsServer"))
    mutex.Unlock()
}

func getUnassignedServiceProfileCount() {
    servers := getUnassociatedServers(ucsData["UCS_ADDRESS"], ucsData["UCS_COOKIE"])

    mutex.Lock()
    ucsData["UCS_UNASSIGNED_SP_COUNT"] = strconv.Itoa(getElementCount(servers,[]string{"configResolveClass","outConfigs"},"lsServer"))
    mutex.Unlock()
}

func getAllMacPoolCount() {
    macPools := getAllMacPools(ucsData["UCS_ADDRESS"], ucsData["UCS_COOKIE"])

    mutex.Lock()
    ucsData["UCS_MAC_POOL_COUNT"] = strconv.Itoa(getElementCount(macPools,[]string{"configScope","outConfigs"},"macpoolAddr"))
    mutex.Unlock()
}

func getTotalCPU(data map[string]string) {
    //total := 0.0
    //activeTotal := 0.0
    //inactiveTotal := 0.0
    //activeCoresTotal := 0
    //inactiveCoresTotal := 0
    xmlOn := ""
    xmlOff := ""

    for x := 0; x < len(data)/2; x++ {
        if data["UCS_COMPUTE_"+strconv.Itoa(x)+"_OPERPOWER"] == "on" {
            for y := 1; y < 9; y++ {
                xmlOn += "<dn value=\""+data["UCS_COMPUTE_"+strconv.Itoa(x)+"_DN"]+"/board/cpu-"+strconv.Itoa(y)+"\" />"
            }
        } else {
            for y := 1; y < 9; y++ {
                xmlOff += "<dn value=\""+data["UCS_COMPUTE_"+strconv.Itoa(x)+"_DN"]+"/board/cpu-"+strconv.Itoa(y)+"\" />"
            }
        }
    }

    cpuOn := getCPUDetail(ucsData["UCS_ADDRESS"], ucsData["UCS_COOKIE"], xmlOn)
    cpuOff := getCPUDetail(ucsData["UCS_ADDRESS"], ucsData["UCS_COOKIE"], xmlOff)
    cpuSpeedOn := getElementOccurence(cpuOn,[]string{"configResolveDns","outConfigs"},"processorUnit","speed")
    cpuCoresOn := getElementOccurence(cpuOn,[]string{"configResolveDns","outConfigs"},"processorUnit","cores")
    cpuSpeedOff := getElementOccurence(cpuOff,[]string{"configResolveDns","outConfigs"},"processorUnit","speed")
    cpuCoresOff := getElementOccurence(cpuOff,[]string{"configResolveDns","outConfigs"},"processorUnit","cores")

    jsonCPUSpeedOn, jtmerr := json.Marshal(cpuSpeedOn)
    jsonCPUCoresOn, jmserr := json.Marshal(cpuCoresOn)
    jsonCPUSpeedOff, jmerr := json.Marshal(cpuSpeedOff)
    jsonCPUCoresOff, jcerr := json.Marshal(cpuCoresOff)

    coresOn := getTotal(cpuCoresOn)
    coresOff := getTotal(cpuCoresOff)
    speedOn := getTotal(cpuSpeedOn)
    speedOff := getTotal(cpuSpeedOff)

    corOn,_ := strconv.Atoi(coresOn)
    corOf,_ := strconv.Atoi(coresOff)
    speOn,_ := strconv.Atoi(speedOn)
    speOf,_ := strconv.Atoi(speedOff)

    mutex.Lock()
    if jtmerr == nil {ucsData["UCS_BLADE_CPU_SPEEDS_ON_JSON"] = string(jsonCPUSpeedOn)}
    if jmserr == nil {ucsData["UCS_BLADE_CPU_CORES_ON_JSON"] = string(jsonCPUCoresOn)}
    if jmerr == nil {ucsData["UCS_BLADE_CPU_SPEEDS_OFF_JSON"] = string(jsonCPUSpeedOff)}
    if jcerr == nil {ucsData["UCS_BLADE_CPU_CORES_OFF_JSON"] = string(jsonCPUCoresOff)}
    ucsData["UCS_BLADE_CPU_CORES_ON"] = coresOn
    ucsData["UCS_BLADE_CPU_CORES_OFF"] = coresOff
    ucsData["UCS_BLADE_CPU_SPEED_ON"] = speedOn + "Ghz"
    ucsData["UCS_BLADE_CPU_SPEED_OFF"] = speedOff + "Ghz"
    ucsData["UCS_BLADE_CPU_CORES_TOTAL"] = strconv.Itoa(corOn + corOf)
    ucsData["UCS_BLADE_CPU_SPEED_TOTAL"] = strconv.Itoa(speOn + speOf) + "Ghz"
    mutex.Unlock()

    //fmt.Println(total)
    //fmt.Println(activeTotal)
    //fmt.Println(inactiveTotal)
    //fmt.Println(activeCoresTotal)
    //fmt.Println(inactiveCoresTotal)

}

func getAllBladeCount() {
    blades := getAllBlades(ucsData["UCS_ADDRESS"], ucsData["UCS_COOKIE"])

    mutex.Lock()
    ucsData["UCS_BLADE_COUNT"] = strconv.Itoa(getElementCount(blades,[]string{"configResolveClass","outConfigs"},"computeBlade"))
    mutex.Unlock()

    stateCount := getElementOccurence(blades,[]string{"configResolveClass","outConfigs"},"computeBlade","availability")
    totalMemory := getElementOccurence(blades,[]string{"configResolveClass","outConfigs"},"computeBlade","availableMemory")
    memorySpeed := getElementOccurence(blades,[]string{"configResolveClass","outConfigs"},"computeBlade","memorySpeed")
    models := getElementOccurence(blades,[]string{"configResolveClass","outConfigs"},"computeBlade","model")
    adapters := getElementOccurence(blades,[]string{"configResolveClass","outConfigs"},"computeBlade","numOfAdaptors")
    cores := getElementOccurence(blades,[]string{"configResolveClass","outConfigs"},"computeBlade","numOfCores")
    cpus := getElementOccurence(blades,[]string{"configResolveClass","outConfigs"},"computeBlade","numOfCpus")
    ethifs := getElementOccurence(blades,[]string{"configResolveClass","outConfigs"},"computeBlade","numOfEthHostIfs")
    fcifs := getElementOccurence(blades,[]string{"configResolveClass","outConfigs"},"computeBlade","numOfFcHostIfs")
    power := getElementOccurence(blades,[]string{"configResolveClass","outConfigs"},"computeBlade","operPower")
    cpuDNs := getElementArray(blades,[]string{"configResolveClass","outConfigs"},"computeBlade",[]string{"dn","operPower"})
    getTotalCPU(cpuDNs)

    jsonTotalMemory, jtmerr := json.Marshal(totalMemory)
    jsonMemorySpeed, jmserr := json.Marshal(memorySpeed)
    jsonModels, jmerr := json.Marshal(models)
    jsonCores, jcerr := json.Marshal(cores)
    jsonCPUs, jcperr := json.Marshal(cpus)
    jsonAdapters, jaerr := json.Marshal(adapters)
    jsonEthIfs, jeperr := json.Marshal(ethifs)
    jsonFcIfs, jfperr := json.Marshal(fcifs)

    sumMem, _ := strconv.Atoi(getTotal(totalMemory))
    sumMem = sumMem * 1000000
    sumMemory := humanize.Bytes(uint64(sumMem))
    sumMem = getElementOccuranceVariable(blades,[]string{"configResolveClass","outConfigs"},"computeBlade",[]string{"operPower","on"},"availableMemory")
    sumMem = sumMem * 1000000
    totalActiveMemory := humanize.Bytes(uint64(sumMem))
    sumMem = getElementOccuranceVariable(blades,[]string{"configResolveClass","outConfigs"},"computeBlade",[]string{"operPower","off"},"availableMemory")
    sumMem = sumMem * 1000000
    totalInactiveMemory := humanize.Bytes(uint64(sumMem))

    averageMemory := getAverageTotal(totalMemory,ucsData["UCS_BLADE_COUNT"])
    averageMemorySpeed := getAverageTotal(memorySpeed,ucsData["UCS_BLADE_COUNT"])
    averageModel := getHighest(models)
    averageAdapters := getHighest(adapters)
    averageCores := getAverageTotal(cores,ucsData["UCS_BLADE_COUNT"])
    averageCPUS := getAverageTotal(cpus,ucsData["UCS_BLADE_COUNT"])
    averageEthernet := getHighest(ethifs)
    averageFC := getHighest(fcifs)

    minimumMemory := getMinimumTotal(totalMemory)
    minimumMemorySpeed := getMinimumTotal(memorySpeed)
    minimumAdapters := getMinimumTotal(adapters)
    minimumCores := getMinimumTotal(cores)
    minimumCPUS := getMinimumTotal(cpus)
    minimumEthernet := getMinimumTotal(ethifs)
    minimumFC := getMinimumTotal(fcifs)

    maximumMemory := getMaximumTotal(totalMemory)
    maximumMemorySpeed := getMaximumTotal(memorySpeed)
    maximumAdapters := getMaximumTotal(adapters)
    maximumCores := getMaximumTotal(cores)
    maximumCPUS := getMaximumTotal(cpus)
    maximumEthernet := getMaximumTotal(ethifs)
    maximumFC := getMaximumTotal(fcifs)

    mutex.Lock()
    ucsData["UCS_BLADE_MEMORY_TOTAL"] = sumMemory
    ucsData["UCS_BLADE_AVAILABLE"] = strconv.Itoa(stateCount["available"])
    ucsData["UCS_BLADE_UNAVAILABLE"] = strconv.Itoa(stateCount["unavailable"])
    ucsData["UCS_BLADE_POWER_ON"] = strconv.Itoa(power["on"])
    ucsData["UCS_BLADE_POWER_OFF"] = strconv.Itoa(power["off"])
    ucsData["UCS_BLADE_ACTIVE_MEMORY"] = totalActiveMemory
    ucsData["UCS_BLADE_INACTIVE_MEMORY"] = totalInactiveMemory
    if jtmerr == nil {ucsData["UCS_BLADE_TOTAL_MEMORY"] = string(jsonTotalMemory)}
    if jmserr == nil {ucsData["UCS_BLADE_MEMORY_SPEED"] = string(jsonMemorySpeed)}
    if jmerr == nil {ucsData["UCS_BLADE_MODELS"] = string(jsonModels)}
    if jcerr == nil {ucsData["UCS_BLADE_CORES"] = string(jsonCores)}
    if jcperr == nil {ucsData["UCS_BLADE_CPUS"] = string(jsonCPUs)}
    if jaerr == nil {ucsData["UCS_BLADE_ADAPTERS"] = string(jsonAdapters)}
    if jeperr == nil {ucsData["UCS_BLADE_ETHERNET"] = string(jsonEthIfs)}
    if jfperr == nil {ucsData["UCS_BLADE_FC"] = string(jsonFcIfs)}
    ucsData["UCS_BLADE_AVERAGE_MEMORY"] = averageMemory
    ucsData["UCS_BLADE_AVERAGE_MEMORY_SPEED"] = averageMemorySpeed
    ucsData["UCS_BLADE_AVERAGE_MODEL"] = averageModel
    ucsData["UCS_BLADE_AVERAGE_ADAPTERS"] = averageAdapters
    ucsData["UCS_BLADE_AVERAGE_CORES"] = averageCores
    ucsData["UCS_BLADE_AVERAGE_CPUS"] = averageCPUS
    ucsData["UCS_BLADE_AVERAGE_ETHERNET"] = averageEthernet
    ucsData["UCS_BLADE_AVERAGE_FIBRE"] = averageFC
    ucsData["UCS_BLADE_MINIMUM_MEMORY"] = minimumMemory
    ucsData["UCS_BLADE_MINIMUM_MEMORY_SPEED"] = minimumMemorySpeed
    ucsData["UCS_BLADE_MINIMUM_ADAPTERS"] = minimumAdapters
    ucsData["UCS_BLADE_MINIMUM_CORES"] = minimumCores
    ucsData["UCS_BLADE_MINIMUM_CPUS"] = minimumCPUS
    ucsData["UCS_BLADE_MINIMUM_ETHERNET"] = minimumEthernet
    ucsData["UCS_BLADE_MINIMUM_FC"] = minimumFC
    ucsData["UCS_BLADE_MAXIMUM_MEMORY"] = maximumMemory
    ucsData["UCS_BLADE_MAXIMUM_MEMORY_SPEED"] = maximumMemorySpeed
    ucsData["UCS_BLADE_MAXIMUM_ADAPTERS"] = maximumAdapters
    ucsData["UCS_BLADE_MAXIMUM_CORES"] = maximumCores
    ucsData["UCS_BLADE_MAXIMUM_CPUS"] = maximumCPUS
    ucsData["UCS_BLADE_MAXIMUM_ETHERNET"] = maximumEthernet
    ucsData["UCS_BLADE_MAXIMUM_FC"] = maximumFC
    mutex.Unlock()

}



func getTotal(data map[string]int) string {
    iTotal := 0
    fTotal := 0.0
    integer := false

    for k,v := range data {
        num,err := strconv.Atoi(k)
        num2,err2 := strconv.ParseFloat(k, 64)

        if err == nil && err2 == nil {
            iTotal += num*v
            integer = true
        } else if err !=nil && err2 == nil {
            fTotal += num2*float64(v)
        }
    }

    if integer {
        return strconv.Itoa(iTotal)
    } else {
        return strconv.Itoa(round(fTotal))
    }
}

func getMinimumTotal(data map[string]int) string {
    lowest := 512000
    for v,_ := range data {
        num,err := strconv.Atoi(v)

        if err == nil {
            if num < lowest {
                lowest = num
            }
        }
    }
    return strconv.Itoa(lowest)
}

func getMaximumTotal(data map[string]int) string {
    highest := 0
    for v,_ := range data {
        num,err := strconv.Atoi(v)

        if err == nil {
            if num > highest {
                highest = num
            }
        }
    }
    return strconv.Itoa(highest)
}

func getAverageTotal(data map[string]int, count string) string {
    total := 0

    for k,v := range data {
        tmp,err := strconv.Atoi(k)
        if err == nil {
            total += tmp * v
        }
    }

    num,err := strconv.Atoi(count)

    if err == nil {
        if total > num {
            return strconv.Itoa(total / num)
        }
        return strconv.Itoa(total)
    }
    return ""
}

func getHighest(data map[string]int) string {
    highest := 0
    result := ""

    for k,v := range data {
        if v > highest {
            highest = v
            result = k
        }
    }

    return result
}

func round(f float64) int {
	return int(f + math.Copysign(0.5, f))
}

func getAllFaultCount() {
    faults := getAllFaults(ucsData["UCS_ADDRESS"], ucsData["UCS_COOKIE"])

    mutex.Lock()
    ucsData["UCS_FAULT_COUNT"] = strconv.Itoa(getElementCount(faults,[]string{"configResolveClass","outConfigs"},"faultInst"))
    mutex.Unlock()

    faultSeverities := getElementOccurence(faults,[]string{"configResolveClass","outConfigs"},"faultInst","severity")

    mutex.Lock()
    for k,v := range faultSeverities {
        ucsData["UCS_FAULT_"+strings.ToUpper(k)] = strconv.Itoa(v)
    }
    mutex.Unlock()
}

func fillBlankData() {
    ucsData["LastRefreshDateTime"] = ""
    ucsData["UCS_ADDRESS"] = ""
    ucsData["UCS_USERNAME"] = ""
    ucsData["UCS_PASSWORD"] = ""
    ucsData["UCS_COOKIE"] = ""
    ucsData["UCS_VERSION"] = ""
    ucsData["UCS_PRIVILEDGE"] = ""
    ucsData["UCS_ASSIGNED_SP_COUNT"] = ""
    ucsData["UCS_UNASSIGNED_SP_COUNT"] = ""
    ucsData["UCS_UNASSIGNED_SP_COUNT"] = ""
    ucsData["UCS_MAC_POOL_COUNT"] = ""
    ucsData["UCS_BLADE_AVAILABLE"] = ""
    ucsData["UCS_BLADE_UNAVAILABLE"] = ""
    ucsData["UCS_BLADE_POWER_ON"] = ""
    ucsData["UCS_BLADE_POWER_OFF"] = ""
    ucsData["UCS_BLADE_TOTAL_MEMORY"] = ""
    ucsData["UCS_BLADE_MEMORY_SPEED"] = ""
    ucsData["UCS_BLADE_MODELS"] = ""
    ucsData["UCS_BLADE_CORES"] = ""
    ucsData["UCS_BLADE_CPUS"] = ""
    ucsData["UCS_BLADE_ADAPTERS"] = ""
    ucsData["UCS_BLADE_ETHERNET"] = ""
    ucsData["UCS_BLADE_FC"] = ""
    ucsData["UCS_BLADE_AVERAGE_MEMORY"] = ""
    ucsData["UCS_BLADE_AVERAGE_MEMORY_SPEED"] = ""
    ucsData["UCS_BLADE_AVERAGE_MODEL"] = ""
    ucsData["UCS_BLADE_AVERAGE_ADAPTERS"] = ""
    ucsData["UCS_BLADE_AVERAGE_CORES"] = ""
    ucsData["UCS_BLADE_AVERAGE_CPUS"] = ""
    ucsData["UCS_BLADE_AVERAGE_ETHERNET"] = ""
    ucsData["UCS_BLADE_AVERAGE_FIBRE"] = ""
    ucsData["UCS_BLADE_MINIMUM_MEMORY"] = ""
    ucsData["UCS_BLADE_MINIMUM_MEMORY_SPEED"] = ""
    ucsData["UCS_BLADE_MINIMUM_ADAPTERS"] = ""
    ucsData["UCS_BLADE_MINIMUM_CORES"] = ""
    ucsData["UCS_BLADE_MINIMUM_CPUS"] = ""
    ucsData["UCS_BLADE_MINIMUM_ETHERNET"] = ""
    ucsData["UCS_BLADE_MINIMUM_FC"] = ""
    ucsData["UCS_BLADE_MAXIMUM_MEMORY"] = ""
    ucsData["UCS_BLADE_MAXIMUM_MEMORY_SPEED"] = ""
    ucsData["UCS_BLADE_MAXIMUM_ADAPTERS"] = ""
    ucsData["UCS_BLADE_MAXIMUM_CORES"] = ""
    ucsData["UCS_BLADE_MAXIMUM_CPUS"] = ""
    ucsData["UCS_BLADE_MAXIMUM_ETHERNET"] = ""
    ucsData["UCS_BLADE_MAXIMUM_FC"] = ""
    ucsData["UCS_BLADE_CPU_SPEEDS_ON_JSON"] = ""
    ucsData["UCS_BLADE_CPU_CORES_ON_JSON"] = ""
    ucsData["UCS_BLADE_CPU_SPEEDS_OFF_JSON"] = ""
    ucsData["UCS_BLADE_CPU_CORES_OFF_JSON"] = ""
    ucsData["UCS_BLADE_CPU_CORES_ON"] = ""
    ucsData["UCS_BLADE_CPU_CORES_OFF"] = ""
    ucsData["UCS_BLADE_CPU_SPEED_ON"] = ""
    ucsData["UCS_BLADE_CPU_SPEED_OFF"] = ""
    ucsData["UCS_BLADE_CPU_CORES_TOTAL"] = ""
    ucsData["UCS_BLADE_CPU_SPEED_TOTAL"] = ""
}

func SetUCSAddress(address string) {
    mutex.Lock()
    ucsData["UCS_ADDRESS"] = address
    mutex.Unlock()
}

func SetUCSUsername(username string) {
    mutex.Lock()
    ucsData["UCS_USERNAME"] = username
    mutex.Unlock()
}

func SetUCSPassword(password string) {
    if password != "" {
        pass, err := encryptString(password)
        if err == nil {
            mutex.Lock()
            ucsData["UCS_PASSWORD"] = pass
            mutex.Unlock()
        }
    }
}

func ConfigDisplay(w http.ResponseWriter, r *http.Request) {
    /*
    if r.Method == "GET" {
        ip := "IP Address"
        username := "Username"
        password := "Password"
        if _, err := os.Stat("config.json"); err == nil {
            buf := bytes.NewBuffer(nil)
            f, _ := os.Open("config.json") // Error handling elided for brevity.
            io.Copy(buf, f)           // Error handling elided for brevity.
            f.Close()
            var dat map[string]interface{}
            if err := json.Unmarshal(buf.Bytes(), &dat); err != nil {
                fmt.Println(err)
            }
            ip = dat["ip"].(string)
            username = dat["username"].(string)
            password = dat["password"].(string)
        }
        display := setupHeader()
        display += "<form method=\"POST\">"
        display += "<div class=\"row\">"
        display += "<div class=\"col-xs-5\">"
        display += "<h4>IP Address of UCS Manager:</h4>"
        display += "</div>"
        display += "<div class=\"col-xs-7\">"
        display += "<input class=\"form-control form-control-lg\" name=\"ip\" type=\"text\" value=\""+ip+"\" >"
        display += "</div>"
        display += "</div>"
        display += "</br>"
        display += "<div class=\"row\">"
        display += "<div class=\"col-xs-5\">"
        display += "<h4>Username for UCS Manager:</h4>"
        display += "</div>"
        display += "<div class=\"col-xs-7\">"
        display += "<input class=\"form-control form-control-lg\" name=\"username\" type=\"text\" value=\""+username+"\">"
        display += "</div>"
        display += "</div>"
        display += "</br>"
        display += "<div class=\"row\">"
        display += "<div class=\"col-xs-5\">"
        display += "<h4>Password UCS Manager:</h4>"
        display += "</div>"
        display += "<div class=\"col-xs-7\">"
        display += "<input class=\"form-control form-control-lg\" name=\"password\" type=\"password\" value=\""+password+"\">"
        display += "</div>"
        display += "</div>"
        display += "</br></br>"
        display += "<div class=\"row\">"
        display += "<div class=\"col-xs-12\">"
        display += "<button type=\"submit\" class=\"btn btn-success btn-block\">Submit</button>"
        display += "</div>"
        display += "</div>"
        display += "</form>"
        display += setupFooter()
        fmt.Fprintf(w, display)
    } else {
        r.ParseForm()  //Parse url parameters passed, then parse the response packet for the POST body (request body)
        ip := r.Form["ip"][0]
        username := r.Form["username"][0]
        password := r.Form["password"][0]
        if ip != "" && username != "" && password != "" {
            config := Config{
                IP: ip,
                Username: username,
                Password: password,
            }
            c, _ := json.Marshal(config)
            f, err := os.Create("config.json")
            if err != nil {
                fmt.Println("There has been an error saving your config file.")
            } else {
                defer f.Close()
                _, err := f.Write(c)
                if err != nil {
                    fmt.Println(err)
                } else {
                    fmt.Println("Config saved successfully.")
                    fmt.Fprintf(w,"The config file has been successfully updated.")
                }
            }
        } else {
            fmt.Fprintf(w, "There were some values missing and so no update has occured to the config file.")
        }
    }
    */
}
