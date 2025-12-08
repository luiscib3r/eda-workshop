-- name: CreateFilePage :exec
INSERT INTO ocr.file_pages (id, file_id, page_image_key, page_number)
VALUES ($1, $2, $3, $4);

-- name: GetFilePagesByFileID :many
SELECT *
FROM ocr.file_pages
WHERE file_id = $1
ORDER BY page_number ASC;

-- name: DeleteFilePagesByFileID :exec
DELETE FROM ocr.file_pages
WHERE file_id = $1;

-- name: UpdateFilePageText :one
UPDATE ocr.file_pages
SET text_content = $2
WHERE id = $1
RETURNING *;