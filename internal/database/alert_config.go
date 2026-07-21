package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jindal-parth/kavach/internal/models"
)

// CreateAlertConfig creates a new alert configuration
func CreateAlertConfig(userID, alertType, destination string) (*models.AlertConfig, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO alert_configs (id, user_id, alert_type, destination, is_enabled, created_at)
		VALUES (?, ?, ?, ?, 1, ?)
	`

	_, err := DB.Exec(query, id, userID, alertType, destination, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert config: %w", err)
	}

	return &models.AlertConfig{
		ID:          id,
		UserID:      userID,
		AlertType:   alertType,
		Destination: destination,
		IsEnabled:   true,
		CreatedAt:   now,
	}, nil
}

// GetAlertConfigsByUserID retrieves alert configs for a user
func GetAlertConfigsByUserID(userID string) ([]models.AlertConfig, error) {
	query := `
		SELECT id, user_id, alert_type, destination, is_enabled, created_at
		FROM alert_configs WHERE user_id = ? ORDER BY created_at DESC
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query alert configs: %w", err)
	}
	defer rows.Close()

	var configs []models.AlertConfig
	for rows.Next() {
		var config models.AlertConfig
		err := rows.Scan(
			&config.ID, &config.UserID, &config.AlertType, &config.Destination,
			&config.IsEnabled, &config.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert config: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, nil
}

// GetEnabledAlertConfigs retrieves only enabled alert configs for a user
func GetEnabledAlertConfigs(userID string) ([]models.AlertConfig, error) {
	query := `
		SELECT id, user_id, alert_type, destination, is_enabled, created_at
		FROM alert_configs WHERE user_id = ? AND is_enabled = 1 ORDER BY created_at DESC
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query enabled alert configs: %w", err)
	}
	defer rows.Close()

	var configs []models.AlertConfig
	for rows.Next() {
		var config models.AlertConfig
		err := rows.Scan(
			&config.ID, &config.UserID, &config.AlertType, &config.Destination,
			&config.IsEnabled, &config.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert config: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, nil
}

// DeleteAlertConfig deletes an alert configuration
func DeleteAlertConfig(configID string) error {
	query := `DELETE FROM alert_configs WHERE id = ?`
	_, err := DB.Exec(query, configID)
	return err
}
