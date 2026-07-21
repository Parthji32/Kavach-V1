package models

// Auth Requests
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	User      User   `json:"user"`
	ExpiresAt int64  `json:"expires_at"`
	ChallengeRequired bool `json:"challenge_required,omitempty"`
	ChallengeMessage  string `json:"challenge_message,omitempty"`
}

// Token Requests
type CreateTokenRequest struct {
	TokenType   string `json:"token_type" validate:"required,oneof=url api_key document dns email"`
	Description string `json:"description"`
}

type CreateTokenResponse struct {
	Token Token `json:"token"`
}

// Alert Config Requests
type CreateAlertConfigRequest struct {
	AlertType   string `json:"alert_type" validate:"required,oneof=webhook email slack"`
	Destination string `json:"destination" validate:"required"`
}

// Dashboard Stats
type DashboardStats struct {
	TotalTokens      int                    `json:"total_tokens"`
	ActiveTokens     int                    `json:"active_tokens"`
	TotalAttackers   int                    `json:"total_attackers"`
	HighRiskCount    int                    `json:"high_risk_count"`
	EventsLast24h    int                    `json:"events_last_24h"`
	RecentAttackers  []AttackerWithRisk     `json:"recent_attackers"`
	RecentEvents     []TriggerEventResponse `json:"recent_events"`
}

type AttackerWithRisk struct {
	ID        string `json:"id"`
	IPAddress string `json:"ip_address"`
	RiskScore int    `json:"risk_score"`
	LastSeen  string `json:"last_seen"`
	IsBlocked bool   `json:"is_blocked"`
}

type TriggerEventResponse struct {
	ID        string `json:"id"`
	TokenType string `json:"token_type"`
	EventType string `json:"event_type"`
	RiskScore int    `json:"risk_score"`
	Timestamp string `json:"timestamp"`
}

// Error Response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
