-- name: CreateFile :one
INSERT INTO storage.files (id, file_name, file_size, file_type, bucket_name)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFiles :many
SELECT 
    *,
    COUNT(*) OVER() AS total
FROM storage.files
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountFiles :one
SELECT COUNT(*) FROM storage.files;

-- name: GetFileByID :one
SELECT * FROM storage.files
WHERE id = $1;

-- name: DeleteFile :exec
DELETE FROM storage.files
WHERE id = $1;

-- name: DeleteFilesByIDs :exec
DELETE FROM storage.files
WHERE id = ANY($1::text[]);