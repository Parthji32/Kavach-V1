package classifier

import (
	"log"
	"math"
	"strings"
	"time"
)

// AdvancedClassifier implements 7-dimensional ML-ready risk scoring
type AdvancedClassifier struct {
	weights map[string]float64
	// Track request patterns per IP for timing analysis
	ipRequestHistory map[string][]time.Time
}

// NewAdvancedClassifier creates a new classifier with ML-ready weights
func NewAdvancedClassifier() *AdvancedClassifier {
	return &AdvancedClassifier{
		// Weights sum to 1.0 for normalization - easily tunable for ML
		weights: map[string]float64{
			"ip_reputation":       0.25,  // Known bad IPs
			"request_rate":        0.15,  // Velocity/DoS patterns
			"payload_analysis":    0.15,  // Injection attacks
			"header_fingerprint":  0.12,  // Bot/automation detection
			"behavioral_anomaly":  0.12,  // Path traversal, admin access
			"geolocation":         0.12,  // VPN/proxy/unusual location
			"timing_pattern":      0.09,  // Machine-like request patterns
		},
		ipRequestHistory: make(map[string][]time.Time),
	}
}

// ClassificationScoreAdvanced represents 7D risk assessment
type ClassificationScoreAdvanced struct {
	OverallRisk        int     // 0-100
	IPReputation       int     // 0-100
	RequestRate        int     // 0-100
	PayloadAnalysis    int     // 0-100
	HeaderFingerprint  int     // 0-100
	BehavioralAnomaly  int     // 0-100
	GeolocationRisk    int     // 0-100
	TimingPattern      int     // 0-100
	IsHoneypotTrigger  bool
	IsSuspicious       bool
	RecommendedAction  string // 'allow', 'flag', 'block', 'challenge'
	TriggeredTokenID   string
	RiskFactors        []string // Human-readable risk reasons
}

// Classify analyzes traffic with 7 dimensions
func (ac *AdvancedClassifier) Classify(
	ipAddress string,
	userAgent string,
	method string,
	path string,
	headers map[string]string,
	body string,
	userID string,
) *ClassificationScoreAdvanced {

	score := &ClassificationScoreAdvanced{}

	// 1. Check honeypot trigger first (immediate block)
	if ac.checkHoneypotTrigger(body, headers) {
		score.IsHoneypotTrigger = true
		score.OverallRisk = 99
		score.RecommendedAction = "block"
		score.RiskFactors = []string{"honeypot_token_detected"}
		return score
	}

	// Calculate each dimension
	d1 := ac.analyzeIPReputation(ipAddress)
	d2 := ac.analyzeRequestRate(ipAddress)
	d3 := ac.analyzePayload(body, method)
	d4 := ac.analyzeHeaders(headers, userAgent)
	d5 := ac.analyzeBehavior(path, method)
	d6 := ac.analyzeGeolocation(ipAddress, headers)
	d7 := ac.analyzeTimingPattern(ipAddress)

	score.IPReputation = d1
	score.RequestRate = d2
	score.PayloadAnalysis = d3
	score.HeaderFingerprint = d4
	score.BehavioralAnomaly = d5
	score.GeolocationRisk = d6
	score.TimingPattern = d7

	// Weighted score calculation
	score.OverallRisk = int(
		math.Round(
			float64(d1)*ac.weights["ip_reputation"] +
				float64(d2)*ac.weights["request_rate"] +
				float64(d3)*ac.weights["payload_analysis"] +
				float64(d4)*ac.weights["header_fingerprint"] +
				float64(d5)*ac.weights["behavioral_anomaly"] +
				float64(d6)*ac.weights["geolocation"] +
				float64(d7)*ac.weights["timing_pattern"],
		),
	)

	// Determine if suspicious (threshold: 55)
	score.IsSuspicious = score.OverallRisk > 55

	// Build risk factors list
	if d1 > 70 {
		score.RiskFactors = append(score.RiskFactors, "high_ip_reputation_risk")
	}
	if d2 > 75 {
		score.RiskFactors = append(score.RiskFactors, "suspicious_request_rate")
	}
	if d3 > 60 {
		score.RiskFactors = append(score.RiskFactors, "payload_injection_detected")
	}
	if d4 > 70 {
		score.RiskFactors = append(score.RiskFactors, "bot_like_headers")
	}
	if d5 > 65 {
		score.RiskFactors = append(score.RiskFactors, "suspicious_behavior")
	}
	if d6 > 60 {
		score.RiskFactors = append(score.RiskFactors, "unusual_geolocation")
	}
	if d7 > 70 {
		score.RiskFactors = append(score.RiskFactors, "machine_like_timing")
	}

	// Recommend action
	if score.OverallRisk > 80 {
		score.RecommendedAction = "block"
	} else if score.OverallRisk > 65 {
		score.RecommendedAction = "challenge" // MFA/CAPTCHA
	} else if score.OverallRisk > 55 {
		score.RecommendedAction = "flag"
	} else {
		score.RecommendedAction = "allow"
	}

	return score
}

