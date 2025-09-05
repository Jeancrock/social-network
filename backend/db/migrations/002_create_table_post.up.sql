CREATE TABLE IF NOT EXISTS posts(
    id TEXT PRIMARY KEY,
    userId TEXT,
    content TEXT,
    created TEXT,
    FOREIGN KEY(userId) REFERENCES users(id)
)