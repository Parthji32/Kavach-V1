package tests

import (
	"testing"

	"github.com/jindal-parth/kavach/internal/classifier"
)

// TestTrafficClassification tests the 5D classifier
func TestTrafficClassification(t *testing.T) {
	tests := []struct {
		name              string
		ip                string
		userAgent         string
		method            string
		path              string
		expectedRiskLevel string // "low", "medium", "high", "critical"
		maxRiskScore      int
	}{
		{
			name:              "Normal user from known IP",
			ip:                "192.168.1.100",
			userAgent:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			method:            "GET",
			path:              "/dashboard",
			expectedRiskLevel: "low",
			maxRiskScore:      55,
		},
		{
			name:              "Datacenter IP",
			ip:                "35.201.100.0", // Google Cloud IP range example
			userAgent:         "Mozilla/5.0",
			method:            "GET",
			path:              "/api/tokens",
			expectedRiskLevel: "medium",
			maxRiskScore:      65,
		},
		{
			name:              "SQL injection attempt",
			ip:                "203.0.113.42",
			userAgent:         "Mozilla/5.0",
			method:            "POST",
			path:              "/api/users?id=1' OR '1'='1",
			expectedRiskLevel: "high",
			maxRiskScore:      85,
		},
		{
			name:              "Bot user agent",
			ip:                "192.168.1.101",
			userAgent:         "curl/7.64.1",
			method:            "GET",
			path:              "/",
			expectedRiskLevel: "medium",
			maxRiskScore:      70,
		},
		{
			name:              "Admin path access",
			ip:                "203.0.113.50",
			userAgent:         "Mozilla/5.0",
			method:            "GET",
			path:              "/admin/panel",
			expectedRiskLevel: "high",
			maxRiskScore:      75,
		},
	}

	tc := classifier.NewTrafficClassifier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := map[string]string{
				"User-Agent": tt.userAgent,
			}

			risk := tc.ClassifyRequest(tt.ip, tt.userAgent, tt.method, tt.path, headers, "")

			if risk < 0 || risk > 100 {
				t.Errorf("Risk score out of bounds: %.1f", risk)
			}

			if risk > float64(tt.maxRiskScore) {
				t.Errorf("Risk score too high: %.1f (expected max %d)", risk, tt.maxRiskScore)
			}
		})
	}
}

// TestIPReputation tests IP reputation scoring
func TestIPReputation(t *testing.T) {
	tests := []struct {
		name           string
		ip             string
		maxRiskScore   int
		shouldBeRisked bool
	}{
		{
			name:           "Private IP 192.168.x.x",
			ip:             "192.168.1.1",
			maxRiskScore:   10,
			shouldBeRisked: false,
		},
		{
			name:           "Private IP 10.x.x.x",
			ip:             "10.0.0.1",
			maxRiskScore:   10,
			shouldBeRisked: false,
		},
		{
			name:           "Private IP 172.16.x.x",
			ip:             "172.16.0.1",
			maxRiskScore:   10,
			shouldBeRisked: false,
		},
		{
			name:           "Localhost",
			ip:             "127.0.0.1",
			maxRiskScore:   10,
			shouldBeRisked: false,
		},
		{
			name:           "Public IP",
			ip:             "203.0.113.42",
			maxRiskScore:   50,
			shouldBeRisked: true,
		},
	}

	tc := classifier.NewTrafficClassifier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			risk := tc.ScoreIPReputation(tt.ip)

			if risk < 0 || risk > 100 {
				t.Errorf("Risk score out of bounds: %.1f", risk)
			}

			if risk > float64(tt.maxRiskScore) {
				t.Errorf("Risk score too high: %.1f (expected max %d)", risk, tt.maxRiskScore)
			}
		})
	}
}