// analyzeIPReputation checks IP against threat databases
func (ac *AdvancedClassifier) analyzeIPReputation(ipAddress string) int {
	// TODO: Integrate with:
	// - AbuseIPDB API
	// - MaxMind GeoIP2
	// - AlienVault OTX
	// - Shodan API
	
	knownBadIPs := map[string]int{
		"192.0.2.1":     95,
		"203.0.113.1":   90,
		"198.51.100.1":  85,
	}

	if score, found := knownBadIPs[ipAddress]; found {
		return score
	}

	// Private IPs = lower risk
	if isPrivateIP(ipAddress) {
		return 5
	}

	// Default: moderate risk for unknown public IPs
	return 35
}

// analyzeRequestRate detects DoS/scanning patterns
func (ac *AdvancedClassifier) analyzeRequestRate(ipAddress string) int {
	now := time.Now()
	windowStart := now.Add(-1 * time.Minute)

	// Clean old entries
	if history, exists := ac.ipRequestHistory[ipAddress]; exists {
		var filtered []time.Time
		for _, t := range history {
			if t.After(windowStart) {
				filtered = append(filtered, t)
			}
		}
		ac.ipRequestHistory[ipAddress] = filtered
	}

	// Add current request
	ac.ipRequestHistory[ipAddress] = append(ac.ipRequestHistory[ipAddress], now)

	requestCount := len(ac.ipRequestHistory[ipAddress])

	// Scoring:
	// <5 req/min = normal (0-20)
	// 5-20 req/min = suspicious (20-60)
	// 20-50 req/min = very suspicious (60-85)
	// >50 req/min = DoS pattern (85-100)

	if requestCount < 5 {
		return 5
	} else if requestCount < 20 {
		return int(20 + (requestCount-5)*2)
	} else if requestCount < 50 {
		return int(60 + (requestCount-20))
	} else {
		return 85 + int(math.Min(15, float64(requestCount-50)))
	}
}

// analyzePayload detects injection attacks
func (ac *AdvancedClassifier) analyzePayload(body string, method string) int {
	score := 0

	// SQL injection patterns
	sqlPatterns := []string{"union", "select", "insert", "delete", "drop", "exec", "--", "/*", "*/"}
	for _, pattern := range sqlPatterns {
		if strings.Contains(strings.ToLower(body), pattern) {
			score += 20
		}
	}

	// XSS patterns
	xssPatterns := []string{"<script", "javascript:", "onerror=", "onload=", "onclick="}
	for _, pattern := range xssPatterns {
		if strings.Contains(strings.ToLower(body), pattern) {
			score += 20
		}
	}

	// Command injection
	cmdPatterns := []string{";", "&&", "|", "$(", "`"}
	for _, pattern := range cmdPatterns {
		if strings.Contains(body, pattern) {
			score += 10
		}
	}

	// Large payloads (potential XXE, DoS)
	if len(body) > 10*1024*1024 {
		score += 40
	} else if len(body) > 1*1024*1024 {
		score += 20
	}

	// POST/PUT with suspicious size
	if (method == "POST" || method == "PUT") && len(body) > 100*1024 {
		score += 15
	}

	return int(math.Min(100, float64(score)))
}

