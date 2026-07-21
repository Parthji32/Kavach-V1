package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jindal-parth/kavach/internal/classifier"
	"github.com/jindal-parth/kavach/internal/database"
	"github.com/jindal-parth/kavach/internal/middleware"
	"github.com/jindal-parth/kavach/internal/models"
	"github.com/jindal-parth/kavach/internal/services"
)

// Register handles user registration
func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
	}

	// Validate input
	if !ValidateEmail(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid email",
			Message: "Please provide a valid email address",
		})
	}
	
	if !ValidatePassword(req.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Weak password",
			Message: "Password must be at least 8 characters with uppercase, lowercase, digit, and special character",
		})
	}

	// Register user
	user, err := services.RegisterUser(req.Email, req.Password, req.FullName)
	if err != nil {
		log.Printf("Registration failed: %v", err)
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Error:   "Registration failed",
			Message: err.Error(),
		})
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
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Verification required",
			Message: "This login attempt appears suspicious. Please try again from a known device.",
		})
	}

	// Authenticate user
	response, err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		log.Printf("Login failed for %s: %v", req.Email, err)
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

	c.Cookie(&fiber.Cookie{Name: "token", Value: response.Token, MaxAge: 604800, HTTPOnly: false, Secure: false})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Login successful",
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
