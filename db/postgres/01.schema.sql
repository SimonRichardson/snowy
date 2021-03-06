CREATE TABLE IF NOT EXISTS ledgers (
  id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  parent_id               UUID NOT NULL,
  resource_id             UUID NOT NULL,
  resource_address        TEXT NOT NULL,
  resource_size           BIGINT NOT NULL,
  resource_content_type   TEXT NOT NULL,
  author_id               TEXT NOT NULL,
  name                    TEXT NOT NULL,
  tags                    TEXT[] NOT NULL,
  created_on              TIMESTAMPTZ NOT NULL,
  deleted_on              TIMESTAMPTZ NOT NULL
);
