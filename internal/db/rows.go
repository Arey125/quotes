package db

import "database/sql"

func ForEachRow(rows *sql.Rows, f func(*sql.Rows) error) error {
	for rows.Next() {
		err := f(rows)
		if err != nil {
			return err
		}
	}
	return rows.Err()
}

func Collect[T any](rows *sql.Rows, f func(*sql.Rows, *T) error) ([]T, error) {
	defer rows.Close()
	items := make([]T, 0)
	err := ForEachRow(rows, func(rows *sql.Rows) error {
		var item T
		err := f(rows, &item)
		if err != nil {
			return err
		}

		items = append(items, item)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return items, nil
}
