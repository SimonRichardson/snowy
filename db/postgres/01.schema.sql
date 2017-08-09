CREATE TABLE IF NOT EXISTS documents (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  resource_id   UUID NOT NULL,
  author_id     TEXT NOT NULL,
  name          TEXT NOT NULL,
  tags          TEXT[] NOT NULL,
  created_on    TIMESTAMPTZ NOT NULL,
  deleted_on    TIMESTAMPTZ NOT NULL
);
