package main

import (
  "fmt"
  "log"
  "net/http"
  "strings"
  "runtime"
  "os/exec"

  "github.com/PuerkitoBio/goquery"
)

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
  	case "linux":
  		err = exec.Command("xdg-open", url).Start()
  	case "windows":
  		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
  	case "darwin":
  		err = exec.Command("open", url).Start()
  	default:
  		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// https://pkg.go.dev/github.com/PuerkitoBio/goquery
func ExploreBook() {
  // Request the HTML page.
  res, err := http.Get("https://human.capital/")
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  // Find the review items
  doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Text()
		fmt.Printf("Review %d: %s\n", i, title)

    OUTTER:
    for j := 0; j < len(s.Nodes); j++ {
        for k := 0; k < len(s.Nodes[j].Attr); k++ {
          //fmt.Printf("SAttr %d: %s\n", k, s.Nodes[j].Attr[k])
          href := s.Nodes[j].Attr[k].Val
          if (strings.HasPrefix(href, "http")) {
              fmt.Printf("SVal %d: %s\n", k, href)
              openbrowser(href)
              continue OUTTER
          }
        }
    }
	})
}

func main() {
  ExploreBook()
}
