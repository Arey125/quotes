package users

import (
	"database/sql"
	"encoding/gob"
	"errors"

	sq "github.com/Masterminds/squirrel"
)

func init() {
	gob.Register(User{})
}

var ErrUserNotFound = errors.New("user not found")

type User struct {
	Id           int
	GoogleUserId string
	Name         string
}

type UsersModel struct {
	db *sql.DB
}

func NewModel(db *sql.DB) UsersModel {
	return UsersModel{db}
}

func (m *UsersModel) Get(id int) (*User, error) {
	row := sq.Select("id", "google_user_id", "name").
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1).
		RunWith(m.db).
		QueryRow()

	user := User{}
	err := row.Scan(&user.Id, &user.GoogleUserId, &user.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UsersModel) GetByGoogleUserId(googleUserId string) (*User, error) {
	row := sq.Select("id", "google_user_id", "name").
		From("users").
		Where(sq.Eq{"google_user_id": googleUserId}).
		Limit(1).
		RunWith(m.db).
		QueryRow()

	user := User{}
	err := row.Scan(&user.Id, &user.GoogleUserId, &user.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UsersModel) Add(user User) error {
	_, err := sq.Insert("users").
		Columns("google_user_id", "name").
		Values(user.GoogleUserId, user.Name).
		RunWith(m.db).
		Exec()

	return err
}
