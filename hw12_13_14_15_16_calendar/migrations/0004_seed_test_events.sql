-- наполнение таблицы событий тестовыми данными
BEGIN;

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
BEGIN
    -- перебор по юзерам
    FOR i IN 1..user_count LOOP
        user_id := 'user-' || LPAD(i::TEXT, 4, '0'); -- пример "user-0001"
        
        -- перебор по событиям юзеров
        FOR j IN 1..events_per_user LOOP
            -- сборка полей шаблну
            event_title := 'title_' || user_id || '_' || j; -- пример "title_user-0001_1"
            event_desc := event_title || '_desc'; -- пример "title_user-0001_1_desc"
            event_duration := (60 + random() * 172799)::INTEGER; -- период от 1 мин до 2 суток
            event_reminder := (60 + random() * 172799)::INTEGER; 
            event_start := NOW() + (random() * 1095 - 547.5) * INTERVAL '1 day'; -- период ±1.5 года
            
            INSERT INTO events (
                title,
                start_time,
                duration,
                description,
                user_id,
                reminder
            ) VALUES (
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

-- 100 юзеров, по 100 эвентов на каждого
SELECT seed_test_events(100, 100);

DROP FUNCTION seed_test_events;

COMMIT;