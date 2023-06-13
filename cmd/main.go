package main

import (
	"fmt"
	"net/http"
	"releases-parser/utils"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	baseURL := fmt.Sprintf("%s://%s", utils.ConnScheme, utils.ConnHost)
	port := "8080"

	// Crawl
	fmt.Printf("Starting crawl across %s\n", baseURL)
	utils.CrawlLinks(baseURL)

	sort.Strings(utils.Result)

	result := make(map[string][]string)

	for _, s := range utils.Result {
		parts := strings.Split(s, "/")
		if len(parts) >= 3 {
			key := parts[1]
			result[key] = append(result[key], s)
		}
	}

	// Serve
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")

	r.GET("/hc", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/", func(c *gin.Context) {
		keys := make([]string, 0, len(result))
		for key := range result {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Keys": keys,
		})
	})

	r.GET("/values/:key", func(c *gin.Context) {
		key := c.Param("key")
		values, exists := result[key]
		if !exists {
			c.String(http.StatusNotFound, "Key not found")
			return
		}

		c.HTML(http.StatusOK, "values.tmpl", gin.H{
			"Title":  key,
			"Values": values,
		})
	})

	r.GET("/download/*link", func(c *gin.Context) {
		link := c.Param("link")
		utils.DownloadHandler(c.Writer, c.Request, link)
	})

	r.Run(fmt.Sprintf(":%s", port))
}
