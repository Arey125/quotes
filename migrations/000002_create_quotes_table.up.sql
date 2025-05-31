CREATE TABLE quotes (
    id  INTEGER PRIMARY KEY,
	content TEXT,
	created_by INTEGER,
	created_at TEXT,

    FOREIGN KEY(created_by) REFERENCES users(id)
);

CREATE INDEX quotes_created_at ON quotes(created_at);
