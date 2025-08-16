-- удаление таблицы событий и ее индексов

BEGIN;
DROP TABLE IF EXISTS events;
DROP FUNCTION IF EXISTS immutable_tstzrange;
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS btree_gist;
COMMIT;
