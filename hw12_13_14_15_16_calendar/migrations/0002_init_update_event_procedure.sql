-- создание хранимой процедуры обновления события обновления события

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

CREATE OR REPLACE FUNCTION public.update_event(
    p_id UUID,
    p_title TEXT,
    p_start_time TIMESTAMPTZ,
    p_duration INTEGER,
    p_description TEXT,
    p_user_id TEXT,
    p_reminder INTEGER,
    OUT status_code INTEGER,
    OUT error_message TEXT
) 
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE LOG 'Update event attempt. Event ID: %, User ID: %', p_id, p_user_id;

    -- проверка существования события
    IF NOT EXISTS (SELECT 1 FROM events WHERE id = p_id) THEN
        RAISE NOTICE 'Event not found. ID: %', p_id;
        status_code := 404;
        error_message := 'EVENT_NOT_FOUND';
        RETURN;
    END IF;

    -- проверка владельца
    IF NOT EXISTS (SELECT 1 FROM events WHERE id = p_id AND user_id = p_user_id) THEN
        RAISE NOTICE 'Ownership conflict. Event ID: %, Requested User: %', p_id, p_user_id;
        status_code := 403;
        error_message := 'OWNERSHIP_CONFLICT';
        RETURN;
    END IF;

    -- проверка конфликта времени с использованием IMMUTABLE-функции
    IF EXISTS (
        SELECT 1 FROM events
        WHERE user_id = p_user_id
          AND id != p_id
          AND immutable_tstzrange(start_time, duration) 
              && immutable_tstzrange(p_start_time, p_duration)
    ) THEN
        RAISE NOTICE 'Time conflict detected for event: %', p_id;
        status_code := 409;
        error_message := 'TIME_CONFLICT';
        RETURN;
    END IF;

    -- обновление
    UPDATE events SET
        title = p_title,
        start_time = p_start_time,
        duration = p_duration,
        description = p_description,
        reminder = p_reminder
    WHERE id = p_id;
    
    RAISE LOG 'Event updated successfully. ID: %', p_id;
    status_code := 200;
    error_message := 'SUCCESS';
EXCEPTION
    WHEN others THEN
        RAISE EXCEPTION 'Update failed for event %. Error: %', p_id, SQLERRM;
        status_code := 500;
        error_message := 'INTERNAL_ERROR: ' || SQLERRM;
END;
$$;

COMMENT ON FUNCTION public.update_event IS 'Обновляет событие с проверкой прав доступа и временных конфликтов';

COMMIT;