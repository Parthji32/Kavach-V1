package classifier

import (
	"log"
	"strings"

	"github.com/jindal-parth/kavach/internal/models"
)

// TrafficClassifier scores incoming traffic
type TrafficClassifier struct {
	tokenCache map[string]*models.Token
}

// NewTrafficClassifier creates a new classifier instance
func NewTrafficClassifier() *TrafficClassifier {
	return &TrafficClassifier{
		tokenCache: make(map[string]*models.Token),
	}
}

// ClassificationScore represents the security assessment of traffic
type ClassificationScore struct {
	OverallRisk          int // 0-100
	IPReputation         int // 0-100
	RequestRate          int // 0-100
	PayloadAnalysis      int // 0-100
	HeaderFingerprint    int // 0-100
	BehavioralAnomaly    int // 0-100
	IsHoneypotTrigger    bool
	IsSuspicious         bool
	RecommendedAction    string // 'allow', 'flag', 'block'
	TriggeredTokenID     string
}

// Classify analyzes incoming traffic
func (tc *TrafficClassifier) Classify(
	ipAddress string,
	userAgent string,
	method string,
	path string,
	headers map[string]string,
	body string,
) *ClassificationScore {

	score := &ClassificationScore{}

	// 1. Check if request contains honeypot token
	score.IsHoneypotTrigger = tc.checkHoneypotTrigger(body, headers)
	if score.IsHoneypotTrigger {
		score.OverallRisk = 95
		score.RecommendedAction = "block"
		return score
	}

	// 2. IP Reputation (0-100)
	score.IPReputation = tc.analyzeIPReputation(ipAddress)

	// 3. Request Rate analysis (placeholder - needs rate limiter integration)
	score.RequestRate = 10 // Low risk by default

	// 4. Payload Analysis
	score.PayloadAnalysis = tc.analyzePayload(body, method)

	// 5. Header & Fingerprint Analysis
	score.HeaderFingerprint = tc.analyzeHeaders(headers, userAgent)

	// 6. Behavioral Anomaly
	score.BehavioralAnomaly = tc.analyzeBehavior(path, method)

	// Calculate overall risk
	score.OverallRisk = (
		score.IPReputation*2 +
			score.RequestRate +
			score.PayloadAnalysis +
			score.HeaderFingerprint +
			score.BehavioralAnomaly) / 6

	// Determine if suspicious
	score.IsSuspicious = score.OverallRisk > 50

	// Recommend action
	if score.OverallRisk > 80 {
		score.RecommendedAction = "block"
	} else if score.OverallRisk > 60 {
		score.RecommendedAction = "flag"
	} else {
		score.RecommendedAction = "allow"
	}

	return score
}

// checkHoneypotTrigger checks if request contains a honeypot token
func (tc *TrafficClassifier) checkHoneypotTrigger(body string, headers map[string]string) bool {
	// Check in URL path
	if strings.Contains(body, "sk_") || strings.Contains(body, "token=") {
		return true
	}

	// Check in body
	if strings.Contains(body, "sk_") || strings.Contains(body, "token") {
		return true
	}

	// Check in headers
	for key, val := range headers {
		if strings.ToLower(key) == "authorization" && strings.Contains(val, "sk_") {
			return true
		}
		if strings.Contains(val, "doc_") && strings.Contains(val, "_confidential") {
			return true
		}
	}

	return false
}

// analyzeIPReputation scores an IP address (0-100)
func (tc *TrafficClassifier) analyzeIPReputation(ipAddress string) int {
	// Placeholder: In production, integrate with IP reputation service
	// (e.g., abuse.ch, MaxMind, AlienVault OTX, etc.)

	// Known bad IPs (example)
	badIPs := map[string]bool{
		"192.0.2.1":     true, // Example: known attacker
		"203.0.113.1":   true,
		"198.51.100.1":  true,
	}

	if badIPs[ipAddress] {
		return 95
	}

	// Check if private IP
	if isPrivateIP(ipAddress) {
		return 10
	}

	// Default moderate risk for unknown public IPs
	return 40
}

// isPrivateIP checks if IP is private
func isPrivateIP(ip string) bool {
	privateRanges := []string{
		"127.", "10.", "172.", "192.168.", "::1",
	}
	for _, prefix := range privateRanges {
		if strings.HasPrefix(ip, prefix) {
			return true
		}
	}
	return false
}

// analyzePayload scores request body (0-100)
func (tc *TrafficClassifier) analyzePayload(body string, method string) int {
	score := 0

	// Detect SQL injection patterns
	sqlPatterns := []string{"union", "select", "insert", "delete", "drop", "exec", "execute"}
	for _, pattern := range sqlPatterns {
		if strings.Contains(strings.ToLower(body), pattern) {
			score += 20
		}
	}

	// Detect XSS patterns
	xssPatterns := []string{"<script", "javascript:", "onerror=", "onload="}
	for _, pattern := range xssPatterns {
		if strings.Contains(strings.ToLower(body), pattern) {
			score += 20
		}
	}

	// Detect large payloads (potential DoS)
	if len(body) > 10*1024*1024 { // 10MB
		score += 30
	}

	// POST/PUT with large body is higher risk
	if method == "POST" || method == "PUT" {
		if len(body) > 100*1024 { // 100KB
			score += 15
		}
	}

	if score > 100 {
		score = 100
	}

	return score
}

// analyzeHeaders scores headers & fingerprinting (0-100)
func (tc *TrafficClassifier) analyzeHeaders(headers map[string]string, userAgent string) int {
	score := 0

	// Check for missing headers
	expectedHeaders := []string{"Host", "User-Agent", "Accept"}
	for _, header := range expectedHeaders {
		found := false
		for key := range headers {
			if strings.EqualFold(key, header) {
				found = true
				break
			}
		}
		if !found {
			score += 10
		}
	}

	// Check for suspicious User-Agent
	suspiciousAgents := []string{"bot", "crawler", "scanner", "nikto", "nmap", "masscan"}
	for _, agent := range suspiciousAgents {
		if strings.Contains(strings.ToLower(userAgent), agent) {
			score += 25
		}
	}

	// Check for automation tools
	automationHeaders := []string{"X-Scanner", "X-Forwarded-For", "X-Abuse-Email"}
	for _, header := range automationHeaders {
		if _, found := headers[header]; found {
			score += 15
		}
	}

	if score > 100 {
		score = 100
	}

	return score
}

// analyzeBehavior scores behavioral patterns (0-100)
func (tc *TrafficClassifier) analyzeBehavior(path string, method string) int {
	score := 0

	// Check for path traversal attempts
	if strings.Contains(path, "..") || strings.Contains(path, "%2e%2e") {
		score += 40
	}

	// Check for admin/config endpoints
	adminPaths := []string{"/admin", "/config", "/debug", "/.env", "/backup"}
	for _, adminPath := range adminPaths {
		if strings.HasPrefix(strings.ToLower(path), adminPath) {
			score += 30
		}
	}

	// Unusual HTTP methods
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		score += 20
	}

	// Very long URLs
	if len(path) > 2048 {
		score += 20
	}

	if score > 100 {
		score = 100
	}

	return score
}

// LogClassification logs the classification for debugging
func (tc *TrafficClassifier) LogClassification(score *ClassificationScore) {
	log.Printf(
		"[CLASSIFY] Risk: %d | Action: %s | IP: %d | Payload: %d | Header: %d | Behavior: %d",
		score.OverallRisk,
		score.RecommendedAction,
		score.IPReputation,
		score.PayloadAnalysis,
		score.HeaderFingerprint,
		score.BehavioralAnomaly,
	)
}
