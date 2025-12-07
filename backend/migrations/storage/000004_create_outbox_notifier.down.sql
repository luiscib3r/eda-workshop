DROP TRIGGER IF EXISTS outbox_insert_trigger ON storage.outbox;
DROP FUNCTION IF EXISTS storage.notify_outbox_event();
