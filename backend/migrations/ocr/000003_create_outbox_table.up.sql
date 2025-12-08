CREATE TABLE IF NOT EXISTS ocr.outbox (
    event_id uuid PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_unpublished 
    ON ocr.outbox (created_at) 
    WHERE published_at IS NULL;

CREATE OR REPLACE FUNCTION ocr.notify_outbox_event()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('ocr_outbox_channel', NEW.event_id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER outbox_insert_trigger
AFTER INSERT ON ocr.outbox
FOR EACH ROW
EXECUTE FUNCTION ocr.notify_outbox_event();