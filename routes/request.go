package routes

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var storedSsrfURLs = make(chan string, 100)

func init() {
	go func() {
		for url := range storedSsrfURLs {
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error fetching %s: %s\n", url, err.Error())
				continue
			}
			resp.Body.Close()
		}
	}()
}

// SetupRequestRoutes configures SSRF/request routes
func SetupRequestRoutes(r *gin.Engine) {
	r.POST("/api/request", makeRequest)
	r.POST("/api/request2", makeRequest2)
	r.POST("/api/request_different_port", makeRequestDifferentPort)
	r.POST("/api/stored_ssrf", storedSSRF)
	r.POST("/api/stored_ssrf_2", storedSSRF2)
}

func makeRequest(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	resp, err := http.Get(req.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func makeRequestDifferentPort(c *gin.Context) {
	var req struct {
		URL  string `json:"url"`
		Port string `json:"port"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	// Replace port in URL
	parts := strings.Split(req.URL, ":")
	if len(parts) >= 3 {
		parts[len(parts)-1] = req.Port
	}
	urlWithPort := strings.Join(parts, ":")

	resp, err := http.Get(urlWithPort)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
		return
	}
	defer resp.Body.Close()

	body := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			body = append(body, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	c.String(http.StatusOK, string(body))
}

func makeRequest2(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	resp, err := http.Get(req.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
		return
	}
	defer resp.Body.Close()

	body := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			body = append(body, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	c.String(http.StatusOK, string(body))
}

func storedSSRF(c *gin.Context) {
	var req struct {
		URLIndex *int `json:"urlIndex"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	urls := []string{
		"http://evil-stored-ssrf-hostname/latest/api/token",
		"http://metadata.google.internal/latest/api/token",
		"http://metadata.goog/latest/api/token",
		"http://169.254.169.254/latest/api/token",
	}

	index := 0
	if req.URLIndex != nil {
		index = *req.URLIndex
	}
	url := urls[index%len(urls)]

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "output": err.Error()})
		return
	}
	defer resp.Body.Close()

	body := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			body = append(body, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "output": string(body)})
}

func storedSSRF2(c *gin.Context) {
	storedSsrfURLs <- "http://evil-stored-ssrf-hostname/latest/api/token"
	c.JSON(http.StatusOK, gin.H{"success": true, "output": "Request successful (Stored SSRF 2)"})
}
