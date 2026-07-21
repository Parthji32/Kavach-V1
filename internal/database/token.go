package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jindal-parth/kavach/internal/models"
)

// CreateToken creates a new honeypot token
func CreateToken(userID, tokenType, tokenValue, description string) (*models.Token, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO tokens (id, user_id, token_type, token_value, description, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, 1, ?)
	`

	_, err := DB.Exec(query, id, userID, tokenType, tokenValue, description, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &models.Token{
		ID:          id,
		UserID:      userID,
		TokenType:   tokenType,
		TokenValue:  tokenValue,
		Description: description,
		IsActive:    true,
		CreatedAt:   now,
	}, nil
}

// GetTokensByUserID retrieves all tokens for a user
func GetTokensByUserID(userID string, limit, offset int) ([]models.Token, error) {
	query := `
		SELECT id, user_id, token_type, token_value, description, is_active, created_at, triggered_at
		FROM tokens WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query tokens: %w", err)
	}
	defer rows.Close()

	var tokens []models.Token
	for rows.Next() {
		var token models.Token
		err := rows.Scan(
			&token.ID, &token.UserID, &token.TokenType, &token.TokenValue,
			&token.Description, &token.IsActive, &token.CreatedAt, &token.TriggeredAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan token: %w", err)
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// GetTokenByValue retrieves a token by its value
func GetTokenByValue(tokenValue string) (*models.Token, error) {
	var token models.Token

	query := `
		SELECT id, user_id, token_type, token_value, description, is_active, created_at, triggered_at
		FROM tokens WHERE token_value = ?
	`

	err := DB.QueryRow(query, tokenValue).Scan(
		&token.ID, &token.UserID, &token.TokenType, &token.TokenValue,
		&token.Description, &token.IsActive, &token.CreatedAt, &token.TriggeredAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return &token, nil
}

// GetTokenByID retrieves a token by its ID
func GetTokenByID(tokenID string) (*models.Token, error) {
	var token models.Token

	query := `
		SELECT id, user_id, token_type, token_value, description, is_active, created_at, triggered_at
		FROM tokens WHERE id = ?
	`

	err := DB.QueryRow(query, tokenID).Scan(
		&token.ID, &token.UserID, &token.TokenType, &token.TokenValue,
		&token.Description, &token.IsActive, &token.CreatedAt, &token.TriggeredAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token by ID: %w", err)
	}

	return &token, nil
}

// DeactivateToken deactivates a token
func DeactivateToken(tokenID string) error {
	query := `UPDATE tokens SET is_active = 0 WHERE id = ?`
	_, err := DB.Exec(query, tokenID)
	return err
}

// MarkTokenTriggered marks a token as triggered
func MarkTokenTriggered(tokenID string) error {
	query := `UPDATE tokens SET triggered_at = ? WHERE id = ?`
	_, err := DB.Exec(query, time.Now(), tokenID)
	return err
}

// GetActiveTokensCount returns count of active tokens for user
func GetActiveTokensCount(userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM tokens WHERE user_id = ? AND is_active = 1`
	err := DB.QueryRow(query, userID).Scan(&count)
	return count, err
}
