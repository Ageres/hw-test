-- создание таблицы событий и индексов

BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE events (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       TEXT        NOT NULL CHECK (title <> ''),
    start_time  TIMESTAMPTZ NOT NULL,
    duration    INTEGER     NOT NULL CHECK (duration > 0),
    description TEXT,
    user_id     TEXT        NOT NULL CHECK (user_id <> ''),
    reminder    INTEGER     CHECK (reminder IS NULL OR reminder > 0)
);
CREATE INDEX idx_user_id ON events (user_id);
CREATE INDEX idx_start_time_idx ON events (start_time);
COMMIT;