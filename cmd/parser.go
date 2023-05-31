package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var (
	resultMutex sync.Mutex
	result      []string
)

func crawlLinks(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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
						saveLinkToMemory(link)
					} else {
						crawlLinks(absoluteURL)
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

func saveLinkToMemory(link string) {
	resultMutex.Lock()
	defer resultMutex.Unlock()

	if !strings.HasSuffix(link, "HA256SUMS") && !strings.HasSuffix(link, ".sig") &&
		(strings.Contains(link, "_linux") || strings.Contains(link, "_darwin") || strings.Contains(link, "_windows")) {
		result = append(result, link)
	}
}

func writeLinksToFile(filename string, links []string) error {
	baseHTML := "<a href={{.}}>{{.}}</a><br>"

	var htmlStrings []string

	for _, link := range links {
		html := strings.Replace(baseHTML, "{{.}}", link, -1)
		htmlStrings = append(htmlStrings, html)
	}

	outputHTML := strings.Join(htmlStrings, "\n")

	err := os.WriteFile(filename, []byte(outputHTML), 0644)
	if err != nil {
		return err
	}

	return nil
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

func main() {
	baseURL := "https://releases.hashicorp.com/"
	output := "index.html"
	port := ":8080"

	fmt.Printf("Begin to parse %s\n", baseURL)
	crawlLinks(baseURL)

	sort.Strings(result)

	if err := writeLinksToFile(output, result); err != nil {
		log.Fatalf("Error writing links to file: %s", err)
	}

	fmt.Printf("Done!\nLinks written to file: %s\n", output)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, output)
	})

	http.HandleFunc("/hc", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	fmt.Printf("Server listening on port %s...\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
