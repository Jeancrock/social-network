CREATE TABLE IF NOT EXISTS followers(
    id TEXT PRIMARY KEY,
    userId TEXT,       -- l'utilisateur suivi
    followerId TEXT,   -- l'utilisateur qui suit
    FOREIGN KEY(userId) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(followerId) REFERENCES users(id) ON DELETE CASCADE
)