package main

import (
    "fmt"
    "log"
    "github.com/PuerkitoBio/goquery"
)

func main() {
    //ExampleScrape("http://www.cisco.com/c/en/us/support/servers-unified-computing/ucs-manager/products-release-notes-list.html")
    ExampleScrape("https://software.cisco.com/download/release.html?mdfid=283612660&catid=282558030&softwareid=283655658&release=3.1(2b)&relind=AVAILABLE&rellifecycle=&reltype=latest")
}

func ExampleScrape(url string) {
  doc, err := goquery.NewDocument(url)
  if err != nil {
    log.Fatal(err)
  }

  // Find the review items
  //doc.Find("#fw-content .visitedlinks .wide-v2 .lll-cq .low-level ul:nth-child(2) ul:nth-child(2) ul li:nth-child(1)").Each(func(i int, s *goquery.Selection) {
  // band := s.Find("a").Text()
  //fmt.Println(band)
  //})
  doc.Find(".ptl-full .outerDiv .csWrapper .csDashboard #ptl-sdpriimagedetails-div-id div table tbody tr div div #treeAndLegend .treeOuter").Each(func(i int, s *goquery.Selection) {
      band := s.Find("div").Text()
      fmt.Println(band)
  })

}

/*
/html/body/div[1]/div/table/tbody/tr/td/div/div[3]/div/div/div/div[1]/div/table/tbody/tr/td[1]/div/div/div/div/div/div[2]/div[2]/div[1]/a

#fw-mb #fw-mb-w1 #framework-base-main tbody tr #framework-column-center .ptl-full div #sdpridiv.outerDiv .csWrapper .csDashboard #ptl-sdpriimagedetails-div-id div table tbody tr td div div #treeAndLegend .treeOuter .tree #dmdfTree0.clip #dmdfTree1.clip #div2.treeNode

html body#wps.cdc-fw.modal-open div#fw-mb div#fw-mb-w1 table#framework-base-main tbody tr td#framework-column-center div.ptl-full div div#sdpridiv.outerDiv div.csWrapper div.csDashboard div#ptl-sdpriimagedetails-div-id div table tbody tr td div div div#treeAndLegend div.treeOuter div.tree div#dmdfTree0.clip div#dmdfTree1.clip div#div2.treeNode a#smdfTree2.nodeSel

*/
