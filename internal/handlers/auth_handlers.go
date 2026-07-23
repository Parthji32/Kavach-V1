package handlers

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jindal-parth/kavach/internal/classifier"
	"github.com/jindal-parth/kavach/internal/database"
	"github.com/jindal-parth/kavach/internal/middleware"
	"github.com/jindal-parth/kavach/internal/models"
	"github.com/jindal-parth/kavach/internal/services"
)

// isHTMX returns true if the request was made via HTMX
func isHTMX(c *fiber.Ctx) bool {
	return c.Get("HX-Request") == "true"
}

// authErrorHTML returns a styled HTML error fragment for HTMX callers
func authErrorHTML(message string) string {
	return fmt.Sprintf(
		`<div class="p-3 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 text-sm">%s</div>`,
		message,
	)
}

// authSuccessHTML returns a styled HTML success fragment for HTMX callers
func authSuccessHTML(message string) string {
	return fmt.Sprintf(
		`<div class="p-3 rounded-lg bg-green-500/10 border border-green-500/30 text-green-400 text-sm">%s</div>`,
		message,
	)
}

// Register handles user registration
func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		if isHTMX(c) {
			return c.Status(fiber.StatusBadRequest).SendString(authErrorHTML("Invalid request. Please fill all fields."))
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
	}

	// Validate input
	if !ValidateEmail(req.Email) {
		if isHTMX(c) {
			return c.Status(fiber.StatusBadRequest).SendString(authErrorHTML("Please provide a valid email address."))
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid email",
			Message: "Please provide a valid email address",
		})
	}

	if !ValidatePassword(req.Password) {
		if isHTMX(c) {
			return c.Status(fiber.StatusBadRequest).SendString(
				authErrorHTML("Password must be at least 8 characters with uppercase, lowercase, digit, and special character."))
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Weak password",
			Message: "Password must be at least 8 characters with uppercase, lowercase, digit, and special character",
		})
	}

	// Register user
	user, err := services.RegisterUser(req.Email, req.Password, req.FullName)
	if err != nil {
		log.Printf("Registration failed: %v", err)
		if isHTMX(c) {
			return c.Status(fiber.StatusConflict).SendString(authErrorHTML("Registration failed: " + err.Error()))
		}
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Error:   "Registration failed",
			Message: err.Error(),
		})
	}

	// For HTMX: show success and redirect to login
	if isHTMX(c) {
		c.Set("HX-Redirect", "/login")
		return c.Status(fiber.StatusCreated).SendString(authSuccessHTML("Account created! Redirecting to login..."))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    user,
		"message": "User registered successfully",
	})
}

// Login handles user login with risk assessment
func Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		if isHTMX(c) {
			return c.Status(fiber.StatusBadRequest).SendString(authErrorHTML("Invalid request. Please provide email and password."))
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
	}

	// Get client details
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	// Classify login attempt using advanced classifier
	adv := classifier.NewAdvancedClassifier()
	headers := make(map[string]string)
	c.Request().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})

	classification := adv.Classify(ipAddress, userAgent, "POST", "/api/auth/login", headers, "", req.Email)

	// Log the classification
	adv.LogClassification(classification)

	// If risk very high (75+), reject immediately
	if classification.OverallRisk > 75 {
		log.Printf("Login attempt blocked for %s - risk score=%d", req.Email, classification.OverallRisk)
		if isHTMX(c) {
			return c.Status(fiber.StatusUnauthorized).SendString(
				authErrorHTML("This login attempt appears suspicious. Please try again from a known device."))
		}
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Verification required",
			Message: "This login attempt appears suspicious. Please try again from a known device.",
		})
	}

	// Authenticate user
	response, err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		log.Printf("Login failed for %s: %v", req.Email, err)
		if isHTMX(c) {
			return c.Status(fiber.StatusUnauthorized).SendString(authErrorHTML("Invalid email or password."))
		}
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Login failed",
			Message: err.Error(),
		})
	}

	// If risk moderate (50-75), add challenge flag
	if classification.OverallRisk > 50 {
		response.ChallengeRequired = true
		response.ChallengeMessage = "Please verify with a code sent to your email"
	}

	// Set HTTP-only secure cookie with the JWT token
	secure := os.Getenv("ENV") == "production"
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    response.Token,
		MaxAge:   604800, // 7 days
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
		Path:     "/",
	})

	// For HTMX: redirect to dashboard via HX-Redirect header
	if isHTMX(c) {
		c.Set("HX-Redirect", "/app")
		return c.Status(fiber.StatusOK).SendString(authSuccessHTML("Login successful! Redirecting..."))
	}

	// For JSON API consumers
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"data":     response,
		"message":  "Login successful",
		"redirect": "/app",
	})
}

// GetProfile returns current user's profile
func GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in token",
		})
	}

	// Get user from database
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "Not found",
			Message: "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// TrustDevice marks current device as trusted for user
func TrustDevice(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	fpGen := &services.FingerprintGenerator{}
	deviceFingerprint := fpGen.Generate(ipAddress, userAgent, "", "")

	// Add to trusted devices
	err := services.AddTrustedDevice(userID, deviceFingerprint)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to trust device",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Device marked as trusted",
	})
}
