package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var (
	resultMutex sync.Mutex
	Result      []string
)

// Public
func CrawlLinks(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ul := doc.Find("ul").First()
	links := ul.Find("a")

	linkChan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(25)

	for i := 0; i < 25; i++ {
		go func() {
			for link := range linkChan {
				absoluteURL := getAbsoluteURL(url, link)
				if !isRelativeURL(absoluteURL) {
					if strings.HasPrefix(link, url) {
						saveLinkToMemory(link, "https://releases.hashicorp.com")
					} else {
						CrawlLinks(absoluteURL)
					}
				}
			}
			wg.Done()
		}()
	}

	links.Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && href != "../" {
			linkChan <- href
		}
	})

	close(linkChan)
	wg.Wait()
}

// Private
func saveLinkToMemory(link, baseUrl string) {
	resultMutex.Lock()
	defer resultMutex.Unlock()

	if !strings.HasSuffix(link, "HA256SUMS") && !strings.HasSuffix(link, ".sig") &&
		(strings.Contains(link, "_linux") || strings.Contains(link, "_darwin") || strings.Contains(link, "_windows")) {
		link := strings.TrimPrefix(link, baseUrl)
		Result = append(Result, link)
	}
}

func getAbsoluteURL(baseURL, href string) string {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Error parsing base URL:", err)
		return ""
	}

	relative, err := url.Parse(href)
	if err != nil {
		fmt.Println("Error parsing relative URL:", err)
		return ""
	}

	absoluteURL := base.ResolveReference(relative).String()
	return absoluteURL
}

func isRelativeURL(urlString string) bool {
	u, err := url.Parse(urlString)
	if err != nil {
		return false
	}

	return !u.IsAbs()
}
