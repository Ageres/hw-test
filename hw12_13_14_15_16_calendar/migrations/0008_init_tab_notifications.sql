-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications (
    id          UUID        PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;
-- +goose StatementEnd