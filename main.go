package main

import (
	"fmt"

	"zen-demo-go/database"
	"zen-demo-go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize database
	database.InitDatabase()

	// HTML routes
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/pages/create", func(c *gin.Context) {
		c.File("./static/create.html")
	})

	r.GET("/pages/execute", func(c *gin.Context) {
		c.File("./static/execute_command.html")
	})

	r.GET("/pages/read", func(c *gin.Context) {
		c.File("./static/read_file.html")
	})

	r.GET("/pages/request", func(c *gin.Context) {
		c.File("./static/request.html")
	})

	// Setup API routes
	routes.SetupTestRoutes(r)
	routes.SetupRequestRoutes(r)
	routes.SetupFileRoutes(r)
	routes.SetupPetRoutes(r)
	routes.SetupExecuteRoutes(r)

	r.Static("/css", "./static/public/css")
	r.Static("/js", "./static/public/js")
	r.Static("/New-Grotesk", "./static/public/New-Grotesk")
	r.GET("/aikido_logo.svg", func(c *gin.Context) {
		c.File("./static/public/aikido_logo.svg")
	})

	// Start server
	fmt.Println("Server is running on http://localhost:3000")
	r.Run(":3000")
}
