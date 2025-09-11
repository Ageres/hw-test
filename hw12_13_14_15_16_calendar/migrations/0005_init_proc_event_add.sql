-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION public.add_event(
    p_title TEXT,
    p_start_time TIMESTAMPTZ,
    p_duration INTEGER,
    p_description TEXT,
    p_user_id TEXT,
    p_reminder INTEGER,
    OUT status_code INTEGER,
    OUT error_message TEXT,
    OUT event_id TEXT,
    OUT conflict_event_id TEXT
) 
LANGUAGE plpgsql
AS $$
DECLARE
    v_conflict_event_id UUID;
    v_event_id UUID;
    timeout_constraint TEXT := 'statement_timeout';
BEGIN
    SET LOCAL statement_timeout = '60s';
    event_id := '';
    conflict_event_id := '';
    RAISE LOG 'Add event attempt. User ID: %, Title: %', p_user_id, p_title;

    BEGIN
        PERFORM pg_advisory_xact_lock(hashtext(p_user_id));
        
        SELECT id INTO v_conflict_event_id
        FROM events
        WHERE user_id = p_user_id
          AND immutable_tstzrange(start_time, duration) 
              && immutable_tstzrange(p_start_time, p_duration)
        LIMIT 1
        FOR UPDATE SKIP LOCKED;
        
        IF v_conflict_event_id IS NOT NULL THEN
            RAISE NOTICE 'Time conflict detected for user: %', p_user_id;
            status_code := 409;
            error_message := 'TIME_CONFLICT';
            conflict_event_id := v_conflict_event_id::TEXT;
            RETURN;
        END IF;

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
        RETURNING id INTO v_event_id;
        
        event_id := v_event_id::TEXT;
        RAISE LOG 'Event added successfully. ID: %', event_id;
        status_code := 200;
        error_message := 'SUCCESS';
        
    EXCEPTION
        WHEN SQLSTATE '57014' THEN
            RAISE EXCEPTION 'Add event timeout for user: %', p_user_id;
            status_code := 504;
            error_message := 'TIMEOUT: Operation took too long';
        WHEN others THEN
            RAISE EXCEPTION 'Add failed. Error: %', SQLERRM;
            status_code := 500;
            error_message := 'INTERNAL_ERROR: ' || SQLERRM;
    END;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS public.add_event;
-- +goose StatementEnd