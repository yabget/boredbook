package main

import (
	"boredbook/constants"
	"boredbook/explorer"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getSitesToSkip() map[string]bool {
	skipSites := make(map[string]bool)

	skipExploreFile, err := os.Open(constants.SKIP_EXPLORE_FILENAME)
	defer skipExploreFile.Close()

	if err != nil {
		log.Printf("Could not find or open %s: %s\n", constants.SKIP_EXPLORE_FILENAME, err)
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
	f, err := os.Open(constants.BOOKMARKS_FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(constants.URLS_TO_EXPLORE_FILENAME)
	if err != nil {
		log.Fatal(err)
	}

	urlsFile, err := os.OpenFile(constants.URLS_TO_EXPLORE_FILENAME,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		for j := 0; j < len(s.Nodes); j++ {
			for k := 0; k < len(s.Nodes[j].Attr); k++ {
				href := s.Nodes[j].Attr[k].Val
				if strings.HasPrefix(href, constants.HTTP_PREFIX) {
					if _, err := urlsFile.Write([]byte(href + "\n")); err != nil {
						urlsFile.Close()
						log.Fatal(err)
					}
				}
			}
		}
		return true
	})

}

func getSitesToExplore(skipSites map[string]bool) []string {
	// https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go
	urlsFile, err := os.Open(constants.URLS_TO_EXPLORE_FILENAME)
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

func main() {

	fmt.Printf("Extracting urls from %s\n", constants.BOOKMARKS_FILENAME)
	extractURLsFromHTML()

	fmt.Printf("Gathering sites to skip during exploring.\n")
	skipSites := getSitesToSkip()

	fmt.Printf("Retrieving sites to explore.\n")
	sites := getSitesToExplore(skipSites)

	fmt.Printf("Starting exploration of sites.\n")
	explorer.ExploreSites(sites)
}
