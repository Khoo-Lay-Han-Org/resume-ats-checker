-- name: FindUserByPublicIdWithSession :one
SELECT u.*, s.id as session_id, s.public_id as session_public_id, s.session_key, s.expires_at as session_expires_at
FROM users u
LEFT JOIN sessions s ON s.user_id = u.id AND (s.expires_at IS NULL OR s.expires_at > NOW())
WHERE u.public_id = @public_id AND (u.expires_at IS NULL OR u.expires_at > NOW()) LIMIT 1;

-- name: FindAllClientsWithSession :many
SELECT u.*, s.id as session_id, s.public_id as session_public_id, s.session_key, s.expires_at as session_expires_at
FROM users u
LEFT JOIN sessions s ON s.user_id = u.id AND (s.expires_at IS NULL OR s.expires_at > NOW())
WHERE u.user_type = 'client'::user_type AND (u.expires_at IS NULL OR u.expires_at > NOW());

-- name: FindAllAdminsWithSession :many
SELECT u.*, s.id as session_id, s.public_id as session_public_id, s.session_key, s.expires_at as session_expires_at
FROM users u
LEFT JOIN sessions s ON s.user_id = u.id AND (s.expires_at IS NULL OR s.expires_at > NOW())
WHERE u.user_type IN ('admin'::user_type, 'super-admin'::user_type) AND (u.expires_at IS NULL OR u.expires_at > NOW());

-- name: FindUserByEmailWithExpired :one
SELECT * FROM users WHERE email = @email LIMIT 1;
