// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
)

type RefreshToken struct {
	Token     string
	CreatedAt string
	UpdatedAt string
	UserID    string
	ExpiresAt string
	RevokedAt sql.NullString
}

type User struct {
	ID             string
	CreatedAt      string
	UpdatedAt      string
	Email          string
	HashedPassword string
}
