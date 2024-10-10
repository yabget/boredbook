package explorer

import (
	"boredbook/browser"
	"boredbook/constants"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
)

func ExploreSites(sites []string) {
	// https://pkg.go.dev/os#example-OpenFile-Append
	skipNextTimeFile, err := os.OpenFile(constants.SKIP_EXPLORE_FILENAME,
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
			browser.Open(siteToExplore)
			ExploreSite(siteToExplore)
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

// https://pkg.go.dev/github.com/PuerkitoBio/goquery
func ExploreSite(url string) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		log.Printf("Skipping %s", url)
		return
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
				if strings.HasPrefix(href, constants.HTTP_PREFIX) && !openedSite[href] {
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

					browser.Open(href)
					openedSite[href] = true
					openedSitesCount++
					continue OUTTER
				}
			}
		}
		return true
	})
}
