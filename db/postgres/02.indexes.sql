CREATE INDEX documents_tags ON documents USING GIN(tags);
