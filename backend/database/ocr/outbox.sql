-- name: CreateOutboxEvent :exec
INSERT INTO ocr.outbox (event_id, event_type, payload)
VALUES ($1, $2, $3);

-- name: GetOutboxUnpublishedEvents :many
SELECT event_id, event_type, payload, created_at
FROM ocr.outbox
WHERE published_at IS NULL
ORDER BY created_at
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: MarkEventAsPublished :exec
UPDATE ocr.outbox
SET published_at = NOW()
WHERE event_id = $1;