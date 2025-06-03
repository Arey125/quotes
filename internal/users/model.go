package users

import (
	"database/sql"
	"encoding/gob"
	"errors"
	"quotes/internal/db"

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

type UserPermissions struct {
	CanReadQuotes        bool
	CanWriteQuotes       bool
	CanChangePermissions bool
}

func (p UserPermissions) HasPermission(perm Permisson) bool {
	if perm == PermissonQuotesRead {
		return p.CanReadQuotes
	}
	if perm == PermissonQuotesWrite {
		return p.CanWriteQuotes
	}
	if perm == PermissonUserPermissions {
		return p.CanChangePermissions
	}
	panic("permission does not exist")
}

type UserWithPermissions struct {
	User        User
	Permissions UserPermissions
}

type Model struct {
	db *sql.DB
}

func NewModel(db *sql.DB) Model {
	return Model{db}
}

func (m *Model) Get(id int) (*User, error) {
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

func (m *Model) All() ([]User, error) {
	rows, err := sq.Select("id", "google_user_id", "name").
		From("users").
		RunWith(m.db).
		Query()

	if err != nil {
		return nil, err
	}

	return db.Collect(rows, func(rows *sql.Rows, u *User) error {
		return rows.Scan(&u.Id, &u.GoogleUserId, &u.Name)
	})
}

func (m *Model) GetByGoogleUserId(googleUserId string) (*User, error) {
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

func (m *Model) Add(user User) error {
	_, err := sq.Insert("users").
		Columns("google_user_id", "name").
		Values(user.GoogleUserId, user.Name).
		RunWith(m.db).
		Exec()

	return err
}

func (m *Model) AddPermission(userId int, perm Permisson) error {
	getInsertedRow := sq.Select().
		Column("?", userId).
		Column("id").
		From("permissions").
		Where(sq.Eq{"slug": perm}).
		Limit(1)

	_, err := sq.Insert("user_permissions").
		Options("OR IGNORE").
		Columns("user_id", "permission_id").
		Select(getInsertedRow).
		Values(userId, perm).
		RunWith(m.db).
		Exec()

	return err
}

func (m *Model) RemovePermission(userId int, perm Permisson) error {
	_, err := sq.Delete("user_permissions").
		Where(`user_id = ? and permission_id = 
		(select id from permissions where slug = ? limit 1)`, userId, perm).
		RunWith(m.db).
		Exec()

	return err
}

func (m *Model) HasPermission(userId int, perm Permisson) (bool, error) {
	row := sq.Select("COUNT(1)").
		From("user_permissions").
		Join("permissions on user_permissions.permission_id = permissions.id").
		Where(sq.Eq{"user_id": userId, "slug": perm}).
		RunWith(m.db).
		QueryRow()

	res := 0
	err := row.Scan(&res)
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func (m *Model) GetUserWithPermissions(userId int) (*UserWithPermissions, error) {
	user, err := m.Get(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	canReadQuotes, err := m.HasPermission(user.Id, PermissonQuotesRead)
	if err != nil {
		return nil, err
	}
	canWriteQuotes, err := m.HasPermission(user.Id, PermissonQuotesWrite)
	if err != nil {
		return nil, err
	}
	canChangePermissions, err := m.HasPermission(user.Id, PermissonUserPermissions)
	if err != nil {
		return nil, err
	}

	return &UserWithPermissions{
		User: *user,
		Permissions: UserPermissions{
			CanReadQuotes:        canReadQuotes,
			CanWriteQuotes:       canWriteQuotes,
			CanChangePermissions: canChangePermissions,
		},
	}, nil
}
