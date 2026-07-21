package tests

import (
	"strings"
	"testing"

	"github.com/jindal-parth/kavach/internal/services"
)

// TestGenerateTokenValue tests token generation for all types
func TestGenerateTokenValue(t *testing.T) {
	tests := []struct {
		name      string
		tokenType string
		wantError bool
		prefix    string
		minLength int
	}{
		{
			name:      "Generate URL token",
			tokenType: "url",
			wantError: false,
			prefix:    "sk_",
			minLength: 32,
		},
		{
			name:      "Generate API key token",
			tokenType: "api_key",
			wantError: false,
			prefix:    "sk_",
			minLength: 32,
		},
		{
			name:      "Generate document token",
			tokenType: "document",
			wantError: false,
			prefix:    "sk_",
			minLength: 32,
		},
		{
			name:      "Generate DNS token",
			tokenType: "dns",
			wantError: false,
			prefix:    "sk_",
			minLength: 32,
		},
		{
			name:      "Generate email token",
			tokenType: "email",
			wantError: false,
			prefix:    "sk_",
			minLength: 32,
		},
		{
			name:      "Invalid token type",
			tokenType: "invalid",
			wantError: true,
		},
		{
			name:      "Empty token type",
			tokenType: "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := services.GenerateTokenValue(tt.tokenType)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify token format
			if !strings.HasPrefix(token, tt.prefix) {
				t.Errorf("Expected prefix %s, token: %s", tt.prefix, token)
			}

			if len(token) < tt.minLength {
				t.Errorf("Token too short: %d chars (expected %d)", len(token), tt.minLength)
			}

			// Verify token contains only alphanumeric + underscore
			for _, ch := range token {
				if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_') {
					t.Errorf("Invalid character in token: %c", ch)
				}
			}
		})
	}
}

// TestTokenUniqueness tests that generated tokens are unique
func TestTokenUniqueness(t *testing.T) {
	tokenSet := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		token, err := services.GenerateTokenValue("url")
		if err != nil {
			t.Errorf("Token generation failed: %v", err)
			continue
		}

		if tokenSet[token] {
			t.Errorf("Duplicate token generated: %s", token)
		}
		tokenSet[token] = true
	}

	if len(tokenSet) != iterations {
		t.Errorf("Expected %d unique tokens, got %d", iterations, len(tokenSet))
	}
}

// TestTokenDistribution tests token entropy (should look random)
func TestTokenDistribution(t *testing.T) {
	tokens := make([]string, 100)
	charFreq := make(map[rune]int)

	for i := 0; i < 100; i++ {
		token, err := services.GenerateTokenValue("url")
		if err != nil {
			t.Errorf("Token generation failed: %v", err)
			continue
		}
		tokens[i] = token

		for _, ch := range token {
			charFreq[ch]++
		}
	}

	// Check that we have good distribution (more than 5 different characters)
	if len(charFreq) < 5 {
		t.Errorf("Poor token distribution: only %d unique characters", len(charFreq))
	}

	// Check that no single character dominates (shouldn't appear in >40% of positions)
	totalChars := 0
	for _, count := range charFreq {
		totalChars += count
	}

	for ch, count := range charFreq {
		percentage := float64(count) / float64(totalChars) * 100
		if percentage > 40 {
			t.Errorf("Character %c appears too frequently: %.1f%%", ch, percentage)
		}
	}
}

// TestTokenTypeVariation tests that different token types produce valid tokens
func TestTokenTypeVariation(t *testing.T) {
	types := []string{"url", "api_key", "document", "dns", "email"}

	for _, tokenType := range types {
		token, err := services.GenerateTokenValue(tokenType)
		if err != nil {
			t.Errorf("Failed to generate %s token: %v", tokenType, err)
			continue
		}

		if token == "" {
			t.Errorf("Empty token generated for type: %s", tokenType)
		}

		// All types should have the same format (sk_ prefix)
		if !strings.HasPrefix(token, "sk_") {
			t.Errorf("Invalid prefix for %s token: %s", tokenType, token)
		}
	}
}

// BenchmarkTokenGeneration benchmarks token generation performance
func BenchmarkTokenGeneration(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		services.GenerateTokenValue("url")
	}
}

// BenchmarkTokenGenerationAllTypes benchmarks all token types
func BenchmarkTokenGenerationAllTypes(b *testing.B) {
	types := []string{"url", "api_key", "document", "dns", "email"}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, tokenType := range types {
			services.GenerateTokenValue(tokenType)
		}
	}
}
