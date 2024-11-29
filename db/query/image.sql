-- name: CreateImageConversion :one
INSERT INTO image_conversions (
    user_id,
    image_name,
    extracted_text
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetImageConversionsByUser :many
SELECT * FROM image_conversions
WHERE user_id = $1;

-- name: GetImageConversionByID :one
SELECT 
    conversion_id,
    user_id,
    image_name,
    extracted_text,
    created_at,
    updated_at
FROM image_conversions
WHERE conversion_id = $1
LIMIT 1;

-- name: UpdateImageConversion :one
UPDATE image_conversions
SET
    image_name = COALESCE($2, image_name),
    extracted_text = COALESCE($3, extracted_text),
    updated_at = NOW()
WHERE conversion_id = $1
RETURNING *;

-- name: DeleteImageConversion :exec
DELETE FROM image_conversions
WHERE conversion_id = $1;
