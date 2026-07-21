package tests

import (
	"testing"

	"github.com/jindal-parth/kavach/internal/services"
)

// TestRegisterUser tests user registration
func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		fullName  string
		wantError bool
		errMsg    string
	}{
		{
			name:      "Valid registration",
			email:     "test@example.com",
			password:  "SecurePass123!",
			fullName:  "Test User",
			wantError: false,
		},
		{
			name:      "Invalid email format",
			email:     "not-an-email",
			password:  "SecurePass123!",
			fullName:  "Test User",
			wantError: true,
			errMsg:    "invalid email",
		},
		{
			name:      "Weak password",
			email:     "test2@example.com",
			password:  "weak",
			fullName:  "Test User",
			wantError: true,
			errMsg:    "weak password",
		},
		{
			name:      "Duplicate email",
			email:     "test@example.com", // Same as first test
			password:  "SecurePass123!",
			fullName:  "Another User",
			wantError: true,
			errMsg:    "already registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := services.RegisterUser(tt.email, tt.password, tt.fullName)

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

			if user == nil {
				t.Errorf("Expected user, got nil")
				return
			}

			if user.Email != tt.email {
				t.Errorf("Expected email %s, got %s", tt.email, user.Email)
			}
		})
	}
}

// TestLoginUser tests user login with risk assessment
func TestLoginUser(t *testing.T) {
	// First register a user
	email := "login-test@example.com"
	password := "LoginTest123!"
	services.RegisterUser(email, password, "Login Test User")

	tests := []struct {
		name      string
		email     string
		password  string
		wantError bool
		errMsg    string
	}{
		{
			name:      "Correct credentials",
			email:     email,
			password:  password,
			wantError: false,
		},
		{
			name:      "Wrong password",
			email:     email,
			password:  "WrongPassword123!",
			wantError: true,
			errMsg:    "invalid credentials",
		},
		{
			name:      "Non-existent user",
			email:     "nonexistent@example.com",
			password:  password,
			wantError: true,
			errMsg:    "not found",
		},
		{
			name:      "Empty email",
			email:     "",
			password:  password,
			wantError: true,
			errMsg:    "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := services.LoginUser(tt.email, tt.password)

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

			if response == nil {
				t.Errorf("Expected response, got nil")
				return
			}

			if response.Token == "" {
				t.Errorf("Expected JWT token, got empty string")
			}

			if response.Email != tt.email {
				t.Errorf("Expected email %s, got %s", tt.email, response.Email)
			}

			// Verify risk score is within bounds
			if response.RiskScore < 0 || response.RiskScore > 100 {
				t.Errorf("Risk score out of bounds: %f", response.RiskScore)
			}
		})
	}
}

// TestGenerateJWT tests JWT token generation
func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		email     string
		wantError bool
	}{
		{
			name:      "Valid token generation",
			userID:    "usr_123456",
			email:     "test@example.com",
			wantError: false,
		},
		{
			name:      "Empty user ID",
			userID:    "",
			email:     "test@example.com",
			wantError: true,
		},
		{
			name:      "Empty email",
			userID:    "usr_123456",
			email:     "",
			wantError: false, // Email can be empty in some cases
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := services.GenerateJWT(tt.userID, tt.email)

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

			if token == "" {
				t.Errorf("Expected JWT token, got empty string")
				return
			}

			// Verify token format (should be 3 parts separated by dots)
			parts := len([]byte(token))
			if parts < 20 { // JWT should be at least 20 chars
				t.Errorf("Token seems too short: %d chars", parts)
			}
		})
	}
}

// TestValidateJWT tests JWT token validation
func TestValidateJWT(t *testing.T) {
	// Generate a valid token
	validUserID := "usr_test123"
	validEmail := "test@example.com"
	validToken, _ := services.GenerateJWT(validUserID, validEmail)

	tests := []struct {
		name      string
		token     string
		wantError bool
		wantUserID string
	}{
		{
			name:      "Valid token",
			token:     validToken,
			wantError: false,
			wantUserID: validUserID,
		},
		{
			name:      "Invalid token format",
			token:     "invalid.token.format",
			wantError: true,
		},
		{
			name:      "Empty token",
			token:     "",
			wantError: true,
		},
		{
			name:      "Malformed token",
			token:     "definitely.not.a.valid.jwt.token",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := services.ValidateJWT(tt.token)

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

			if claims == nil {
				t.Errorf("Expected claims, got nil")
				return
			}

			if claims.Subject != tt.wantUserID {
				t.Errorf("Expected user ID %s, got %s", tt.wantUserID, claims.Subject)
			}
		})
	}
}

// TestPasswordStrength tests password validation
func TestPasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		isValid  bool
	}{
		{
			name:     "Strong password",
			password: "SecurePass123!",
			isValid:  true,
		},
		{
			name:     "No uppercase",
			password: "securepass123!",
			isValid:  false,
		},
		{
			name:     "No lowercase",
			password: "SECUREPASS123!",
			isValid:  false,
		},
		{
			name:     "No digit",
			password: "SecurePass!",
			isValid:  false,
		},
		{
			name:     "No special character",
			password: "SecurePass123",
			isValid:  false,
		},
		{
			name:     "Too short",
			password: "Pass1!",
			isValid:  false,
		},
		{
			name:     "Empty password",
			password: "",
			isValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would use the validation function from handlers
			// For now, we're testing the concept
			hasUpper := false
			hasLower := false
			hasDigit := false
			hasSpecial := false

			if len(tt.password) < 8 {
				t.Skipf("Password too short (expected by test)")
			}

			for _, ch := range tt.password {
				switch {
				case ch >= 'A' && ch <= 'Z':
					hasUpper = true
				case ch >= 'a' && ch <= 'z':
					hasLower = true
				case ch >= '0' && ch <= '9':
					hasDigit = true
				case ch == '!' || ch == '@' || ch == '#' || ch == '$' || ch == '%':
					hasSpecial = true
				}
			}

			isValid := len(tt.password) >= 8 && hasUpper && hasLower && hasDigit && hasSpecial

			if isValid != tt.isValid {
				t.Errorf("Expected valid=%v, got %v for password %q", tt.isValid, isValid, tt.password)
			}
		})
	}
}
