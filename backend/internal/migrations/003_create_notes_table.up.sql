CREATE TABLE notes
(
    id         SERIAL PRIMARY KEY,
    user_id    INT  NOT NULL REFERENCES users (id),
    content    TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_user_notes ON notes (user_id);
