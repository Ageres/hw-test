-- создание хранимой процедуры добавления события

BEGIN;

CREATE EXTENSION IF NOT EXISTS btree_gist;

-- IMMUTABLE-функция для работы с диапазонами
CREATE OR REPLACE FUNCTION immutable_tstzrange(start_time TIMESTAMPTZ, duration INTEGER)
RETURNS TSTZRANGE
LANGUAGE SQL IMMUTABLE
AS $$
    SELECT tstzrange(start_time, start_time + (duration * INTERVAL '1 second'))
$$;

CREATE INDEX IF NOT EXISTS idx_events_id_user ON events (id, user_id);
CREATE INDEX IF NOT EXISTS idx_events_user_time ON events USING gist (
    user_id,
    immutable_tstzrange(start_time, duration)
);

CREATE OR REPLACE FUNCTION public.add_event(
    p_title TEXT,
    p_start_time TIMESTAMPTZ,
    p_duration INTEGER,
    p_description TEXT,
    p_user_id TEXT,
    p_reminder INTEGER,
    OUT status_code INTEGER,
    OUT error_message TEXT,
    OUT event_id UUID
) 
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE LOG 'Add event attempt. User ID: %, Title: %', p_user_id, p_title;

    event_id := '00000000-0000-0000-0000-000000000000'::UUID;

    -- проверка конфликта времени
    IF EXISTS (
        SELECT 1 FROM events
        WHERE user_id = p_user_id
          AND immutable_tstzrange(start_time, duration) 
              && immutable_tstzrange(p_start_time, p_duration)
    ) THEN
        RAISE NOTICE 'Time conflict detected for user: %', p_user_id;
        status_code := 409;
        error_message := 'TIME_CONFLICT';
        RETURN;
    END IF;

    -- вставка события
    INSERT INTO events (
        title, 
        start_time, 
        duration, 
        description, 
        user_id, 
        reminder
    ) VALUES (
        p_title,
        p_start_time,
        p_duration,
        p_description,
        p_user_id,
        p_reminder
    )
    RETURNING id INTO event_id;
    
    RAISE LOG 'Event added successfully. ID: %', event_id;
    status_code := 200;
    error_message := 'SUCCESS';
EXCEPTION
    WHEN others THEN
        RAISE EXCEPTION 'Add failed. Error: %', SQLERRM;
        status_code := 500;
        error_message := 'INTERNAL_ERROR: ' || SQLERRM;
        event_id := NULL;
END;
$$;

COMMENT ON FUNCTION public.add_event IS 'Добавляет новое событие с проверкой временных конфликтов';

COMMIT;
