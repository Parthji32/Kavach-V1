package services

import (
	"time"
)

// UserContext tracks user behavior to allow known users
type UserContext struct {
	UserID            string
	LastLoginIP       string
	LastLoginTime     time.Time
	TrustedDevices    []string // Device fingerprints
	IsVerifiedUser    bool     // Email verified
	TotalSuccessLogins int
}

// GetOrCreateUserContext retrieves user context
func GetOrCreateUserContext(userID string) (*UserContext, error) {
	// TODO: Store in cache or database
	return &UserContext{
		UserID:         userID,
		IsVerifiedUser: true,
		TotalSuccessLogins: 0,
	}, nil
}

// IsKnownUser checks if user should get frictionless access
func IsKnownUser(userID, ipAddress, deviceFingerprint string) bool {
	ctx, err := GetOrCreateUserContext(userID)
	if err != nil {
		return false
	}

	// Known user if:
	// 1. Same IP as last login
	if ctx.LastLoginIP == ipAddress && time.Since(ctx.LastLoginTime) < 24*time.Hour {
		return true
	}

	// 2. Device in trusted list
	for _, trusted := range ctx.TrustedDevices {
		if trusted == deviceFingerprint {
			return true
		}
	}

	// 3. Many successful logins from diverse IPs (established user)
	if ctx.TotalSuccessLogins > 50 {
		return true
	}

	return false
}

// RecordSuccessfulLogin updates user context after successful login
func RecordSuccessfulLogin(userID, ipAddress, deviceFingerprint string) error {
	// TODO: Update database/cache
	return nil
}

// AddTrustedDevice adds device to user's trusted list
func AddTrustedDevice(userID, deviceFingerprint string) error {
	// TODO: Store in database
	return nil
}

// CheckPasswordStrength validates password meets security requirements
func CheckPasswordStrength(password string) (bool, string) {
	if len(password) < 12 {
		return false, "Password must be at least 12 characters"
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		} else if char == '!' || char == '@' || char == '#' || char == '$' || char == '%' {
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return false, "Password must contain uppercase, lowercase, digit, and special character"
	}

	return true, ""
}
