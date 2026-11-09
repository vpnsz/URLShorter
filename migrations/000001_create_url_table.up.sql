CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    short_url VARCHAR(2555) NOT NULL
);

CREATE INDEX url_index on urls(url);
