package quotes

import (
	"database/sql"
	"quotes/internal/db"
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

func scanQuote(scanner sq.RowScanner, q *Quote) error {
	createdAtStr := ""
	err := scanner.Scan(&q.Id, &q.Content, &q.CreatedBy, &createdAtStr)
	if err != nil {
		return err
	}
	createdAt, err := time.Parse(sqlite3.SQLiteTimestampFormats[0], createdAtStr)

	if err != nil {
		return err
	}

	q.CreatedAt = createdAt
	return nil
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
	err := scanQuote(row, &quote)
	if err != nil {
		return nil, err 
	}
	return &quote, nil
}

func (m *Model) All() ([]Quote, error) {
	rows, err := selectQuotes().
		RunWith(m.db).
		Query()

	if err != nil {
		return nil, err
	}

	return db.Collect(rows, func(r *sql.Rows, q *Quote) error {
		return scanQuote(r, q)
	})
}

func (m *Model) Add(quote Quote) error {
	_, err := sq.Insert("quotes").
		Columns("content", "created_by", "created_at").
		Values(quote.Content, quote.CreatedBy, time.Now()).
		RunWith(m.db).
		Exec()

	return err
}
