-- name: FindJwtKeyByUserId :one
SELECT * FROM jwt_keys
WHERE user_id = @user_id AND (expires_at IS NULL OR expires_at > NOW()) LIMIT 1;

-- name: FindAllJwtKeys :many
SELECT * FROM jwt_keys
WHERE (expires_at IS NULL OR expires_at > NOW());

-- name: CreateJwtKey :one
INSERT INTO jwt_keys (user_id, key)
VALUES (@user_id, @key)
RETURNING *;

-- name: DeleteJwtKeyByUserId :exec
DELETE FROM jwt_keys WHERE user_id = @user_id;

-- name: UpdateJwtKey :exec
UPDATE jwt_keys SET key = @key, updated_at = NOW()
WHERE user_id = @user_id;
