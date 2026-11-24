DROP INDEX IF EXISTS url_index;
CREATE UNIQUE INDEX url_index on urls(url);