CREATE TABLE IF NOT EXISTS storage.outbox (
    event_id uuid PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_unpublished 
    ON storage.outbox (created_at) 
    WHERE published_at IS NULL;