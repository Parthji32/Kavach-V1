package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jindal-parth/kavach/internal/models"
	"github.com/jindal-parth/kavach/internal/services"
)

// JWTAuth middleware validates JWT tokens
func JWTAuth(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	var tokenString string

	if authHeader != "" {
		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		}
	}

	// If not in header, check cookie
	if tokenString == "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Missing authorization token",
		})
	}

	// Validate token
	claims, err := services.ValidateJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid or expired token",
		})
	}

	// Store user info in locals for downstream handlers
	c.Locals("userID", claims.UserID)
	c.Locals("email", claims.Email)

	return c.Next()
}

// GetUserID extracts user ID from JWT in request context
func GetUserID(c *fiber.Ctx) string {
	userID := c.Locals("userID")
	if userID == nil {
		return ""
	}
	return userID.(string)
}

// GetEmail extracts email from JWT in request context
func GetEmail(c *fiber.Ctx) string {
	email := c.Locals("email")
	if email == nil {
		return ""
	}
	return email.(string)
}
