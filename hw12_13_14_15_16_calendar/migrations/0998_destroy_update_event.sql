-- удаление процедуры обновления события
BEGIN;
DROP FUNCTION IF EXISTS public.update_event;
DROP INDEX IF EXISTS idx_events_id_user;
DROP INDEX IF EXISTS idx_events_user_time;
DROP FUNCTION IF EXISTS immutable_tstzrange;
DROP EXTENSION IF EXISTS btree_gist;
COMMIT;