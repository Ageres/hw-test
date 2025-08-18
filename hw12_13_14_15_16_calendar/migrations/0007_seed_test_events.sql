-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION seed_test_events(user_count INTEGER, events_per_user INTEGER) 
RETURNS VOID AS $$
DECLARE
    user_id TEXT;
    event_title TEXT;
    event_desc TEXT;
    event_duration INTEGER;
    event_reminder INTEGER;
    event_start TIMESTAMPTZ;
    i INTEGER;
    j INTEGER;
    event_id UUID;
    event_counter INTEGER := 0;
BEGIN
    FOR i IN 1..user_count LOOP
        user_id := 'user-' || LPAD(i::TEXT, 4, '0');
        
        FOR j IN 1..events_per_user LOOP
            event_title := 'title_' || user_id || '_' || j;
            event_desc := event_title || '_desc';
            event_duration := (60 + random() * 172799)::INTEGER;
            event_reminder := (60 + random() * 172799)::INTEGER; 
            event_start := NOW() + (random() * 1095 - 547.5) * INTERVAL '1 day';
            event_id := ('00000000-0000-0000-0000-' || LPAD(event_counter::TEXT, 12, '0'))::UUID;
            event_counter := event_counter + 1;
            
            INSERT INTO events (
                id,
                title,
                start_time,
                duration,
                description,
                user_id,
                reminder
            ) VALUES (
                event_id,
                event_title,
                event_start,
                event_duration,
                event_desc,
                user_id,
                event_reminder
            );
            
            RAISE NOTICE 'Added event for user %: %', user_id, event_title;
        END LOOP;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT seed_test_events(100, 100);
DROP FUNCTION seed_test_events;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE events;
-- +goose StatementEnd