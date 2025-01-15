package models

import "database/sql"

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

type SessionsService struct {
	DB *sql.DB
}

func (ss *SessionsService) Create(userID int) (*Session, error) {
	// TODO: Create the session token
	// TODO: Create the SessionService.Create
	return nil, nil
}

func (ss *SessionsService) User(token string) (*User, error) {
	// TODO: Implement SessionService.User
	return nil, nil
}
