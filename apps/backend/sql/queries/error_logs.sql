-- name: FindAllErrorLogs :many
SELECT * FROM error_logs
WHERE (expires_at IS NULL OR expires_at > NOW())
ORDER BY created_at DESC;

-- name: CreateErrorLog :one
INSERT INTO error_logs (user_id, type, message)
VALUES (@user_id, @type, @message)
RETURNING *;

-- name: UpdateErrorLogByPublicId :exec
UPDATE error_logs SET type = @type, message = @message, updated_at = NOW()
WHERE public_id = @public_id;
