CREATE TABLE things (
    id BIGSERIAL NOT NULL,
    title TEXT NOT NULL,
    mime TEXT NOT NULL,
    size BIGSERIAL NOT NULL,
    url TEXT NOT NULL,
    tags JSONB,
    metadata JSONB
);

CREATE UNIQUE INDEX idx_things_id ON things (id);
CREATE INDEX idx_things_mime ON things (mime);
CREATE INDEX idx_things_tags ON things USING gin(tags);
CREATE INDEX idx_things_metadata ON things USING gin(metadata);
