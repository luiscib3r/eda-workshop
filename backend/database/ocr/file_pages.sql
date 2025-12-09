-- name: CreateFilePage :exec
INSERT INTO ocr.file_pages (id, file_id, page_image_key, page_number)
VALUES ($1, $2, $3, $4);

-- name: GetFilePagesByFileID :many
SELECT 
    *,
    COUNT(*) OVER() AS total
FROM ocr.file_pages
WHERE file_id = $1
ORDER BY page_number ASC
LIMIT $2 OFFSET $3;

-- name: DeleteFilePagesByFileID :exec
DELETE FROM ocr.file_pages
WHERE file_id = $1;

-- name: UpdateFilePageText :exec
UPDATE ocr.file_pages
SET text_content = $2
WHERE id = $1;