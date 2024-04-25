package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	ErrUserNameInvalid = errors.New("user name is invalid, must be in \"firstname lastname\" format")
	AnonymousUser      = &User{}
)

type User struct {
	ID          int64
	FirstName   string
	LastName    string
	PhoneNumber string
}

type UserModel struct {
	db *sql.DB
}

func (m *UserModel) Get(phoneNumber string) (*User, error) {
	stmt := `
		SELECT id, first_name, last_name, phone_number
		FROM users
		WHERE phone_number = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := new(User)
	err := m.db.QueryRowContext(ctx, stmt, phoneNumber).Scan(&user.ID, &user.FirstName, &user.LastName, &user.PhoneNumber)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (m *UserModel) Create(phoneNumber string) (*User, error) {
	stmt := `
		INSERT INTO users (phone_number)
			VALUES ($1)
		RETURNING id, phone_number`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := new(User)
	err := m.db.QueryRowContext(ctx, stmt, phoneNumber).Scan(&user.ID, &user.PhoneNumber)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *UserModel) AssignName(userID int64, fullName string) error {
	stmt := `
		UPDATE users
			SET first_name = $1,
				last_name = $2,
			WHERE id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	names := strings.Split(fullName, " ")
	if len(names) != 2 {
		return ErrUserNameInvalid
	}

	_, err := m.db.ExecContext(ctx, stmt, names[0], names[1], userID)
	if err != nil {
		return err
	}

	return nil
}
