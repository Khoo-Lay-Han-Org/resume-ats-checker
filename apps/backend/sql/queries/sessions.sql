-- name: FindSessionByUserId :one
SELECT * FROM sessions
WHERE user_id = @user_id AND (expires_at IS NULL OR expires_at > NOW()) LIMIT 1;

-- name: FindSessionByPublicId :one
SELECT * FROM sessions
WHERE public_id = @public_id AND (expires_at IS NULL OR expires_at > NOW()) LIMIT 1;

-- name: FindAllSessions :many
SELECT * FROM sessions
WHERE (expires_at IS NULL OR expires_at > NOW());

-- name: FindAllSessionsWithUser :many
SELECT s.*, u.public_id as user_public_id, u.username, u.displayname, u.user_type
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE (s.expires_at IS NULL OR s.expires_at > NOW());

-- name: CreateSession :one
INSERT INTO sessions (user_id, session_key)
VALUES (@user_id, @session_key)
RETURNING *;

-- name: DeleteSessionByUserId :exec
DELETE FROM sessions WHERE user_id = @user_id;

-- name: DeleteAllSessions :exec
DELETE FROM sessions;

-- name: UpdateSession :exec
UPDATE sessions SET session_key = @session_key, updated_at = NOW()
WHERE user_id = @user_id;
