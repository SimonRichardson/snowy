CREATE TABLE IF NOT EXISTS documents (
  id            SERIAL PRIMARY KEY,
  resource_id   CHAR(36) NOT NULL,
  author_id     CHAR(36) NOT NULL,
  name          VARCHAR(255) NOT NULL,
  tags          TEXT[] NOT NULL,
  created_on    DATE NOT NULL,
  deleted_on    DATE NOT NULL
);