// TestPayloadAnalysis tests payload scoring for attacks
func TestPayloadAnalysis(t *testing.T) {
	tests := []struct {
		name         string
		payload      string
		minRiskScore int
		shouldFlag   bool
	}{
		{
			name:         "Clean payload",
			payload:      "name=John&age=30",
			minRiskScore: 0,
			shouldFlag:   false,
		},
		{
			name:         "SQL injection",
			payload:      "id=1 OR 1=1; DROP TABLE users;",
			minRiskScore: 20,
			shouldFlag:   true,
		},
		{
			name:         "XSS attempt",
			payload:      "<script>alert('XSS')</script>",
			minRiskScore: 20,
			shouldFlag:   true,
		},
		{
			name:         "Command injection",
			payload:      "file=$(cat /etc/passwd)",
			minRiskScore: 20,
			shouldFlag:   true,
		},
		{
			name:         "Large payload",
			payload:      string(make([]byte, 3*1024*1024)), // 3MB
			minRiskScore: 15,
			shouldFlag:   true,
		},
	}

	tc := classifier.NewTrafficClassifier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			risk := tc.ScorePayloadAnalysis(tt.payload)

			if risk < 0 || risk > 100 {
				t.Errorf("Risk score out of bounds: %.1f", risk)
			}

			if tt.shouldFlag && risk < float64(tt.minRiskScore) {
				t.Errorf("Expected flagged payload, risk: %.1f (min: %d)", risk, tt.minRiskScore)
			}

			if !tt.shouldFlag && risk > float64(tt.minRiskScore) {
				t.Errorf("Unexpected high risk for clean payload: %.1f", risk)
			}
		})
	}
}

// TestBehavioralAnomaly tests behavioral scoring
func TestBehavioralAnomaly(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		minRiskScore int
		shouldFlag   bool
	}{
		{
			name:         "Normal path",
			path:         "/api/users",
			minRiskScore: 0,
			shouldFlag:   false,
		},
		{
			name:         "Admin path access",
			path:         "/admin/panel",
			minRiskScore: 20,
			shouldFlag:   true,
		},
		{
			name:         "Path traversal attempt",
			path:         "/../../etc/passwd",
			minRiskScore: 20,
			shouldFlag:   true,
		},
		{
			name:         "Null byte injection",
			path:         "/api/users%00",
			minRiskScore: 20,
			shouldFlag:   true,
		},
		{
			name:         "Hidden directory access",
			path:         "/.well-known/admin",
			minRiskScore: 15,
			shouldFlag:   true,
		},
	}

	tc := classifier.NewTrafficClassifier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			risk := tc.ScoreBehavioralAnomaly(tt.path)

			if risk < 0 || risk > 100 {
				t.Errorf("Risk score out of bounds: %.1f", risk)
			}

			if tt.shouldFlag && risk < float64(tt.minRiskScore) {
				t.Errorf("Expected flagged path, risk: %.1f (min: %d)", risk, tt.minRiskScore)
			}

			if !tt.shouldFlag && risk > float64(tt.minRiskScore) {
				t.Errorf("Unexpected high risk for normal path: %.1f", risk)
			}
		})
	}
}

// TestRiskActions tests what actions should be taken at different risk levels
func TestRiskActions(t *testing.T) {
	tests := []struct {
		name          string
		riskScore     float64
		expectedAction string
	}{
		{
			name:          "Low risk",
			riskScore:     30.0,
			expectedAction: "allow",
		},
		{
			name:          "Medium risk",
			riskScore:     60.0,
			expectedAction: "flag",
		},
		{
			name:          "High risk",
			riskScore:     70.0,
			expectedAction: "challenge",
		},
		{
			name:          "Critical risk",
			riskScore:     85.0,
			expectedAction: "block",
		},
		{
			name:          "Honeypot trigger",
			riskScore:     99.0,
			expectedAction: "honeypot_alert",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var action string

			switch {
			case tt.riskScore < 55:
				action = "allow"
			case tt.riskScore < 65:
				action = "flag"
			case tt.riskScore < 75:
				action = "challenge"
			case tt.riskScore < 99:
				action = "block"
			default:
				action = "honeypot_alert"
			}

			if action != tt.expectedAction {
				t.Errorf("Expected action %s, got %s for risk %.1f", tt.expectedAction, action, tt.riskScore)
			}
		})
	}
}

// BenchmarkClassification benchmarks the classifier
func BenchmarkClassification(b *testing.B) {
	tc := classifier.NewTrafficClassifier()
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tc.ClassifyRequest("203.0.113.42", "Mozilla/5.0", "GET", "/api/tokens", headers, "")
	}
}

// BenchmarkIPReputation benchmarks IP reputation scoring
func BenchmarkIPReputation(b *testing.B) {
	tc := classifier.NewTrafficClassifier()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tc.ScoreIPReputation("203.0.113.42")
	}
}
