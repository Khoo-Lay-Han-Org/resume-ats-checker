-- name: FindAdminAuditLogsByUserId :many
SELECT * FROM admin_audit_logs
WHERE user_id = @user_id AND (expires_at IS NULL OR expires_at > NOW())
ORDER BY created_at DESC;

-- name: FindAllAdminAuditLogs :many
SELECT * FROM admin_audit_logs
WHERE (expires_at IS NULL OR expires_at > NOW())
ORDER BY created_at DESC;

-- name: CreateAdminAuditLog :one
INSERT INTO admin_audit_logs (user_id, type, message)
VALUES (@user_id, @type, @message)
RETURNING *;

-- name: UpdateAdminAuditLogByPublicId :exec
UPDATE admin_audit_logs SET type = @type, message = @message, updated_at = NOW()
WHERE public_id = @public_id;
