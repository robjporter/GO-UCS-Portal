package main

import (
  "fmt"
  "log"

  "github.com/PuerkitoBio/goquery"
)

func ExampleScrape() {
  doc, err := goquery.NewDocument("https://software.cisco.com/download/release.html?mdfid=283612660&catid=282558030&softwareid=283655658&release=3.1(2b)&relind=AVAILABLE&rellifecycle=&reltype=latest")
  if err != nil {
    log.Fatal(err)
  }

  // Find the review items
  doc.Find("#ptl-sdpriimagedetails-div-id div table tbody tr td div div #treeAndLegend").Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
    band := s.Find("div").Text()
    title := s.Find("i").Text()
    fmt.Printf("Review %d: %s - %s\n", i, band, title)
  })
}

func main() {
  ExampleScrape()
}
