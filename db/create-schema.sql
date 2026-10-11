CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE educations ( 
    id           INTEGER PRIMARY KEY, 
    user_id    INTEGER, 
    description  TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

