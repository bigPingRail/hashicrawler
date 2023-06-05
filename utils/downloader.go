package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request, link string) {
	u := &url.URL{
		Scheme: "https",
		Host:   "releases.hashicorp.com",
		Path:   link,
	}

	buffer := new(bytes.Buffer)

	response, err := http.Get(u.String())
	if err != nil {
		fmt.Println("Failed to download file:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Failed to download file. Status code:", response.StatusCode)
		return
	}

	_, err = io.Copy(buffer, response.Body)
	if err != nil {
		fmt.Println("Failed to save file:", err)
		return
	}

	lastIndex := strings.LastIndex(link, "/")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", link[lastIndex+1:]))
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, buffer)
	if err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		return
	}
}
