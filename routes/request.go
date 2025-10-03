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
	r.POST("/api/request_different_port", makeRequestDifferentPort)
}

func makeRequest(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"output":  "Invalid request",
		})
		return
	}

	resp, err := http.Get(req.URL)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "Zen has blocked a server-side request forgery") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"output":  errMsg,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"output":  fmt.Sprintf("Error: %s", errMsg),
		})
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"output":  string(body),
	})
}

func makeRequestDifferentPort(c *gin.Context) {
	var req struct {
		URL  string `json:"url"`
		Port int    `json:"port"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"output":  "Invalid request",
		})
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
		if strings.Contains(errMsg, "Zen has blocked a server-side request forgery") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"output":  errMsg,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"output":  fmt.Sprintf("Error: %s", errMsg),
		})
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"output":  string(body),
	})
}
