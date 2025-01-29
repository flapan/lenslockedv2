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
	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		row = ss.DB.QueryRow(`
		INSERT INTO sessions(user_id, token_hash) 
		VALUES($1, $2) 
		RETURNING id;`, session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)
	var user User
	row := ss.DB.QueryRow(`
		SELECT users.id, email, password
		FROM sessions JOIN users ON users.id = sessions.user_id
		WHERE sessions.token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("User: %w", err)
	}
	return &user, nil
}

func (ss SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256(([]byte(token)))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
