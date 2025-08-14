-- migrations/00001_init_events.up.sql

CREATE TABLE events (
    id          TEXT        PRIMARY KEY,
    title       TEXT        NOT NULL,
    start_time  TIMESTAMPTZ NOT NULL,  
    duration    BIGINT      NOT NULL,
    description TEXT,
    user_id     TEXT        NOT NULL,
    reminder    BIGINT,                 
);

CREATE INDEX user_id_idx ON events (user_id);
CREATE INDEX start_time_idx ON events (start_time);
CREATE INDEX user_id_start_time_idx ON events (user_id, start_time);
