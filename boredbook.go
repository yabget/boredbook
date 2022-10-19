package main

import (
  "fmt"
  "log"
  "net/http"
  "strings"
  "runtime"
  "os/exec"
  "bufio"
  "os"
  "time"

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
func ExploreBook(url string) {
  // Request the HTML page.
  res, err := http.Get(url)
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

  openedSite := make(map[string]bool)

  // Find the review items
  doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title

    OUTTER:
    for j := 0; j < len(s.Nodes); j++ {
        for k := 0; k < len(s.Nodes[j].Attr); k++ {
          //fmt.Printf("SAttr %d: %s\n", k, s.Nodes[j].Attr[k])
          href := s.Nodes[j].Attr[k].Val
          if (strings.HasPrefix(href, "http") && !openedSite[href]) {
              openbrowser(href)
              openedSite[href] = true
              continue OUTTER
          }
        }
    }
	})
}

func main() {

  // https://gosamples.dev/read-user-input/
  fmt.Printf("Extracting urls from bookmarks.html\n")

  err := exec.Command("./extractURLs.sh", "bookmarks.html").Start()
  if err != nil {
    log.Fatal(err)
  }

  time.Sleep(2 * time.Second)

  // https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go
  urlsFile, err := os.Open("urls.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer urlsFile.Close()

  sites := make([]string, 0, 0)

  scanner := bufio.NewScanner(urlsFile)
  for scanner.Scan() {
    url := scanner.Text()
    sites = append(sites, url)
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  for i := 0; i < len(sites); i++ {
    siteToExplore := sites[i]
    openbrowser(sites[i])
    fmt.Printf("Do you want to explore %s ? (yes/no/exit)\n", siteToExplore)
    var yesNoExit string
    _, err := fmt.Scanln(&yesNoExit)
    if err != nil {
      log.Fatal(err)
    }

    if yesNoExit == "yes" {
      fmt.Printf("Visiting site: %s\n", siteToExplore)
      ExploreBook(siteToExplore)
    } else if yesNoExit == "no" {
      fmt.Printf("Skipping site.\n")
    } else {
      fmt.Printf("Exiting program. Goodbye.\n")
      break;
    }
  }
}
