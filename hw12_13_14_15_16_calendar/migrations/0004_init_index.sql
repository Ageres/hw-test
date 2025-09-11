-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_user_id ON events (user_id);
CREATE INDEX idx_start_time_idx ON events (start_time);
CREATE INDEX idx_user_id_start_time ON events (user_id, start_time);
CREATE INDEX idx_events_id_user ON events (id, user_id);
CREATE INDEX idx_events_user_time ON events USING gist (user_id, immutable_tstzrange(start_time, duration));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_events_user_time;
DROP INDEX idx_events_id_user;
DROP INDEX idx_user_id_start_time;
DROP INDEX idx_start_time_idx;
DROP INDEX idx_user_id;
-- +goose StatementEnd