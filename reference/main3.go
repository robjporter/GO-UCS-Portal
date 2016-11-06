package main

import (
    "fmt"
    "net/http"
    "github.com/yhat/scrape"
    "golang.org/x/net/html"
    "golang.org/x/net/html/atom"
)

func main() {
    ExampleScrape("https://software.cisco.com/download/release.html?mdfid=283612660&catid=282558030&softwareid=283655658&release=3.1(2b)&relind=AVAILABLE&rellifecycle=&reltype=latest")
}

func ExampleScrape(url string) {
    resp, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    root, err := html.Parse(resp.Body)
    if err != nil {
        panic(err)
    }

    // define a matcher
    matcher := func(n *html.Node) bool {
        // must check for nil values
        if n.DataAtom == atom.A && n.Parent != nil && n.Parent.Parent != nil {
            return scrape.Attr(n.Parent.Parent, "class") == "treenode"
        }
        return false
    }
    // grab all articles and print them
    articles := scrape.FindAll(root, matcher)
    for i, article := range articles {
        fmt.Printf("%2d %s (%s)\n", i, scrape.Text(article), scrape.Attr(article, "href"))
    }
}
