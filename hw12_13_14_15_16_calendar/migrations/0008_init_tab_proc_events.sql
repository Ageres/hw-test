-- +goose Up
-- +goose StatementBegin
CREATE TABLE proc_events (
    id          UUID        PRIMARY KEY,
    user_id     TEXT        NOT NULL CHECK (user_id <> '')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS proc_events;
-- +goose StatementEnd