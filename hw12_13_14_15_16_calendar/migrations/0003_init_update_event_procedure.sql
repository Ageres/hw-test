-- создание хранимой процедуры обновления события

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
    OUT error_message TEXT,
    OUT conflict_event_id TEXT,
    OUT conflict_user_id TEXT
) 
LANGUAGE plpgsql
AS $$
DECLARE
    v_conflict_event_id UUID;
    v_actual_user_id TEXT;
BEGIN
    -- Устанавливаем таймаут выполнения (60 секунд)
    SET LOCAL statement_timeout = '60s';
    
    -- Инициализация выходных параметров
    conflict_event_id := '';
    conflict_user_id := '';
    RAISE LOG 'Update event attempt. Event ID: %, User ID: %', p_id, p_user_id;

    -- Начинаем транзакционный блок
    BEGIN
        -- Блокировка пользователя для предотвращения гонки условий
        PERFORM pg_advisory_xact_lock(hashtext(p_user_id));
        
        -- проверка существования события с блокировкой
        IF NOT EXISTS (SELECT 1 FROM events WHERE id = p_id FOR UPDATE) THEN
            RAISE NOTICE 'Event not found. ID: %', p_id;
            status_code := 404;
            error_message := 'EVENT_NOT_FOUND';
            RETURN;
        END IF;

        -- проверка владельца
        SELECT user_id INTO v_actual_user_id
        FROM events 
        WHERE id = p_id;
        
        IF v_actual_user_id != p_user_id THEN
            RAISE NOTICE 'Ownership conflict. Event ID: %, Requested User: %, Actual User: %', 
                         p_id, p_user_id, v_actual_user_id;
            status_code := 403;
            error_message := 'OWNERSHIP_CONFLICT';
            conflict_user_id := v_actual_user_id;
            RETURN;
        END IF;

        -- проверка конфликта времени с блокировкой
        SELECT id INTO v_conflict_event_id
        FROM events
        WHERE user_id = p_user_id
          AND id != p_id
          AND immutable_tstzrange(start_time, duration) 
              && immutable_tstzrange(p_start_time, p_duration)
        LIMIT 1
        FOR UPDATE SKIP LOCKED;
        
        IF v_conflict_event_id IS NOT NULL THEN
            RAISE NOTICE 'Time conflict detected for event: %', p_id;
            status_code := 409;
            error_message := 'TIME_CONFLICT';
            conflict_event_id := v_conflict_event_id::TEXT;
            RETURN;
        END IF;

        -- обновление события
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
        WHEN SQLSTATE '57014' THEN -- Код ошибки для statement_timeout
            RAISE EXCEPTION 'Update event timeout for event: %', p_id;
            status_code := 504;
            error_message := 'TIMEOUT: Operation took too long';
        WHEN others THEN
            RAISE EXCEPTION 'Update failed for event %. Error: %', p_id, SQLERRM;
            status_code := 500;
            error_message := 'INTERNAL_ERROR: ' || SQLERRM;
    END;
END;
$$;

COMMENT ON FUNCTION public.update_event IS 'Обновляет событие с проверкой прав доступа и временных конфликтов';

COMMIT;