package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jindal-parth/kavach/internal/database"
	"github.com/jindal-parth/kavach/internal/middleware"
	"github.com/jindal-parth/kavach/internal/models"
)

// CreateAlertConfig creates a new alert configuration
func CreateAlertConfig(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	var req models.CreateAlertConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
	}

	// Validate input
	if !ValidateURL(req.Destination) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid destination",
			Message: "Destination must be a valid HTTPS URL",
		})
	}
	
	// If Slack, validate slack webhook format
	if req.AlertType == "slack" && !ValidateSlackWebhook(req.Destination) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid Slack webhook",
			Message: "Must be a valid Slack webhook URL (hooks.slack.com)",
		})
	}

	// Validate alert type
	validTypes := map[string]bool{"webhook": true, "email": true, "slack": true}
	if !validTypes[req.AlertType] {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid alert type",
			Message: "Supported types: webhook, email, slack",
		})
	}

	// Create alert config
	config, err := database.CreateAlertConfig(userID, req.AlertType, req.Destination)
	if err != nil {
		log.Printf("Failed to create alert config: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to create alert config",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    config,
		"message": "Alert config created successfully",
	})
}

// ListAlertConfigs returns all alert configs for user
func ListAlertConfigs(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	configs, err := database.GetAlertConfigsByUserID(userID)
	if err != nil {
		log.Printf("Failed to list alert configs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to fetch alert configs",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    configs,
		"message": "Alert configs retrieved successfully",
	})
}

// DeleteAlertConfig deletes an alert configuration
func DeleteAlertConfig(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	configID := c.Params("configID")
	if configID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Config ID is required",
		})
	}

	// Delete alert config
	if err := database.DeleteAlertConfig(configID); err != nil {
		log.Printf("Failed to delete alert config: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to delete alert config",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Alert config deleted successfully",
	})
}
