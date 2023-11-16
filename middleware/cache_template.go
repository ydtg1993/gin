package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"xo/core"

	"github.com/gin-gonic/gin"
)

// CustomResponseWriter is a custom response writer that captures the response body.
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write overrides the Write method to capture the response body.
func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	if w.body == nil {
		w.body = &bytes.Buffer{}
	}
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// CacheEntry represents a cached HTML entry with a timestamp
type CacheEntry struct {
	Content   template.HTML `json:"html_content"`
	Timestamp time.Time     `json:"timestamp"`
}

// CacheHTMLMiddleware is a Gin middleware that caches HTML pages to a cache directory
// and serves the cached page on subsequent requests with an expiration time.
func CacheHTMLMiddleware() gin.HandlerFunc {
	cacheDir := core.Config.GetString("template.cache_dirname")
	return func(c *gin.Context) {
		// Generate a unique cache key based on the request path
		cacheKey := strings.Replace(c.Request.URL.Path, "/", "_", -1)
		cachePath := filepath.Join(cacheDir, cacheKey+".html")

		// Create a custom response writer to capture the response body
		w := &CustomResponseWriter{c.Writer, nil}
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
					"title":  core.Config.GetString("app.name"),
					"static": template.HTML(staticContent),
					"error":  "page cannot be found"})
				c.Abort()
			}
		}()
		// If not cached or expired, proceed with the request and cache the HTML page afterward
		c.Next()

		// After the request is processed, check if the response is HTML
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

// readCacheEntry reads a cached entry from the file system
func readCacheEntry(cachePath string) (*CacheEntry, error) {
	fileContent, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var entry CacheEntry
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

	entry := CacheEntry{
		Content:   template.HTML(content),
		Timestamp: time.Now().Add(core.Config.GetDuration("template.cache_expire") * time.Minute),
	}

	entryJSON, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cachePath, entryJSON, 0644)
}
