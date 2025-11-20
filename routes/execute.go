package routes

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

// SetupExecuteRoutes configures shell command execution routes
func SetupExecuteRoutes(r *gin.Engine) {
	r.POST("/api/execute", executeCommandPost)
	r.GET("/api/execute/:command", executeCommandGet)
}

func executeCommandPost(c *gin.Context) {
	var req struct {
		UserCommand string `json:"userCommand"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	result, err := executeShellCommand(req.UserCommand)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		return
	}
	c.String(http.StatusOK, result)
}

func executeCommandGet(c *gin.Context) {
	command := c.Param("command")
	result, err := executeShellCommand(command)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		return
	}
	c.String(http.StatusOK, result)
}

func executeShellCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
