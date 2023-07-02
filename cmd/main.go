package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"releases-parser/utils"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Parse arguments
	flag.Parse()

	// Crawl
	baseURL := fmt.Sprintf("%s://%s", utils.ConnScheme, utils.ConnHost)

	go utils.StartLoadingAnimation(baseURL)
	utils.CrawlLinks(baseURL)
	utils.StopLoadingAnimation()

	sort.Strings(utils.Result)

	result := make(map[string][]string)

	for _, s := range utils.Result {
		parts := strings.Split(s, "/")
		if len(parts) >= 3 {
			key := parts[1]
			result[key] = append(result[key], s)
		}
	}

	// Create a context for the server
	ctx := context.Background()

	// Create a channel to receive the interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Create the Gin router
	r := gin.Default()
	r.SetTrustedProxies(nil)
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

	// Start the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", *utils.Port),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server failed to start: %v", err)
		}
	}()

	// Wait for the interrupt signal
	<-quit
	fmt.Println("Server shutting down...")

	// Create a context with a timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully with the given timeout
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Server shutdown failed: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
