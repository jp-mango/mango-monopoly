package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	UserName       string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) CheckUsers(tableField string) bool {
	if tableField == "" {
		return false
	}

	query := `
		SELECT id FROM users
		WHERE username = $1 OR email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	err := m.DB.QueryRowContext(ctx, query, tableField).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Println("DB query error:", err)
		return false
	}
	return true
}

func (m *UserModel) Insert(username, email, password string) error {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil
	}

	query := `
		INSERT INTO users (username, email, hashed_password, created)
		VALUES ($1, $2, $3, NOW())`

	args := []any{username, email, string(hashedPW)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPW []byte

	stmt := `
		SELECT id, hashed_password FROM users
		WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// TODO:Figure out error returned here to search for it in handler
	err := m.DB.QueryRowContext(ctx, stmt, email).Scan(&id, &hashedPW)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPW, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
