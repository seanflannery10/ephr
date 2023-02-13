// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package queries

import (
	"time"
)

type Movie struct {
	ID        int64
	CreatedAt time.Time
	Title     string
	Year      int32
	Runtime   int32
	Genres    []string
	Version   int32
}

type Permission struct {
	ID   int64
	Code string
}

type Token struct {
	Hash   []byte
	UserID int64
	Expiry time.Time
	Scope  string
}

type User struct {
	ID           int64
	CreatedAt    time.Time
	Name         string
	Email        string
	PasswordHash []byte
	Activated    bool
	Version      int32
}

type UsersPermission struct {
	UserID       int64
	PermissionID int64
}
