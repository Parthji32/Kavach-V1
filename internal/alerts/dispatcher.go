package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// AlertDispatcher handles sending alerts to various destinations
type AlertDispatcher struct {
	client *http.Client
}

// NewAlertDispatcher creates a new alert dispatcher
func NewAlertDispatcher() *AlertDispatcher {
	return &AlertDispatcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// WebhookPayload is the JSON payload sent to webhooks
type WebhookPayload struct {
	EventType    string    `json:"event_type"`     // "token_accessed", "attacker_blocked"
	Timestamp    time.Time `json:"timestamp"`
	UserID       string    `json:"user_id"`
	TokenID      string    `json:"token_id"`
	TokenValue   string    `json:"token_value"`   // Redacted (first 20 chars)
	TokenType    string    `json:"token_type"`    // "url", "api_key", etc
	AttackerID   string    `json:"attacker_id"`
	AttackerIP   string    `json:"attacker_ip"`
	RiskScore    int       `json:"risk_score"`    // 0-100
	RiskLevel    string    `json:"risk_level"`    // "low", "medium", "high", "critical"
	DetectedAt   string    `json:"detected_at"`
	Severity     string    `json:"severity"`      // "info", "warning", "critical"
	Message      string    `json:"message"`
}

// SendWebhookAlert sends an alert to a webhook URL
func (ad *AlertDispatcher) SendWebhookAlert(webhookURL string, payload *WebhookPayload) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	log.Printf("[WEBHOOK-PAYLOAD] %s", string(jsonData))
	

	// Retry logic: 3 attempts with exponential backoff
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "KAVACH/1.0.0")
		req.Header.Set("X-KAVACH-Signature", "sha256=test") // TODO: Add HMAC signature

		resp, err := ad.client.Do(req)
		if err != nil {
			lastErr = err
			// Exponential backoff: 1s, 2s, 4s
			time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("[WEBHOOK-SENT] %s | Status: %d | Risk: %s", webhookURL, resp.StatusCode, payload.RiskLevel)
			return nil
		}

		lastErr = fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return fmt.Errorf("webhook delivery failed after 3 attempts: %w", lastErr)
}

// SendSlackAlert sends a formatted alert to Slack webhook
func (ad *AlertDispatcher) SendSlackAlert(slackWebhookURL string, payload *WebhookPayload) error {
	if slackWebhookURL == "" {
		return fmt.Errorf("slack webhook URL is empty")
	}

	// Determine color based on risk level
	color := "#36a64f" // green
	if payload.RiskLevel == "high" {
		color = "#ff9800" // orange
	} else if payload.RiskLevel == "critical" {
		color = "#f44336" // red
	}

	// Slack message format
	slackMsg := map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"color":      color,
				"title":      "🚨 KAVACH Alert: Honeypot Triggered",
				"title_link": "http://localhost:3000/app", // TODO: Make configurable
				"text":       payload.Message,
				"fields": []map[string]interface{}{
					{"title": "Risk Level", "value": payload.RiskLevel, "short": true},
					{"title": "Risk Score", "value": fmt.Sprintf("%d/100", payload.RiskScore), "short": true},
					{"title": "Attacker IP", "value": payload.AttackerIP, "short": true},
					{"title": "Token Type", "value": payload.TokenType, "short": true},
					{"title": "Token Value", "value": payload.TokenValue, "short": false},
					{"title": "Detected At", "value": payload.DetectedAt, "short": false},
				},
				"ts": time.Now().Unix(),
			},
		},
	}

	jsonData, err := json.Marshal(slackMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	req, err := http.NewRequest("POST", slackWebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ad.client.Do(req)
	if err != nil {
		return fmt.Errorf("slack request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack returned status %d", resp.StatusCode)
	}

	log.Printf("[SLACK-SENT] Alert delivered | Risk: %s", payload.RiskLevel)
	return nil
}

// BuildWebhookPayload creates a webhook payload from attacker and event data
func BuildWebhookPayload(userID, tokenID, tokenValue, tokenType, attackerID, attackerIP string, riskScore int) *WebhookPayload {
	riskLevel := "low"
	if riskScore > 75 {
		riskLevel = "critical"
	} else if riskScore > 60 {
		riskLevel = "high"
	} else if riskScore > 40 {
		riskLevel = "medium"
	}

	severity := "info"
	if riskLevel == "critical" || riskLevel == "high" {
		severity = "warning"
	}
	if riskLevel == "critical" {
		severity = "critical"
	}

	// Redact token value (show only first 20 chars)
	redactedToken := tokenValue
	if len(tokenValue) > 20 {
		redactedToken = tokenValue[:20] + "..."
	}

	return &WebhookPayload{
		EventType:  "token_accessed",
		Timestamp:  time.Now(),
		UserID:     userID,
		TokenID:    tokenID,
		TokenValue: redactedToken,
		TokenType:  tokenType,
		AttackerID: attackerID,
		AttackerIP: attackerIP,
		RiskScore:  riskScore,
		RiskLevel:  riskLevel,
		DetectedAt: time.Now().Format(time.RFC3339),
		Severity:   severity,
		Message:    fmt.Sprintf("Honeypot token (%s) accessed from %s with risk score %d (%s)", tokenType, attackerIP, riskScore, riskLevel),
	}
}
