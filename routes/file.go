package routes

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupFileRoutes configures file reading routes
func SetupFileRoutes(r *gin.Engine) {
	r.GET("/api/read", readFile)
	r.GET("/api/read2", readFile2)
}

func readFile(c *gin.Context) {
	path := c.Query("path")
	fullPath := filepath.Join("static/blogs/", path)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "No such file or directory") ||
			strings.Contains(errMsg, "Is a directory") ||
			strings.Contains(errMsg, "embedded null byte") {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", errMsg))
			return
		}
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", errMsg))
		return
	}

	c.String(http.StatusOK, string(content))
}

func readFile2(c *gin.Context) {
	path := c.Query("path")
	fullPath := filepath.Join("static/blogs/", path)
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "No such file or directory") ||
			strings.Contains(errMsg, "Is a directory") ||
			strings.Contains(errMsg, "embedded null byte") {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", errMsg))
			return
		}
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", errMsg))
		return
	}

	c.String(http.StatusOK, string(content))
}
