package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
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
func exploreSite(url string) {
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

	openedSitesCount := 0
	// Find the review items
	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {

	OUTTER:
		for j := 0; j < len(s.Nodes); j++ {
			for k := 0; k < len(s.Nodes[j].Attr); k++ {
				//fmt.Printf("SAttr %d: %s\n", k, s.Nodes[j].Attr[k])
				href := s.Nodes[j].Attr[k].Val
				if strings.HasPrefix(href, "http") && !openedSite[href] {
					if openedSitesCount%5 == 0 {
						fmt.Printf(
							"You have explored %d sites, do you want to open the next 5? (yes/no)\n",
							openedSitesCount)

						var yesNo string
						_, err := fmt.Scanln(&yesNo)
						if err != nil {
							log.Fatal(err)
						}

						if yesNo == "yes" {
							// continue opening sites to explore
						} else {
							// exit exploring site
							return false
						}
					}

					openbrowser(href)
					openedSite[href] = true
					openedSitesCount++
					continue OUTTER
				}
			}
		}
		return true
	})
}

func getSitesToSkip() map[string]bool {
	skipSites := make(map[string]bool)

	skipExploreFile, err := os.Open("skipExplore.txt")
	defer skipExploreFile.Close()

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

	return skipSites
}

func extractURLsFromHTML() {
	err := exec.Command("./extractURLs.sh", "bookmarks.html").Start()
	if err != nil {
		log.Fatal(err)
	}

	// sleep to let the script run before resuming program
	time.Sleep(2 * time.Second)
}

func getSitesToExplore(skipSites map[string]bool) []string {
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

	return sites
}

func exploreSites(sites []string) {
	// https://pkg.go.dev/os#example-OpenFile-Append
	skipNextTimeFile, err := os.OpenFile("skipExplore.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

EndExplore:
	for i := 0; i < len(sites); i++ {
		siteToExplore := sites[i]

		// https://gosamples.dev/read-user-input/
		fmt.Printf("Do you want to explore %s ? (yes/no/skip/exit)\n", siteToExplore)

		var yesNoSkipExit string
		_, err := fmt.Scanln(&yesNoSkipExit)
		if err != nil {
			log.Fatal(err)
		}

		switch yesNoSkipExit {
		case "yes":
			fmt.Printf("Visiting site: %s\n", siteToExplore)
			openbrowser(siteToExplore)
			exploreSite(siteToExplore)
		case "no":
			fmt.Printf("Not exploring site.\n")
		case "skip":
			if _, err := skipNextTimeFile.Write([]byte(siteToExplore + "\n")); err != nil {
				skipNextTimeFile.Close()
				log.Fatal(err)
			}
		default:
			fmt.Printf("Exiting program. Goodbye.\n")
			break EndExplore
		}

	}
	if err := skipNextTimeFile.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	fmt.Printf("Extracting urls from bookmarks.html\n")
	extractURLsFromHTML()

	fmt.Printf("Gathering sites to skip during exploring.\n")
	skipSites := getSitesToSkip()

	fmt.Printf("Retrieving sites to explore.\n")
	sites := getSitesToExplore(skipSites)

	fmt.Printf("Starting exploration of sites.\n")
	exploreSites(sites)
}
