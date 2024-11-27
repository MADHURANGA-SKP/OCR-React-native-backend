-- name: CreateUsers :one
INSERT INTO users (
    user_name,
    first_name,
    last_name,
    email,
    hashed_password
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUsers :one
SELECT * FROM users
WHERE user_name = $1 LIMIT 1;

-- name: GetUserID :one
SELECT  
    u.user_id,
    u.user_name,
    u.first_name,
    u.last_name,
    u.email,
    u.created_at,
    u.updated_at
FROM users u
LEFT JOIN image_conversions p ON p.user_id = u.user_id
WHERE u.user_id = $1 
LIMIT 1;

-- name: GetUser :one
SELECT  
    u.user_id,
    u.user_name,
    u.first_name,
    u.last_name,
    u.email,
    u.created_at,
    u.updated_at
FROM users u
WHERE u.user_id = $1 
LIMIT 1;

-- name: UpdateUsers :one
UPDATE users
SET
    user_name = COALESCE(sqlc.narg(user_name), user_name),
    first_name = COALESCE(sqlc.narg(first_name), first_name),
    last_name = COALESCE(sqlc.narg(last_name), last_name),
    email = COALESCE(sqlc.narg(email), email),
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    updated_at = COALESCE(sqlc.narg(updated_at), updated_at)
WHERE 
    user_id = sqlc.arg(user_id)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users
WHERE user_id = $1;
