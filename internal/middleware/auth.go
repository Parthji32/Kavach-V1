package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jindal-parth/kavach/internal/models"
	"github.com/jindal-parth/kavach/internal/services"
)

// JWTAuth middleware validates JWT tokens for API routes.
// Returns 401 JSON on failure (for fetch/XHR callers).
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
		// If this is an HTMX request, return an HTML error fragment
		if c.Get("HX-Request") == "true" {
			c.Set("Content-Type", "text/html; charset=utf-8")
			return c.Status(fiber.StatusUnauthorized).SendString(
				`<div class="p-3 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 text-sm">` +
					`<strong>Session expired.</strong> Please <a href="/login" class="underline text-red-300">sign in</a> again.` +
					`</div>`)
		}
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Missing authorization token",
		})
	}

	// Validate token
	claims, err := services.ValidateJWT(tokenString)
	if err != nil {
		if c.Get("HX-Request") == "true" {
			c.Set("Content-Type", "text/html; charset=utf-8")
			return c.Status(fiber.StatusUnauthorized).SendString(
				`<div class="p-3 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 text-sm">` +
					`<strong>Invalid or expired token.</strong> Please <a href="/login" class="underline text-red-300">sign in</a> again.` +
					`</div>`)
		}
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

// JWTAuthPage middleware validates JWT tokens for page routes.
// Redirects to /login on failure (for browser navigation).
func JWTAuthPage(c *fiber.Ctx) error {
	// Get token from cookie (pages use cookie auth)
	tokenString := c.Cookies("token")

	// Fallback: check Authorization header
	if tokenString == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}
	}

	if tokenString == "" {
		return c.Redirect("/login")
	}

	// Validate token
	claims, err := services.ValidateJWT(tokenString)
	if err != nil {
		// Clear invalid cookie
		c.ClearCookie("token")
		return c.Redirect("/login")
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
