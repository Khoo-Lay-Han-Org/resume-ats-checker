-- name: FindAtsByUserId :one
SELECT * FROM ats
WHERE user_id = @user_id AND (expires_at IS NULL OR expires_at > NOW()) LIMIT 1;

-- name: FindAllAts :many
SELECT * FROM ats
WHERE (expires_at IS NULL OR expires_at > NOW());

-- name: UpsertAts :one
INSERT INTO ats (user_id, score, reasoning)
VALUES (@user_id, @score, @reasoning)
ON CONFLICT (user_id) WHERE (expires_at IS NULL OR expires_at > NOW())
DO UPDATE SET score = @score, reasoning = @reasoning, updated_at = NOW()
RETURNING *;

-- name: UpdateAts :exec
UPDATE ats SET
  score = @score,
  reasoning = @reasoning,
  updated_at = NOW()
WHERE user_id = @user_id;
