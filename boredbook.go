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

  skipSites := make(map[string]bool)

  skipExploreFile, err := os.Open("skipExplore.txt")
  if err != nil {
    log.Printf("Could not find or open skipExplore.txt: %s\n", err)
  } else {
    scanner := bufio.NewScanner(skipExploreFile)
    for scanner.Scan() {
      siteToSkip := scanner.Text()
      skipSites[siteToSkip] = true
    }

    if err := scanner.Err(); err != nil {
      log.Fatal(err)
    }
  }
  defer skipExploreFile.Close()

  // https://gosamples.dev/read-user-input/
  fmt.Printf("Extracting urls from bookmarks.html\n")

  err = exec.Command("./extractURLs.sh", "bookmarks.html").Start()
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
    if !skipSites[url] {
      sites = append(sites, url)
    }
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  // https://pkg.go.dev/os#example-OpenFile-Append
  skipNextTimeFile, err := os.OpenFile("skipExplore.txt",
    os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    log.Fatal(err)
  }

  for i := 0; i < len(sites); i++ {
    siteToExplore := sites[i]

    fmt.Printf("Do you want to explore %s ? (yes/no/skip/exit)\n", siteToExplore)
    var yesNoSkipExit string
    _, err := fmt.Scanln(&yesNoSkipExit)
    if err != nil {
      log.Fatal(err)
    }

    if yesNoSkipExit == "yes" {
      openbrowser(sites[i])
      fmt.Printf("Visiting site: %s\n", siteToExplore)
      ExploreBook(siteToExplore)
    } else if yesNoSkipExit == "no" {
      fmt.Printf("Skipping site.\n")
    } else if yesNoSkipExit == "skip" {
      if _, err := skipNextTimeFile.Write([]byte(siteToExplore + "\n")); err != nil {
        skipNextTimeFile.Close()
        log.Fatal(err)
      }
    } else {
      fmt.Printf("Exiting program. Goodbye.\n")
      break;
    }
  }
  if err := skipNextTimeFile.Close(); err != nil {
    log.Fatal(err)
  }
}
