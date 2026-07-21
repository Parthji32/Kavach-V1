package database

import (
	"database/sql"
	"fmt"
)

// GetDashboardStats retrieves dashboard statistics for a user
func GetDashboardStats(userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total tokens
	var totalTokens int
	row := DB.QueryRow("SELECT COUNT(*) FROM tokens WHERE user_id = ?", userID)
	if err := row.Scan(&totalTokens); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get total tokens: %w", err)
	}
	stats["total_tokens"] = totalTokens

	// Active tokens
	var activeTokens int
	row = DB.QueryRow("SELECT COUNT(*) FROM tokens WHERE user_id = ? AND is_active = 1", userID)
	if err := row.Scan(&activeTokens); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get active tokens: %w", err)
	}
	stats["active_tokens"] = activeTokens

	// Total attackers
	var totalAttackers int
	row = DB.QueryRow("SELECT COUNT(*) FROM attackers WHERE user_id = ?", userID)
	if err := row.Scan(&totalAttackers); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get total attackers: %w", err)
	}
	stats["total_attackers"] = totalAttackers

	// High-risk attackers
	var highRiskCount int
	row = DB.QueryRow("SELECT COUNT(*) FROM attackers WHERE user_id = ? AND risk_score > 75", userID)
	if err := row.Scan(&highRiskCount); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get high-risk count: %w", err)
	}
	stats["high_risk_count"] = highRiskCount

	// Events in last 24 hours
	var eventsLast24h int
	row = DB.QueryRow(`
		SELECT COUNT(*) FROM trigger_events 
		WHERE user_id = ? 
		AND timestamp > datetime('now', '-24 hours')
	`, userID)
	if err := row.Scan(&eventsLast24h); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get events last 24h: %w", err)
	}
	stats["events_last_24h"] = eventsLast24h

	return stats, nil
}
