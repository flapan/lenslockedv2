package models

import "database/sql"

type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	CreatedAt string
}

type UserService struct {
	DB *sql.DB
}
