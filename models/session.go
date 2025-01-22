package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/flapan/lenslockedv2/rand"
)

const (
	// The minimum number of bytes required for a secure token
	// This is used to generate a random token for the session
	MinBytesPerToken = 32
)

// Token is only set when creating a new session,
// when looking up a session this will be left empty
// as we only store the hash of the session token in
// our database and we cannot reverse it into a raw token

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}
	// TODO: Store the session in the DB
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO: Implement SessionService.User
	return nil, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256(([]byte(token)))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
