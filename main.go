package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"

	"zen-demo-go/database"
	"zen-demo-go/routes"

	"github.com/AikidoSec/firewall-go/zen"
	"github.com/gin-gonic/gin"
)

func main() {
	zen.Protect()

	r := gin.Default()
	r.ContextWithFallback = true

	// Initialize database
	database.InitDatabase()

	r.Use(func(c *gin.Context) {
		if c.GetHeader("user") != "" {
			zen.SetUser(c, c.GetHeader("user"), "John Doe")
		}
	})
	r.Use(BlockingMiddleware())

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

	// Setup pprof routes on localhost:6060
	pprofMux := http.NewServeMux()
	pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
	pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Start pprof on localhost:6060 with its own mux
	go func() {
		log.Println("pprof available at http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", pprofMux); err != nil {
			log.Fatal("pprof server failed:", err)
		}
	}()

	// Start server
	fmt.Println("Server is running on http://localhost:3000")
	r.Run(":3000")
}

func BlockingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		blockResult := zen.ShouldBlockRequest(c)

		if blockResult != nil {
			if blockResult.Type == "rate-limited" {
				message := "You are rate limited by Zen."
				if blockResult.Trigger == "ip" {
					message += " (Your IP: " + *blockResult.IP + ")"
				}
				c.String(http.StatusTooManyRequests, message)
				c.Abort() // Stop further processing
				return
			} else if blockResult.Type == "blocked" {
				c.String(http.StatusForbidden, "You are blocked by Zen.")
				c.Abort() // Stop further processing
				return
			}
		}

		c.Next()
	}
}
