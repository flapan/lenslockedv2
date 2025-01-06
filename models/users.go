package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	Password  string
	Email     string
	CreatedAt string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashshedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		return nil, fmt.Errorf("Create user: %W", err)
	}
	passwordHash := string(hashshedBytes)

	user := User{
		Email:    email,
		Password: passwordHash,
	}

	row := us.DB.QueryRow(`
	INSERT INTO users (password, email)
	VALUES ($1, $2) RETURNING id`, passwordHash, email)

	err = row.Scan(&user.ID)

	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}
	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}

	row := us.DB.QueryRow(`
	SELECT id, password
	FROM users
	WHERE email = $1`, email)

	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("Authenticate: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("Authenticate: %w", err)
	}

	return &user, nil
}
