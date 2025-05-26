CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BLOB NOT NULL,
	expiry REAL NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    google_user_id TEXT,
    name TEXT
);

CREATE INDEX users_google_user_id_idx ON users(google_user_id);
