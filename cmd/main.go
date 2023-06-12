package main

import (
	"fmt"
	"net/http"
	"releases-parser/utils"
	"sort"

	"github.com/gin-gonic/gin"
)

func main() {
	baseURL := "https://releases.hashicorp.com"

	port := "8080"

	// Crawl
	fmt.Printf("Starting crawl across %s\n", baseURL)
	utils.CrawlLinks(baseURL)

	sort.Strings(utils.Result)

	// Serve
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/hc", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Crawled Links",
			"url":   utils.Result,
		})
	})
	r.GET("/download/*link", func(c *gin.Context) {
		link := c.Param("link")
		utils.DownloadHandler(c.Writer, c.Request, link)
	})

	r.Run(fmt.Sprintf(":%s", port))
}
