package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupRequestRoutes configures SSRF/request routes
func SetupRequestRoutes(r *gin.Engine) {
	r.POST("/api/request", makeRequest)
	r.POST("/api/request2", makeRequest2)
	r.POST("/api/request_different_port", makeRequestDifferentPort)
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
		errMsg := err.Error()
		if strings.Contains(errMsg, "Failed to resolve") {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", errMsg))
			return
		}
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", errMsg))
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

func makeRequestDifferentPort(c *gin.Context) {
	var req struct {
		URL  string `json:"url"`
		Port int    `json:"port"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	// Replace port in URL
	parts := strings.Split(req.URL, ":")
	if len(parts) >= 3 {
		parts[len(parts)-1] = strconv.Itoa(req.Port)
	}
	urlWithPort := strings.Join(parts, ":")

	resp, err := http.Get(urlWithPort)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "Failed to resolve") {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", errMsg))
			return
		}
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", errMsg))
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
		errMsg := err.Error()
		if strings.Contains(errMsg, "Failed to resolve") {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", errMsg))
			return
		}
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", errMsg))
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
