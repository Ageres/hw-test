-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE events (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       TEXT        NOT NULL CHECK (title <> ''),
    start_time  TIMESTAMPTZ NOT NULL,
    duration    INTEGER     NOT NULL CHECK (duration > 0),
    description TEXT,
    user_id     TEXT        NOT NULL CHECK (user_id <> ''),
    reminder    INTEGER     CHECK (reminder IS NULL OR reminder > 0)
);

CREATE OR REPLACE FUNCTION immutable_tstzrange(start_time TIMESTAMPTZ, duration INTEGER)
RETURNS TSTZRANGE
LANGUAGE SQL IMMUTABLE
AS $$
    SELECT tstzrange(start_time, start_time + (duration * INTERVAL '1 second'))
$$;

CREATE INDEX idx_user_id ON events (user_id);
CREATE INDEX idx_start_time_idx ON events (start_time);
CREATE INDEX idx_user_id_start_time ON events (user_id, start_time);
CREATE INDEX idx_events_time_range ON events USING gist 
    (user_id, immutable_tstzrange(start_time, duration));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
DROP FUNCTION IF EXISTS immutable_tstzrange;
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS btree_gist;
-- +goose StatementEnd