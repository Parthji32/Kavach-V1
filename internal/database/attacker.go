package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jindal-parth/kavach/internal/models"
)

// CreateOrUpdateAttacker creates or updates an attacker record
func CreateOrUpdateAttacker(userID, fingerprint, ipAddress, userAgent, os, browser, deviceType string, riskScore int) (*models.Attacker, error) {
	id := uuid.New().String()
	now := time.Now()

	// Try to find existing attacker by fingerprint
	query := `SELECT id FROM attackers WHERE fingerprint = ? AND user_id = ?`
	var existingID string
	err := DB.QueryRow(query, fingerprint, userID).Scan(&existingID)

	if err == nil {
		// Update existing attacker
		updateQuery := `
			UPDATE attackers 
			SET last_seen = ?, risk_score = MAX(risk_score, ?), ip_address = ?, user_agent = ?, os = ?, browser = ?, device_type = ?
			WHERE id = ?
		`
		_, err := DB.Exec(updateQuery, now, riskScore, ipAddress, userAgent, os, browser, deviceType, existingID)
		if err != nil {
			return nil, fmt.Errorf("failed to update attacker: %w", err)
		}

		return &models.Attacker{
			ID:         existingID,
			UserID:     userID,
			IPAddress:  ipAddress,
			UserAgent:  userAgent,
			OS:         os,
			Browser:    browser,
			DeviceType: deviceType,
			Fingerprint: fingerprint,
			RiskScore:  float64(riskScore),
			LastSeen:   now,
		}, nil
	}

	// Create new attacker
	insertQuery := `
		INSERT INTO attackers (id, user_id, fingerprint, ip_address, user_agent, os, browser, device_type, risk_score, first_seen, last_seen)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = DB.Exec(insertQuery, id, userID, fingerprint, ipAddress, userAgent, os, browser, deviceType, riskScore, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create attacker: %w", err)
	}

	return &models.Attacker{
		ID:         id,
		UserID:     userID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		OS:         os,
		Browser:    browser,
		DeviceType: deviceType,
		Fingerprint: fingerprint,
		RiskScore:  float64(riskScore),
		FirstSeen:  now,
		LastSeen:   now,
	}, nil
}

// UpdateAttacker updates an existing attacker record
func UpdateAttacker(attacker *models.Attacker) error {
	query := `
		UPDATE attackers 
		SET user_id = ?, fingerprint = ?, ip_address = ?, user_agent = ?, 
		    device_type = ?, risk_score = ?, risk_level = ?, is_known_user = ?, 
		    is_blocked = ?, last_seen = ?, detection_count = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := DB.Exec(
		query,
		attacker.UserID, attacker.Fingerprint, attacker.IPAddress, attacker.UserAgent,
		attacker.DeviceType, attacker.RiskScore, attacker.RiskLevel, attacker.IsKnownUser,
		attacker.IsBlocked, time.Now(), attacker.DetectionCount, time.Now(),
		attacker.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update attacker: %w", err)
	}
	return nil
}

// GetAttackersByUserID retrieves all attackers for a user
func GetAttackersByUserID(userID string, limit int) ([]models.Attacker, error) {
	query := `
		SELECT id, user_id, ip_address, user_agent, os, browser, device_type, fingerprint, risk_score, is_blocked, first_seen, last_seen
		FROM attackers WHERE user_id = ? ORDER BY last_seen DESC LIMIT ?
	`

	rows, err := DB.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query attackers: %w", err)
	}
	defer rows.Close()

	var attackers []models.Attacker
	for rows.Next() {
		var attacker models.Attacker
		err := rows.Scan(
			&attacker.ID, &attacker.UserID, &attacker.IPAddress, &attacker.UserAgent,
			&attacker.OS, &attacker.Browser, &attacker.DeviceType, &attacker.Fingerprint,
			&attacker.RiskScore, &attacker.IsBlocked, &attacker.FirstSeen, &attacker.LastSeen,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attacker: %w", err)
		}
		attackers = append(attackers, attacker)
	}

	return attackers, nil
}

// GetHighRiskAttackers retrieves attackers with risk score >= threshold
func GetHighRiskAttackers(userID string, threshold int) ([]models.Attacker, error) {
	query := `
		SELECT id, user_id, ip_address, user_agent, os, browser, device_type, fingerprint, risk_score, is_blocked, first_seen, last_seen
		FROM attackers WHERE user_id = ? AND risk_score >= ? ORDER BY risk_score DESC
	`

	rows, err := DB.Query(query, userID, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to query high-risk attackers: %w", err)
	}
	defer rows.Close()

	var attackers []models.Attacker
	for rows.Next() {
		var attacker models.Attacker
		err := rows.Scan(
			&attacker.ID, &attacker.UserID, &attacker.IPAddress, &attacker.UserAgent,
			&attacker.OS, &attacker.Browser, &attacker.DeviceType, &attacker.Fingerprint,
			&attacker.RiskScore, &attacker.IsBlocked, &attacker.FirstSeen, &attacker.LastSeen,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attacker: %w", err)
		}
		attackers = append(attackers, attacker)
	}

	return attackers, nil
}

// GetAttackerCount returns total attacker count for user
func GetAttackerCount(userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM attackers WHERE user_id = ?`
	err := DB.QueryRow(query, userID).Scan(&count)
	return count, err
}

// BlockAttacker marks an attacker as blocked
func BlockAttacker(attackerID string) error {
	query := `UPDATE attackers SET is_blocked = 1 WHERE id = ?`
	_, err := DB.Exec(query, attackerID)
	return err
}
