package database

import (
	"database/sql"
	"fmt"

	"github.com/jindal-parth/kavach/internal/models"
)

// GetAttackerByFingerprint retrieves attacker by fingerprint
func GetAttackerByFingerprint(fingerprint string) (*models.Attacker, error) {
	var attacker models.Attacker

	query := `
		SELECT id, user_id, fingerprint, ip_address, user_agent, risk_score, risk_level,
		       is_known_user, is_blocked, first_seen, last_seen, detection_count, created_at, updated_at
		FROM attackers WHERE fingerprint = ?
	`

	err := DB.QueryRow(query, fingerprint).Scan(
		&attacker.ID, &attacker.UserID, &attacker.Fingerprint, &attacker.IPAddress,
		&attacker.UserAgent, &attacker.RiskScore, &attacker.RiskLevel,
		&attacker.IsKnownUser, &attacker.IsBlocked, &attacker.FirstSeen, &attacker.LastSeen,
		&attacker.DetectionCount, &attacker.CreatedAt, &attacker.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get attacker by fingerprint: %w", err)
	}

	return &attacker, nil
}

