DROP TABLE IF EXISTS ocr.outbox;
DROP TRIGGER IF EXISTS outbox_insert_trigger ON ocr.outbox;
DROP FUNCTION IF EXISTS ocr.notify_outbox_event();