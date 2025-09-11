-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       TEXT        NOT NULL CHECK (title <> ''),
    start_time  TIMESTAMPTZ NOT NULL,
    duration    INTEGER     NOT NULL CHECK (duration > 0),
    description TEXT,
    user_id     TEXT        NOT NULL CHECK (user_id <> ''),
    reminder    INTEGER     CHECK (reminder IS NULL OR reminder > 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd