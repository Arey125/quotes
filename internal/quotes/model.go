package quotes

import (
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/mattn/go-sqlite3"
)

type Quote struct {
	Id        int
	Content   string
	CreatedBy int
	CreatedAt time.Time
}

type Model struct {
	db *sql.DB
}

func NewModel(db *sql.DB) Model {
	return Model{db}
}

func selectQuotes() sq.SelectBuilder {
	return sq.Select("id", "content", "created_by", "created_at").From("quotes")
}

func (m *Model) Get(id int) (*Quote, error) {
	row := selectQuotes().
		Where(sq.Eq{"id": id}).
		Limit(1).
		RunWith(m.db).
		QueryRow()

	quote := Quote{}
	createdAtStr := ""
	err := row.Scan(
		&quote.Id,
		&quote.Content,
		&quote.CreatedBy,
		&createdAtStr,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse(sqlite3.SQLiteTimestampFormats[0], createdAtStr)
	if err != nil {
		return nil, err
	}
	quote.CreatedAt = createdAt
	return &quote, nil
}

func (m *Model) All() ([]Quote, error) {
	rows, err := selectQuotes().
		RunWith(m.db).
		Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quotes := make([]Quote, 0)
	for rows.Next() {
		q := Quote{}
		createdAtStr := ""
		err := rows.Scan(
			&q.Id,
			&q.Content,
			&q.CreatedBy,
			&createdAtStr,
		)
		if err != nil {
			return nil, err
		}
		createdAt, err := time.Parse(sqlite3.SQLiteTimestampFormats[0], createdAtStr)
		if err != nil {
			return nil, err
		}
		q.CreatedAt = createdAt

		quotes = append(quotes, q)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return quotes, nil
}

func (m *Model) Add(quote Quote) error {
	_, err := sq.Insert("quotes").
		Columns("content", "created_by", "created_at").
		Values(quote.Content, quote.CreatedBy, time.Now()).
		RunWith(m.db).
		Exec()

	return err
}
