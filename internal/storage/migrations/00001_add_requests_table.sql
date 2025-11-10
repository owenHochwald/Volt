-- +goose Up
-- +goose StatementBegin
-- SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS requests (
    id INTEGER PRIMARY KEY,
    name TEXT,
    method TEXT NOT NULL,
    url TEXT NOT NULL,
    headers TEXT,
    body TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS requests;
-- +goose StatementEnd
