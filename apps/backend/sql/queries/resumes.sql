-- name: FindResumeByUserId :one
SELECT * FROM resumes
WHERE user_id = @user_id AND (expires_at IS NULL OR expires_at > NOW()) LIMIT 1;

-- name: FindAllResumes :many
SELECT * FROM resumes
WHERE (expires_at IS NULL OR expires_at > NOW());

-- name: UpsertResume :one
INSERT INTO resumes (user_id, template_id, detail)
VALUES (@user_id, @template_id, @detail)
ON CONFLICT (user_id) WHERE (expires_at IS NULL OR expires_at > NOW())
DO UPDATE SET template_id = @template_id, detail = @detail, updated_at = NOW()
RETURNING *;

-- name: UpdateResume :exec
UPDATE resumes SET
  template_id = @template_id,
  detail = @detail,
  updated_at = NOW()
WHERE user_id = @user_id;
