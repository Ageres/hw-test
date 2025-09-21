-- +goose Up
-- +goose StatementBegin
CREATE TABLE proc_events (
    id          UUID        PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS proc_events;
-- +goose StatementEnd