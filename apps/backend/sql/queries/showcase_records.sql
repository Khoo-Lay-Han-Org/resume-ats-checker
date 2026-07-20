-- name: FindShowcaseRecordByUserId :one
SELECT * FROM showcase_records
WHERE user_id = @user_id AND (expires_at IS NULL OR expires_at > NOW()) LIMIT 1;

-- name: FindShowcaseRecordByPublicId :one
SELECT * FROM showcase_records
WHERE public_id = @public_id AND (expires_at IS NULL OR expires_at > NOW()) LIMIT 1;

-- name: FindAllShowcaseRecords :many
SELECT * FROM showcase_records
WHERE (expires_at IS NULL OR expires_at > NOW());

-- name: CreateShowcaseRecord :one
INSERT INTO showcase_records (user_id)
VALUES (@user_id)
RETURNING *;

-- name: UpdateShowcaseRecord :exec
UPDATE showcase_records SET
  name = @name,
  email = @email,
  phone_number = @phone_number,
  address = @address,
  social_media = @social_media,
  job_experience = @job_experience,
  education = @education,
  skill = @skill,
  certificate = @certificate,
  language = @language,
  project = @project,
  updated_at = NOW()
WHERE user_id = @user_id;

-- name: UpdateShowcaseRecordByPublicId :exec
UPDATE showcase_records SET
  name = @name,
  email = @email,
  phone_number = @phone_number,
  address = @address,
  social_media = @social_media,
  job_experience = @job_experience,
  education = @education,
  skill = @skill,
  certificate = @certificate,
  language = @language,
  project = @project,
  updated_at = NOW()
WHERE public_id = @public_id;
