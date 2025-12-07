CREATE OR REPLACE FUNCTION storage.notify_outbox_event()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('storage_outbox_channel', NEW.event_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER outbox_insert_trigger
AFTER INSERT ON storage.outbox
FOR EACH ROW
EXECUTE FUNCTION storage.notify_outbox_event();
