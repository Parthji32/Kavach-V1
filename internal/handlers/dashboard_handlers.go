package handlers

import (
	"strings"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jindal-parth/kavach/internal/database"
	"github.com/jindal-parth/kavach/internal/middleware"
	"github.com/jindal-parth/kavach/internal/models"
	"github.com/jindal-parth/kavach/internal/services"
)

// GetDashboardStats returns dashboard statistics
func GetDashboardStats(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	// Check for honeypot token in URL, headers, or body
	detectedToken := detectHoneypotToken(c)
	if detectedToken != "" {
		log.Printf("[HONEYPOT] Token detected in request: %s", detectedToken)

		// Look up token
		token, err := database.GetTokenByValue(detectedToken)
		if err == nil && token != nil && token.IsActive {
			log.Printf("[HONEYPOT TRIGGERED] Token: %s | User: %s", token.ID, userID)

			// Generate fingerprint
			fpGen := &services.FingerprintGenerator{}
			fingerprint := fpGen.Generate(c.IP(), c.Get("User-Agent"), c.Get("Accept-Language"), c.Get("Accept-Encoding"))

			// Create attacker record
			attacker, err := database.CreateOrUpdateAttacker(
				userID,
				fingerprint,
				c.IP(),
				c.Get("User-Agent"),
				"Unknown",
				"Unknown",
				"bot",
				95, // High risk score
			)
			if err == nil {
				log.Printf("[ATTACKER RECORDED] ID: %s | IP: %s", attacker.ID, attacker.IPAddress)

				// Create trigger event with correct function signature
				_, err := database.CreateTriggerEvent(
					userID,
					token.ID,
					&attacker.ID,
					"token_accessed",
					c.Method(),
					c.Path(),
					"honeypot_triggered",
					95,
				)
				if err != nil {
					log.Printf("[EVENT ERROR] Failed to create event: %v", err)
				}
			}
		}
	}

	// Get token stats
	activeTokens, err := database.GetActiveTokensCount(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to fetch dashboard data",
			Message: "Could not retrieve active tokens count",
		})
	}

	tokens, err := database.GetTokensByUserID(userID, 50, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to fetch dashboard data",
			Message: "Could not retrieve tokens",
		})
	}

	// Get attacker stats
	attackerCount, err := database.GetAttackerCount(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to fetch dashboard data",
			Message: "Could not retrieve attacker count",
		})
	}

	highRiskAttackers, err := database.GetHighRiskAttackers(userID, 70)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to fetch dashboard data",
			Message: "Could not retrieve high-risk attackers",
		})
	}

	// Get recent events
	eventsLast24h, err := database.GetEventsLast24h(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to fetch dashboard data",
			Message: "Could not retrieve recent events count",
		})
	}

	recentEvents, err := database.GetTriggerEventsByUserID(userID, 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to fetch dashboard data",
			Message: "Could not retrieve recent events",
		})
	}

	// Get recent attackers
	recentAttackers, err := database.GetAttackersByUserID(userID, 5)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to fetch dashboard data",
			Message: "Could not retrieve recent attackers",
		})
	}

	// Build response
	stats := models.DashboardStats{
		TotalTokens:    len(tokens),
		ActiveTokens:   activeTokens,
		TotalAttackers: attackerCount,
		HighRiskCount:  len(highRiskAttackers),
		EventsLast24h:  eventsLast24h,
	}

	// Build recent attackers response
	for _, attacker := range recentAttackers {
		stats.RecentAttackers = append(stats.RecentAttackers, models.AttackerWithRisk{
			ID:        attacker.ID,
			IPAddress: attacker.IPAddress,
			RiskScore: int(attacker.RiskScore),
			LastSeen:  attacker.LastSeen.Format(time.RFC3339),
			IsBlocked: attacker.IsBlocked,
		})
	}

	// Build recent events response
	for _, event := range recentEvents {
		token, err := database.GetTokenByID(event.TokenID)
		tokenType := "unknown"
		if err == nil && token != nil && token.TokenType != "" {
			tokenType = token.TokenType
		}

		stats.RecentEvents = append(stats.RecentEvents, models.TriggerEventResponse{
			ID:        event.ID,
			TokenType: tokenType,
			EventType: event.EventType,
			Timestamp: event.Timestamp.Format(time.RFC3339),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    stats,
		"message": "Dashboard stats retrieved successfully",
	})
}

// GetAttackers returns list of attackers
func GetAttackers(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	// Apply pagination
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	// Validate pagination parameters
	if limit < 1 || limit > 500 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	// Get attackers from database
	attackers, err := database.GetAttackersByUserID(userID, limit)
	if err != nil {
		log.Printf("Failed to get attackers: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to fetch attackers",
			Message: err.Error(),
		})
	}

	total := len(attackers)
	end := offset + limit
	if end > total {
		end = total
	}
	paginatedAttackers := attackers[offset:end]

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    paginatedAttackers,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"message": "Attackers retrieved successfully",
	})
}

// GetEvents returns trigger events
func GetEvents(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	// Apply pagination
	limit := c.QueryInt("limit", 100)
	offset := c.QueryInt("offset", 0)

	// Validate pagination parameters
	if limit < 1 || limit > 500 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	// Get events from database
	events, err := database.GetTriggerEventsByUserID(userID, limit)
	if err != nil {
		log.Printf("Failed to get events: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to fetch events",
			Message: err.Error(),
		})
	}

	total := len(events)
	end := offset + limit
	if end > total {
		end = total
	}
	paginatedEvents := events[offset:end]

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    paginatedEvents,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"message": "Events retrieved successfully",
	})
}

// detectHoneypotToken looks for honeypot tokens in request (URL, headers, body)
func detectHoneypotToken(c *fiber.Ctx) string {
	// Check URL query parameters
	if token := c.Query("token"); token != "" && strings.HasPrefix(token, "sk_") {
		return token
	}

	// Check Authorization header for token patterns
	if auth := c.Get("Authorization"); auth != "" {
		if strings.Contains(auth, "sk_") {
			parts := strings.Fields(auth)
			for _, part := range parts {
				if strings.HasPrefix(part, "sk_") && len(part) > 10 {
					return part
				}
			}
		}
	}

	// Check common form field names
	if token := c.FormValue("token"); token != "" && strings.HasPrefix(token, "sk_") {
		return token
	}
	if token := c.FormValue("api_key"); token != "" && strings.HasPrefix(token, "sk_") {
		return token
	}

	return ""
}
