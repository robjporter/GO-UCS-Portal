package ucs

import (
    "fmt"
    "bytes"
    "net/http"
    "io/ioutil"
    "crypto/tls"
)

func getData(url string, body string) string {
    respBody := ""
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    req, err := http.NewRequest("POST", "https://"+url+"/nuova", bytes.NewBuffer([]byte(body)))
    if err != nil {
        return ""
    }
    req.Header.Add("Content-Type", "application/xml; charset=utf-8")
    if demo {
        respBody = getDemoData(body)
    } else {
        resp, err := client.Do(req)
        if err != nil {
            fmt.Println(err)
        }
        defer resp.Body.Close()
        repBody, err := ioutil.ReadAll(resp.Body)
        respBody = string(repBody[:])
    }
	commands[body] = respBody
    //fmt.Println(respBody)
    return respBody
}
