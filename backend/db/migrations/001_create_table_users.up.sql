CREATE TABLE IF NOT EXISTS users(
    id TEXT PRIMARY KEY,
    username TEXT,
    email TEXT UNIQUE,
    password TEXT
)
