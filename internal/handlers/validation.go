package handlers

import (
	"regexp"
	"strings"
)

// ValidateTokenType validates that token type is allowed
func ValidateTokenType(tokenType string) bool {
	allowed := map[string]bool{
		"url":      true,
		"api_key":  true,
		"document": true,
		"dns":      true,
		"email":    true,
	}
	return allowed[tokenType]
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ValidatePassword validates password strength
// Requirements:
// - At least 8 characters
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one digit
// - At least one special character
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, ch := range password {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:,.<>?", ch):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// ValidateDescription validates token description
func ValidateDescription(description string) bool {
	if len(description) < 1 {
		return false // Required
	}
	if len(description) > 500 {
		return false // Too long
	}
	return true
}

// ValidateURL validates webhook URL
func ValidateURL(url string) bool {
	if len(url) < 10 || len(url) > 2000 {
		return false
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	return true
}

// ValidateSlackWebhook validates Slack webhook URL format
func ValidateSlackWebhook(url string) bool {
	if !strings.Contains(url, "hooks.slack.com") {
		return false
	}
	return ValidateURL(url)
}

// ValidatePaginationParams validates and normalizes pagination parameters
func ValidatePaginationParams(limit, offset int) (int, int) {
	// Normalize limit
	if limit < 1 {
		limit = 10
	}
	if limit > 500 {
		limit = 500
	}

	// Normalize offset
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}
