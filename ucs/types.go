package ucs

type Config struct {
    IP        string        `json:"ip"`
    Username    string      `json:"username"`
    Password    string      `json:"password"`
}

type UCS struct {
    responses []string
    cookie string
}

type UCSData struct {
    lastRefreshDateTime         string
    bladeCount                  int
    bladeCountActive            int
    serviceProfileCount         int
    serviceProfileCountActive   int
}

type UCSCompute struct {
    state                       bool
    assigned                    string
    available                   bool
    memory                      string
    chassis                     int
    checkpoint                  string
    dn                          string
    mgmt                        string
    memoryspeed                 string
    model                       string
    name                        string
    adapters                    int
    cores                       int
    cpus                        int
    ethif                       int
    fcif                        int
    power                       bool
    serial                      string
    slot                        int
}
