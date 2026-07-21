package services

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GenerateTokenValue generates a unique honeypot token based on type
func GenerateTokenValue(tokenType string) (string, error) {
	switch tokenType {
	case "url":
		return generateURLToken()
	case "api_key":
		return generateAPIKeyToken()
	case "document":
		return generateDocumentToken()
	case "dns":
		return generateDNSToken()
	case "email":
		return generateEmailToken()
	default:
		return "", fmt.Errorf("unsupported token type: %s", tokenType)
	}
}

// generateURLToken creates a URL-like token
func generateURLToken() (string, error) {
	id := uuid.New().String()[:8]
	return fmt.Sprintf("https://api.internal.com/webhook?token=%s", id), nil
}

// generateAPIKeyToken creates an API key-like token
func generateAPIKeyToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	token := fmt.Sprintf("sk_%x", randomBytes)
	return token, nil
}

// generateDocumentToken creates a document name token
func generateDocumentToken() (string, error) {
	id := uuid.New().String()[:8]
	return fmt.Sprintf("doc_%s_confidential.pdf", id), nil
}

// generateDNSToken creates a DNS subdomain token
func generateDNSToken() (string, error) {
	id := uuid.New().String()[:8]
	return fmt.Sprintf("%s.internal.local", id), nil
}

// generateEmailToken creates an email-like token
func generateEmailToken() (string, error) {
	id := uuid.New().String()[:8]
	return fmt.Sprintf("%s@internal-alerts.com", id), nil
}

// TokenGenerator service
type TokenGenerator struct{}

// GenerateTokens generates multiple tokens for user
func (tg *TokenGenerator) GenerateTokens(userID string, count int) ([]string, error) {
	var tokens []string

	for i := 0; i < count; i++ {
		// Alternate between token types
		types := []string{"url", "api_key", "document", "dns", "email"}
		tokenType := types[i%len(types)]

		token, err := GenerateTokenValue(tokenType)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

// HashToken creates a simple hash of a token for storage
func HashToken(token string) string {
	return strings.TrimSpace(fmt.Sprintf("%x", len(token)))
}
