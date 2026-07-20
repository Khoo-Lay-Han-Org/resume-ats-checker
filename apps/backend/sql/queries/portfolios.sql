-- name: FindPortfolioByUserId :one
SELECT * FROM portfolios
WHERE user_id = @user_id AND (expires_at IS NULL OR expires_at > NOW()) LIMIT 1;

-- name: FindAllPortfolios :many
SELECT * FROM portfolios
WHERE (expires_at IS NULL OR expires_at > NOW());

-- name: UpsertPortfolio :one
INSERT INTO portfolios (user_id, template_id, detail)
VALUES (@user_id, @template_id, @detail)
ON CONFLICT (user_id) WHERE (expires_at IS NULL OR expires_at > NOW())
DO UPDATE SET template_id = @template_id, detail = @detail, updated_at = NOW()
RETURNING *;

-- name: UpdatePortfolio :exec
UPDATE portfolios SET
  template_id = @template_id,
  detail = @detail,
  updated_at = NOW()
WHERE user_id = @user_id;
