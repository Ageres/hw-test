-- migrations/00001_init_events.sql
BEGIN;

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

CREATE INDEX user_id_idx ON events (user_id);
CREATE INDEX start_time_idx ON events (start_time);
CREATE INDEX user_id_start_time_idx ON events (user_id, start_time);

COMMENT ON CONSTRAINT no_time_overlaps ON events IS 
    'Prevents time overlaps for events of the same user';

ALTER TABLE events ADD CONSTRAINT no_time_overlaps
    EXCLUDE USING gist (
        user_id WITH =,
        tstzrange(start_time, start_time + (duration * INTERVAL '1 second')) WITH &&
    );

COMMIT;