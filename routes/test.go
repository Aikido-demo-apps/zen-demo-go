package routes

import (
	"net/http"
	"strconv"

	"github.com/AikidoSec/firewall-go/zen"
	"github.com/gin-gonic/gin"
)

// SetupTestRoutes configures test routes for rate limiting, bot blocking, etc.
func SetupTestRoutes(r *gin.Engine) {
	r.GET("/test_ratelimiting_1", testRateLimiting1)
	r.GET("/test_ratelimiting_2", testRateLimiting2)
	r.GET("/test_bot_blocking", testBotBlocking)
	r.GET("/test_user_blocking", testUserBlocking)
	r.GET("/api/custom-event", customEvent)
}

func customEvent(c *gin.Context) {
	count, err := strconv.Atoi(c.DefaultQuery("count", "1"))
	if err != nil || count < 1 || count > 100 {
		c.String(http.StatusBadRequest, "count must be an integer between 1 and 100")
		return
	}

	for index := 0; index < count; index++ {
		zen.Track(c.Request.Context(), "user.login_failed")
	}
	c.String(http.StatusOK, "Events tracked")
}

func testRateLimiting1(c *gin.Context) {
	c.String(http.StatusOK, "Request successful (Ratelimiting 1)")
}

func testRateLimiting2(c *gin.Context) {
	c.String(http.StatusOK, "Request successful (Ratelimiting 2)")
}

func testBotBlocking(c *gin.Context) {
	c.String(http.StatusOK, "Hello World! Bot blocking enabled on this route.")
}

func testUserBlocking(c *gin.Context) {
	userID := c.GetHeader("user")
	c.String(http.StatusOK, "Hello User with id: "+userID)
}
