package main

import (
    "fmt"
    "log"

    "github.com/msoap/html2data"
)

func main() {
    doc := html2data.FromURL("https://software.cisco.com/download/release.html?mdfid=283612660&catid=282558030&softwareid=283655658&release=3.1(2b)&relind=AVAILABLE&rellifecycle=&reltype=latest")
    // or with config
    // doc := html2data.FromURL("http://example.com", html2data.URLCfg{UA: "userAgent", TimeOut: 10, DontDetectCharset: false})
    if doc.Err != nil {
        log.Fatal(doc.Err)
    }

    // get title
    title, _ := doc.GetDataSingle("title")
    fmt.Println("Title is:", title)

    title, _ = doc.GetDataSingle("title", html2data.Cfg{DontTrimSpaces: true})
    fmt.Println("Title as is, with spaces:", title)

    texts, _ := doc.GetData(map[string]string{"h1": "h1", "links": "div:attr(id)"})
    // get all H1 headers:
    if textOne, ok := texts["h1"]; ok {
        for _, text := range textOne {
            fmt.Println(text)
        }
    }
    // get all urls from links
    if links, ok := texts["links"]; ok {
        for _, text := range links {
            //fmt.Println(text)
            if text == "ptl-sdpriimagedetails-div-id" {
                fmt.Println(text)
            }
        }
    }
}