// analyzeHeaders detects bots and automation
func (ac *AdvancedClassifier) analyzeHeaders(headers map[string]string, userAgent string) int {
	score := 0

	// Missing common headers
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
			score += 15
		}
	}

	// Suspicious User-Agent
	suspiciousAgents := []string{"bot", "crawler", "scanner", "nikto", "nmap", "masscan", "curl", "wget", "python"}
	for _, agent := range suspiciousAgents {
		if strings.Contains(strings.ToLower(userAgent), agent) {
			score += 25
		}
	}

	// Odd header combinations
	if _, hasXmlHttp := headers["X-Requested-With"]; hasXmlHttp {
		if _, hasReferer := headers["Referer"]; !hasReferer {
			score += 10
		}
	}

	return int(math.Min(100, float64(score)))
}

// analyzeBehavior detects malicious patterns
func (ac *AdvancedClassifier) analyzeBehavior(path string, method string) int {
	score := 0

	// Path traversal
	if strings.Contains(path, "..") || strings.Contains(path, "%2e%2e") {
		score += 50
	}

	// Admin/config endpoints
	adminPaths := []string{"/admin", "/config", "/debug", "/.env", "/backup", "/wp-admin", "/sqlAdmin"}
	for _, adminPath := range adminPaths {
		if strings.HasPrefix(strings.ToLower(path), adminPath) {
			score += 40
		}
	}

	// Unusual HTTP methods
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" && method != "HEAD" {
		score += 30
	}

	// Very long URLs
	if len(path) > 2048 {
		score += 25
	}

	// Null bytes
	if strings.Contains(path, "%00") {
		score += 40
	}

	return int(math.Min(100, float64(score)))
}

// analyzeGeolocation detects VPN/proxy/unusual locations
func (ac *AdvancedClassifier) analyzeGeolocation(ipAddress string, headers map[string]string) int {
	// TODO: Integrate MaxMind GeoIP2 or IP2Location
	
	score := 0

	// VPN/Proxy indicators in headers
	vpnIndicators := []string{"X-Forwarded-For", "X-Proxy-Authorization", "X-Original-IP"}
	for _, header := range vpnIndicators {
		if _, found := headers[header]; found {
			score += 20
		}
	}

	// Private IP = less risk
	if isPrivateIP(ipAddress) {
		return 0
	}

	// Default: moderate geolocation risk
	return 25
}

// analyzeTimingPattern detects machine-like behavior
func (ac *AdvancedClassifier) analyzeTimingPattern(ipAddress string) int {
	history, exists := ac.ipRequestHistory[ipAddress]
	if !exists || len(history) < 5 {
		return 10 // Not enough data
	}

	// Check if requests are evenly spaced (machine-like)
	intervals := []time.Duration{}
	for i := 1; i < len(history); i++ {
		intervals = append(intervals, history[i].Sub(history[i-1]))
	}

	if len(intervals) < 2 {
		return 10
	}

	// Calculate standard deviation of intervals
	mean := time.Duration(0)
	for _, interval := range intervals {
		mean += interval
	}
	mean /= time.Duration(len(intervals))

	variance := float64(0)
	for _, interval := range intervals {
		diff := float64(interval - mean)
		variance += diff * diff
	}
	variance /= float64(len(intervals))

	stdDev := math.Sqrt(variance)

	// Very consistent timing = machine-like
	if stdDev < 100 && mean > 0 { // in milliseconds
		return 75
	} else if stdDev < 500 { // in milliseconds
		return 50
	}

	return 20
}

// checkHoneypotTrigger checks for honeypot tokens
func (ac *AdvancedClassifier) checkHoneypotTrigger(body string, headers map[string]string) bool {
	if strings.Contains(body, "sk_") || strings.Contains(body, "doc_") || strings.Contains(body, "internal") {
		return true
	}

	for key, val := range headers {
		if strings.Contains(strings.ToLower(key), "authorization") && strings.Contains(val, "sk_") {
			return true
		}
	}

	return false
}

// Log classification with factors
func (ac *AdvancedClassifier) LogClassification(score *ClassificationScoreAdvanced) {
	log.Printf(
		"[CLASSIFY-7D] Risk=%d Action=%s | IP=%d Req=%d Payload=%d Header=%d Behavior=%d Geo=%d Timing=%d | Factors: %v",
		score.OverallRisk,
		score.RecommendedAction,
		score.IPReputation,
		score.RequestRate,
		score.PayloadAnalysis,
		score.HeaderFingerprint,
		score.BehavioralAnomaly,
		score.GeolocationRisk,
		score.TimingPattern,
		score.RiskFactors,
	)
}
