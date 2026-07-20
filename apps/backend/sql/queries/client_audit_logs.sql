-- name: FindClientAuditLogsByUserId :many
SELECT * FROM client_audit_logs
WHERE user_id = @user_id AND (expires_at IS NULL OR expires_at > NOW())
ORDER BY created_at DESC;

-- name: FindAllClientAuditLogs :many
SELECT * FROM client_audit_logs
WHERE (expires_at IS NULL OR expires_at > NOW())
ORDER BY created_at DESC;

-- name: CreateClientAuditLog :one
INSERT INTO client_audit_logs (user_id, type, message)
VALUES (@user_id, @type, @message)
RETURNING *;

-- name: UpdateClientAuditLogByPublicId :exec
UPDATE client_audit_logs SET type = @type, message = @message, updated_at = NOW()
WHERE public_id = @public_id;
