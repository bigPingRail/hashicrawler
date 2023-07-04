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

	// Start the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", *utils.Port),
		Handler: router(result),
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
