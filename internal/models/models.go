package models

import "time"

// User represents a KAVACH user
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Token represents a honeypot token
type Token struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	UserID      string    `json:"user_id"`
	TokenType   string    `json:"token_type"` // 'url', 'api_key', 'document', 'dns', 'email'
	TokenValue  string    `json:"token_value"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	TriggeredAt *time.Time `json:"triggered_at"`
}

// Attacker represents a detected attacker
type Attacker struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Fingerprint string   `json:"fingerprint"`
	RiskLevel  string    `json:"risk_level"` // 'low', 'medium', 'high', 'critical'
	RiskScore  float64   `json:"risk_score"` // 0-100
	DetectionCount int   `json:"detection_count"`
	IsKnownUser bool    `json:"is_known_user"`
	OS         string    `json:"os"`
	Browser    string    `json:"browser"`
	DeviceType string    `json:"device_type"` // 'desktop', 'mobile', 'bot'
	IsBlocked  bool      `json:"is_blocked"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TriggerEvent represents a token trigger event
type TriggerEvent struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	TokenID        string    `json:"token_id"`
	AttackerID     *string   `json:"attacker_id"`
	EventType      string    `json:"event_type"` // 'token_accessed', 'token_used'
	HTTPMethod     string    `json:"http_method"`
	Endpoint       string    `json:"endpoint"`
	RequestPath    string    `json:"request_path"`
	RequestHeaders string    `json:"request_headers"`
	IPAddress      string    `json:"ip_address"`
	UserAgent      string    `json:"user_agent"`
	RequestPayload string    `json:"request_payload"`
	ResponseStatus int       `json:"response_status"`
	Timestamp      time.Time `json:"timestamp"`
}

// AlertConfig represents user alert configuration
type AlertConfig struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	AlertType   string    `json:"alert_type"` // 'webhook', 'email', 'slack'
	Destination string    `json:"destination"`
	IsEnabled   bool      `json:"is_enabled"`
	CreatedAt   time.Time `json:"created_at"`
}

// SentAlert represents a sent alert log
type SentAlert struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	TriggerEventID string    `json:"trigger_event_id"`
	AlertConfigID  string    `json:"alert_config_id"`
	Status         string    `json:"status"` // 'pending', 'sent', 'failed'
	ErrorMessage   *string   `json:"error_message"`
	SentAt         *time.Time `json:"sent_at"`
}
