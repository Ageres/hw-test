-- удаление таблицы событий и индексов
BEGIN;
DROP TABLE  IF EXISTS events;
DROP EXTENSION  IF EXISTS "uuid-ossp";
COMMIT;