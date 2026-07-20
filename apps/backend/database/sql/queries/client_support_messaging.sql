-- name: FindAllClientSupportMessages :many
SELECT * FROM client_support_messaging
WHERE (expires_at IS NULL OR expires_at > NOW())
ORDER BY created_at DESC;

-- name: FindClientSupportMessageByPublicId :one
SELECT * FROM client_support_messaging
WHERE public_id = @public_id LIMIT 1;

-- name: CreateClientSupportMessage :one
INSERT INTO client_support_messaging (type, content)
VALUES (@type, @content)
RETURNING *;

-- name: UpdateClientSupportMessageByPublicId :exec
UPDATE client_support_messaging SET type = @type, content = @content, updated_at = NOW()
WHERE public_id = @public_id;
