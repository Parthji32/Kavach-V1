package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jindal-parth/kavach/internal/models"
)

// CreateUser creates a new user in the database
func CreateUser(email, passwordHash, fullName string) (*models.User, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO users (id, email, password_hash, full_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := DB.Exec(query, id, email, passwordHash, fullName, now, now)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return nil, errors.New("email already registered")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &models.User{
		ID:        id,
		Email:     email,
		FullName:  fullName,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	query := `SELECT id, email, full_name, created_at, updated_at FROM users WHERE email = ?`
	err := DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.FullName, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id string) (*models.User, error) {
	var user models.User

	query := `SELECT id, email, full_name, created_at, updated_at FROM users WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.FullName, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserPassword retrieves user's password hash (for login verification)
func GetUserPassword(email string) (string, error) {
	var passwordHash string

	query := `SELECT password_hash FROM users WHERE email = ?`
	err := DB.QueryRow(query, email).Scan(&passwordHash)

	if err == sql.ErrNoRows {
		return "", errors.New("user not found")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get password: %w", err)
	}

	return passwordHash, nil
}
