package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

type cacheEntry struct {
	Content   template.HTML `json:"html_content"`
	Timestamp time.Time     `json:"timestamp"`
}

func TemplateCache(c *gin.Context, timeout int) func(func()) {
	return func(f func()) {
		cacheDir := "cache"
		cachePath := filepath.Join(cacheDir, ".html")
		// Create a custom response writer to capture the response body
		w := &responseWriter{c.Writer, nil}
		c.Writer = w

		// Check if the HTML page is already cached
		cacheEntry, err := readCacheEntry(cachePath)
		if err == nil && cacheEntry.Timestamp.After(time.Now()) {
			// If cached and not expired, serve the cached HTML page
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, string(cacheEntry.Content))
			c.Abort()
			return
		}
		defer func() {
			if r := recover(); r != nil {
				// Handle panic during c.Next(): return an error page
				staticContent, _ := ioutil.ReadFile("resources/templates/static.html")
				c.HTML(http.StatusInternalServerError, "error.html", gin.H{
					"title":  Config.GetString("app.name"),
					"static": template.HTML(staticContent),
					"error":  "page cannot be found"})
				c.Abort()
			}
		}()
		f()
		if c.Writer.Status() == 200 && strings.HasPrefix(c.Writer.Header().Get("Content-Type"), "text/html") {
			// Cache the HTML page
			htmlContent := w.body.String()
			err := writeCacheEntry(cachePath, htmlContent)
			if err != nil {
				fmt.Println("Failed to cache HTML:", err)
			}
		}
	}
}

// Write overrides the Write method to capture the response body.
func (w *responseWriter) Write(b []byte) (int, error) {
	if w.body == nil {
		w.body = &bytes.Buffer{}
	}
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// readCacheEntry reads a cached entry from the file system
func readCacheEntry(cachePath string) (*cacheEntry, error) {
	fileContent, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var entry cacheEntry
	if err := json.Unmarshal(fileContent, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// writeCacheEntry writes a cached entry to the file system
func writeCacheEntry(cachePath, content string) error {
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return err
	}

	entry := cacheEntry{
		Content:   template.HTML(content),
		Timestamp: time.Now().Add(Config.GetDuration("template.cache_expire") * time.Minute),
	}

	entryJSON, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cachePath, entryJSON, 0644)
}
