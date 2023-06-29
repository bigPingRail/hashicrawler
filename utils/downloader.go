package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"log"
)

// Private
func downloadFile(u *url.URL) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	response, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file. Status code: %d", response.StatusCode)
	}

	_, err = io.Copy(buffer, response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}
	return buffer, nil
}

func writeToFile(buffer *bytes.Buffer, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return err
	}

	return nil
}

func sendFile(w http.ResponseWriter, filename string, buffer *bytes.Buffer) bool {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err := io.Copy(w, buffer)
	if err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		return true
	}
	return false
}

// Public
func DownloadHandler(w http.ResponseWriter, r *http.Request, link string) {
	u := &url.URL{
		Scheme: ConnScheme,
		Host:   ConnHost,
		Path:   link,
	}
	lastIndex := strings.LastIndex(link, "/")
	filename := link[lastIndex+1:]

	if !*Caching {
		buffer, err := downloadFile(u)
		if err != nil {
			fmt.Println(err)
		}
		sendFile(w, filename, buffer)
	} else {
		cwd, _ := os.Getwd()
		cacheDir := filepath.Join(cwd, "cache")
		os.MkdirAll(cacheDir, 0755)

		// Check cache for file
		_, err := os.Stat(filepath.Join(cacheDir, filename))
		if err != nil {
			buffer, err := downloadFile(u)
			if err != nil {
				fmt.Printf("download error: %s", err)
				return
			}
			writeToFile(buffer, filepath.Join(cacheDir, filename))
			result := sendFile(w, filename, buffer)
			if result {
				return
			}
		}

		fileBytes, err := os.ReadFile(filepath.Join(cacheDir, filename))
		if err != nil {
			log.Printf("file read error, %s", err)
			return
		}
		sendFile(w, filename, bytes.NewBuffer(fileBytes))
	}
}
