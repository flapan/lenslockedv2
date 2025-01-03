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

func (us *UserService) Create(password, email string) (*User, error) {
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
