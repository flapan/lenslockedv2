package models

import (
	"database/sql"
	"fmt"

	"github.com/flapan/lenslockedv2/rand"
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
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	token, err := rand.SessionToken()
	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}
	session := Session{
		UserID: userID,
		Token:  token,
		// TODO: set the token hash	TokenHash: ss.hash(token),
	}
	// TODO: Store the session in the DB
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO: Implement SessionService.User
	return nil, nil
}
