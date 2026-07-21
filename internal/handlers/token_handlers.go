package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jindal-parth/kavach/internal/database"
	"github.com/jindal-parth/kavach/internal/middleware"
	"github.com/jindal-parth/kavach/internal/models"
	"github.com/jindal-parth/kavach/internal/services"
)

// CreateToken creates a new honeypot token
func CreateToken(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	var req models.CreateTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
	}

	// Validate input
	if !ValidateTokenType(req.TokenType) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid token type",
			Message: "Allowed types: url, api_key, document, dns, email",
		})
	}
	// Generate token value
	tokenValue, err := services.GenerateTokenValue(req.TokenType)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid token type",
			Message: err.Error(),
		})
	}

	// Create token in database
	token, err := database.CreateToken(userID, req.TokenType, tokenValue, req.Description)
	if err != nil {
		log.Printf("Failed to create token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to create token",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    token,
		"message": "Token created successfully",
	})
}

// ListTokens retrieves all tokens for current user
func ListTokens(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	// Get pagination parameters
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)
	
	// Validate pagination parameters
	var err error
	limit, offset = ValidatePaginationParams(limit, offset)

	// Validate pagination parameters
	if limit < 1 || limit > 500 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	tokens, err := database.GetTokensByUserID(userID, limit, offset)
	if err != nil {
		log.Printf("Failed to list tokens: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to list tokens",
			Message: err.Error(),
		})
	}

	// Apply pagination
	total := len(tokens)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	paginatedTokens := tokens[offset:end]

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    paginatedTokens,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"message": "Tokens retrieved successfully",
	})
}

// DeleteToken deactivates a token
func DeleteToken(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	tokenID := c.Params("tokenID")
	if tokenID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Token ID is required",
		})
	}

	// Deactivate token
	if err := database.DeactivateToken(tokenID); err != nil {
		log.Printf("Failed to delete token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to delete token",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token deleted successfully",
	})
}

// BulkCreateTokens creates multiple tokens at once
func BulkCreateTokens(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	var req struct {
		Count int `json:"count"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
	}

	if req.Count <= 0 || req.Count > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid count",
			Message: "Count must be between 1 and 100",
		})
	}

	var tokens []models.Token
	types := []string{"url", "api_key", "document", "dns", "email"}

	for i := 0; i < req.Count; i++ {
		tokenType := types[i%len(types)]
		tokenValue, err := services.GenerateTokenValue(tokenType)
		if err != nil {
			continue
		}

		token, err := database.CreateToken(userID, tokenType, tokenValue, "Bulk generated")
		if err != nil {
			continue
		}

		tokens = append(tokens, *token)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    tokens,
		"message": "Tokens created successfully",
	})
}
