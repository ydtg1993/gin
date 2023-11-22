package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
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

// CustomResponseWriter is a custom response writer that captures the response body.
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write overrides the Write method to capture the response body.
func (w *responseWriter) Write(b []byte) (int, error) {
	if w.body == nil {
		w.body = &bytes.Buffer{}
	}
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// cacheEntry represents a cached HTML entry with a timestamp
type cacheEntry struct {
	Content   template.HTML `json:"html_content"`
	Timestamp time.Time     `json:"timestamp"`
}

func CacheHTMLMiddleware(timeout time.Duration, params ...string) gin.HandlerFunc {
	cacheDir := "cache"
	return func(c *gin.Context) {
		domain := c.Request.Host
		path := c.Request.URL.Path
		// Get the full URL without query string
		fullURLWithoutQuery := fmt.Sprintf("%s%s", domain, path)
		queryParams := c.Request.URL.Query()
		for _, param := range params {
			fullURLWithoutQuery += "-" + queryParams.Get(param)
		}
		hash := md5.Sum([]byte(fullURLWithoutQuery))
		cachePath := filepath.Join(cacheDir, hex.EncodeToString(hash[:]))

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

		// If not cached or expired, proceed with the request and cache the HTML page afterward
		c.Next()

		// After the request is processed, check if the response is HTML
		if c.Writer.Status() == 200 && strings.HasPrefix(c.Writer.Header().Get("Content-Type"), "text/html") {
			// Cache the HTML page
			htmlContent := w.body.String()
			err := writeCacheEntry(cachePath, htmlContent, timeout)
			if err != nil {
				fmt.Println("Failed to cache HTML:", err)
			}
		}
	}
}

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

func writeCacheEntry(cachePath, content string, timeout time.Duration) error {
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return err
	}

	entry := cacheEntry{
		Content:   template.HTML(content),
		Timestamp: time.Now().Add(timeout),
	}

	entryJSON, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cachePath, entryJSON, 0644)
}
