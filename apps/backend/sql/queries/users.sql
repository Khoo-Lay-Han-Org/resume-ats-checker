-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = @email AND (expires_at IS NULL OR expires_at > NOW()) LIMIT 1;

-- name: FindUserByPublicId :one
SELECT * FROM users
WHERE public_id = @public_id AND (expires_at IS NULL OR expires_at > NOW()) LIMIT 1;

-- name: FindUserById :one
SELECT * FROM users
WHERE id = @id AND (expires_at IS NULL OR expires_at > NOW()) LIMIT 1;

-- name: FindAllUsersByType :many
SELECT * FROM users
WHERE user_type = @user_type::user_type AND (expires_at IS NULL OR expires_at > NOW());

-- name: FindAllUsers :many
SELECT * FROM users
WHERE (expires_at IS NULL OR expires_at > NOW());

-- name: CountUsersByType :one
SELECT COUNT(*) FROM users
WHERE user_type = @user_type::user_type;

-- name: CreateUser :one
INSERT INTO users (username, email, password, displayname, user_type)
VALUES (@username, @email, @password, @displayname, @user_type::user_type)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users SET 
    username = @username, 
    email = @email, 
    displayname = @displayname, 
    password = @password, 
    user_type = @user_type::user_type, 
    banned_at = @banned_at, 
    deleted_at = @deleted_at, 
    updated_at = NOW()
WHERE id = @id;

-- name: UpdateUserByPublicId :one
UPDATE users SET 
    username = @username, 
    email = @email, 
    displayname = @displayname, 
    user_type = @user_type::user_type, 
    banned_at = @banned_at, 
    deleted_at = @deleted_at, 
    updated_at = NOW()
WHERE public_id = @public_id
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id;

-- name: HardDeleteExpiredUsers :exec
DELETE FROM users WHERE expires_at IS NOT NULL AND expires_at < NOW();
