CREATE TABLE images (
    id TEXT PRIMARY KEY,
    caption TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    location TEXT
);
