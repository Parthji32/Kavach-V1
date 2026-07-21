package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jindal-parth/kavach/internal/models"
)

// CreateTriggerEvent creates a new trigger event record
func CreateTriggerEvent(userID, tokenID string, attackerID *string, eventType, httpMethod, endpoint, requestPayload string, responseStatus int) (*models.TriggerEvent, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO trigger_events (id, user_id, token_id, attacker_id, event_type, http_method, endpoint, request_payload, response_status, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := DB.Exec(query, id, userID, tokenID, attackerID, eventType, httpMethod, endpoint, requestPayload, responseStatus, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create trigger event: %w", err)
	}

	return &models.TriggerEvent{
		ID:             id,
		UserID:         userID,
		TokenID:        tokenID,
		AttackerID:     attackerID,
		EventType:      eventType,
		HTTPMethod:     httpMethod,
		Endpoint:       endpoint,
		RequestPayload: requestPayload,
		ResponseStatus: responseStatus,
		Timestamp:      now,
	}, nil
}

// GetTriggerEventsByUserID retrieves trigger events for a user
func GetTriggerEventsByUserID(userID string, limit int) ([]models.TriggerEvent, error) {
	query := `
		SELECT id, user_id, token_id, attacker_id, event_type, http_method, endpoint, request_payload, response_status, timestamp
		FROM trigger_events WHERE user_id = ? ORDER BY timestamp DESC LIMIT ?
	`

	rows, err := DB.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query trigger events: %w", err)
	}
	defer rows.Close()

	var events []models.TriggerEvent
	for rows.Next() {
		var event models.TriggerEvent
		err := rows.Scan(
			&event.ID, &event.UserID, &event.TokenID, &event.AttackerID,
			&event.EventType, &event.HTTPMethod, &event.Endpoint, &event.RequestPayload,
			&event.ResponseStatus, &event.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trigger event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

// GetEventsLast24h returns count of events in last 24 hours
func GetEventsLast24h(userID string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM trigger_events 
		WHERE user_id = ? AND timestamp > datetime('now', '-24 hours')
	`
	err := DB.QueryRow(query, userID).Scan(&count)
	return count, err
}
