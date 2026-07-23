package main

import (
	"time"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/jindal-parth/kavach/internal/database"
	"github.com/jindal-parth/kavach/internal/handlers"
	"github.com/jindal-parth/kavach/internal/middleware"
	"github.com/jindal-parth/kavach/internal/services"
	"github.com/jindal-parth/kavach/internal/alerts"
	"github.com/jindal-parth/kavach/internal/models"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	// Initialize database
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./kavach.db"
	}

	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("❌ Database initialization failed: %v", err)
	}
	defer database.CloseDB()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		Prefork: false,
		AppName: "KAVACH v1.0.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Serve static files
	app.Static("/static", "./static")

	// Honeypot token detection middleware - MUST BE BEFORE ROUTES
	app.Use(func(c *fiber.Ctx) error {
		// Skip health and auth endpoints
		if c.Path() == "/health" || strings.HasPrefix(c.Path(), "/api/auth") {
			return c.Next()
		}

		// Check for honeypot token in URL params
		if token := c.Query("token"); token != "" {
			if len(token) > 30 {
				log.Printf("[TOKEN-DETECTION] Found in URL param: %s...", token[:30])
			}
			if t, _ := database.GetTokenByValue(token); t != nil {
				log.Printf("[HONEYPOT-DETECTED] Valid token accessed from IP: %s | User: %s", c.IP(), t.UserID)
				createAttackerAndEvent(t.UserID, t.ID, c)
			}
		}

		// Check Authorization header
		if auth := c.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
			token := auth[7:]
			if len(token) > 30 {
				log.Printf("[TOKEN-DETECTION] Found in Authorization header: %s...", token[:30])
			}
			if t, _ := database.GetTokenByValue(token); t != nil {
				log.Printf("[HONEYPOT-DETECTED] Valid token accessed from IP: %s | User: %s", c.IP(), t.UserID)
				createAttackerAndEvent(t.UserID, t.ID, c)
			}
		}

		// Check form data
		if token := c.FormValue("token"); token != "" {
			if len(token) > 30 {
				log.Printf("[TOKEN-DETECTION] Found in form data: %s...", token[:30])
			}
			if t, _ := database.GetTokenByValue(token); t != nil {
				log.Printf("[HONEYPOT-DETECTED] Valid token accessed from IP: %s | User: %s", c.IP(), t.UserID)
				createAttackerAndEvent(t.UserID, t.ID, c)
			}
		}

		return c.Next()
	})

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		// Serve landing page
		landingHTML, err := os.ReadFile("templates/index.html")
		if err == nil {
			c.Set("Content-Type", "text/html; charset=utf-8")
			return c.SendString(string(landingHTML))
		}
		// Fallback to JSON if HTML file not found
		return c.JSON(fiber.Map{
			"name":    "KAVACH",
			"version": "1.0.0",
			"status":  "running",
			"docs":    "https://docs.kavach.security",
		})
	})

	// Page routes
	app.Get("/login", handlers.LoginPage)
	app.Get("/signup", handlers.SignupPage)
	app.Get("/products", func(c *fiber.Ctx) error {
		html, err := os.ReadFile("templates/products.html")
		if err != nil {
			return c.Status(500).SendString("Page not found")
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(string(html))
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		html, err := os.ReadFile("templates/docs.html")
		if err != nil {
			return c.Status(500).SendString("Page not found")
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(string(html))
	})
	app.Get("/vision", func(c *fiber.Ctx) error {
		html, err := os.ReadFile("templates/vision.html")
		if err != nil {
			return c.Status(500).SendString("Page not found")
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(string(html))
	})
	app.Get("/landing", handlers.LandingPage)
	app.Get("/app", handlers.DashboardPage)
	app.Get("/tokens", handlers.TokensPage)
	app.Get("/tokens/new", handlers.NewTokenPage)
	app.Get("/tokens/:tokenID", handlers.TokenDetailPage)
	app.Get("/attackers", handlers.AttackersPage)
	app.Get("/attackers/:attackerID", handlers.AttackerDetailPage)
	app.Get("/alerts", handlers.AlertsPage)
	app.Get("/integrations", handlers.IntegrationsPage)
	app.Get("/settings", handlers.SettingsPage)
	app.Get("/profile", handlers.ProfilePage)

	// New website pages
	app.Get("/how-it-works", func(c *fiber.Ctx) error {
		html, err := os.ReadFile("templates/how-it-works.html")
		if err != nil {
			return c.Status(500).SendString("Page not found")
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(string(html))
	})
	app.Get("/pricing", func(c *fiber.Ctx) error {
		html, err := os.ReadFile("templates/pricing.html")
		if err != nil {
			return c.Status(500).SendString("Page not found")
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(string(html))
	})
	app.Get("/use-cases", func(c *fiber.Ctx) error {
		html, err := os.ReadFile("templates/use-cases.html")
		if err != nil {
			return c.Status(500).SendString("Page not found")
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(string(html))
	})
	app.Get("/faq", func(c *fiber.Ctx) error {
		html, err := os.ReadFile("templates/faq.html")
		if err != nil {
			return c.Status(500).SendString("Page not found")
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(string(html))
	})
	app.Get("/support", func(c *fiber.Ctx) error {
		html, err := os.ReadFile("templates/support.html")
		if err != nil {
			return c.Status(500).SendString("Page not found")
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(string(html))
	})

	// Auth routes
	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", handlers.Register)
	authGroup.Post("/login", handlers.Login)
	authGroup.Post("/trust-device", handlers.TrustDevice)

	// Protected API routes
	apiGroup := app.Group("/api")
	apiGroup.Use(middleware.JWTAuth)

	// User routes
	apiGroup.Get("/user/profile", handlers.GetProfile)

	// Token routes
	tokenGroup := apiGroup.Group("/tokens")
	tokenGroup.Post("", handlers.CreateToken)
	tokenGroup.Get("", handlers.ListTokens)
	tokenGroup.Delete("/:tokenID", handlers.DeleteToken)
	tokenGroup.Post("/bulk", handlers.BulkCreateTokens)

	// Dashboard routes
	dashboardGroup := apiGroup.Group("/dashboard")
	dashboardGroup.Get("/stats", handlers.GetDashboardStats)
	dashboardGroup.Get("/attackers", handlers.GetAttackers)
	dashboardGroup.Get("/events", handlers.GetEvents)

	// Alert routes
	alertGroup := apiGroup.Group("/alerts")
	alertGroup.Post("/config", handlers.CreateAlertConfig)
	alertGroup.Get("/config", handlers.ListAlertConfigs)
	alertGroup.Delete("/config/:configID", handlers.DeleteAlertConfig)

	// Proxy routes
	app.Post("/api/proxy/setup", handlers.ProxySetup)
	app.All("/proxy/*", handlers.ProxyHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("\n🚀 KAVACH Server Starting\n")
	log.Printf("📍 Server: http://localhost:%s\n", port)
	log.Printf("📊 Dashboard: http://localhost:%s/app\n", port)
	log.Printf("📚 API Docs: https://docs.kavach.security\n\n", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}

// Helper function to create attacker and event records when honeypot is triggered
func createAttackerAndEvent(userID, tokenID string, c *fiber.Ctx) {
	// Generate fingerprint
	fpGen := &services.FingerprintGenerator{}
	fingerprint := fpGen.Generate(c.IP(), c.Get("User-Agent"), c.Get("Accept-Language"), c.Get("Accept-Encoding"))

	// Create attacker
	attacker, err := database.CreateOrUpdateAttacker(
		userID, fingerprint, c.IP(), c.Get("User-Agent"),
		"Unknown", "Unknown", "bot", 95,
	)
	if err != nil {
		log.Printf("[ATTACKER-ERROR] Failed to create attacker: %v", err)
		return
	}
	log.Printf("[ATTACKER-CREATED] ID: %s | IP: %s", attacker.ID, attacker.IPAddress)

	// Create trigger event
	_, err = database.CreateTriggerEvent(
		userID, tokenID, &attacker.ID,
		"token_accessed", c.Method(), c.Path(),
		"honeypot_triggered", 95,
	)
	if err != nil {
		log.Printf("[EVENT-ERROR] Failed to create event: %v", err)
		return
	}
	log.Printf("[EVENT-CREATED] Event recorded for user: %s", userID)
	
	// Get the token to send in alerts
	token, _ := database.GetTokenByID(tokenID)
	if token == nil {
		return
	}
	
	// Build and send alerts to configured endpoints
	payload := alerts.BuildWebhookPayload(userID, tokenID, token.TokenValue, token.TokenType, attacker.ID, attacker.IPAddress, 95)
	
	// Get user's alert configs and send to each
	configs, _ := database.GetAlertConfigsByUserID(userID)
	log.Printf("[ALERT-CONFIG] Found %d alert configs for user %s", len(configs), userID)
	for _, config := range configs {
		log.Printf("[ALERT-DISPATCH] Sending %s alert to %s", config.AlertType, config.Destination)
		go sendAlertToConfig(&config, payload)
	}
	time.Sleep(100 * time.Millisecond) // Give goroutines time to start logging
}

// sendAlertToConfig sends alert to a configured endpoint
func sendAlertToConfig(config *models.AlertConfig, payload *alerts.WebhookPayload) {
	dispatcher := alerts.NewAlertDispatcher()
	
	switch config.AlertType {
	case "webhook":
		log.Printf("[WEBHOOK-ATTEMPT] Posting to %s with risk: %s", config.Destination, payload.RiskLevel)
		err := dispatcher.SendWebhookAlert(config.Destination, payload)
		if err != nil {
			log.Printf("[ALERT-ERROR] Webhook failed: %v", err)
		} else {
			log.Printf("[WEBHOOK-SUCCESS] Alert sent successfully")
		}
	case "slack":
		log.Printf("[SLACK-ATTEMPT] Posting to Slack webhook")
		err := dispatcher.SendSlackAlert(config.Destination, payload)
		if err != nil {
			log.Printf("[ALERT-ERROR] Slack failed: %v", err)
		}
	case "email":
		log.Printf("[ALERT-TODO] Email alerts not yet implemented for: %s", config.Destination)
	}
}
