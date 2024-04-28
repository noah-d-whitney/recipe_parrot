package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrUserNameInvalid = errors.New("user name is invalid, must be in \"firstname lastname\" format")
	AnonymousUser      = &User{}
)

type User struct {
	ID          int64
	FirstName   *string
	LastName    *string
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
	insertUser := `
		INSERT INTO users (phone_number)
			VALUES ($1)
		RETURNING id, phone_number`

	insertNewList := `
		INSERT INTO lists (user_id)
		VALUES ($1)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	user := new(User)
	err = tx.QueryRowContext(ctx, insertUser, phoneNumber).Scan(&user.ID, &user.PhoneNumber)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	_, err = tx.ExecContext(ctx, insertNewList, user.ID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return user, nil
}

func (m *UserModel) AssignName(userID int64, firstName, lastName string) error {
	stmt := `
		UPDATE users
			SET first_name = $1,
				last_name = $2
			WHERE id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, stmt, firstName, lastName, userID)
	if err != nil {
		return err
	}

	return nil
}
