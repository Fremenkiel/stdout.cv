CREATE TABLE users (
    id    INTEGER PRIMARY KEY, 
    name  TEXT
)

CREATE TABLE educations ( 
    id           INTEGER PRIMARY KEY, 
    user_id    INTEGER, 
    description  TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

