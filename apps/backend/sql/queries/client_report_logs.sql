-- name: FindAllClientReportLogs :many
SELECT * FROM client_report_logs
WHERE (expires_at IS NULL OR expires_at > NOW())
ORDER BY created_at DESC;

-- name: FindClientReportLogByPublicId :one
SELECT * FROM client_report_logs
WHERE public_id = @public_id LIMIT 1;

-- name: CreateClientReportLog :one
INSERT INTO client_report_logs (reporting_user_id, target_user_id, type)
VALUES (@reporting_user_id, @target_user_id, @type)
RETURNING *;

-- name: UpdateClientReportLogByPublicId :exec
UPDATE client_report_logs SET type = @type, updated_at = NOW()
WHERE public_id = @public_id;

-- name: FindClientReportLogWithUsers :many
SELECT crl.*, ru.public_id as reporting_user_public_id, tu.public_id as target_user_public_id
FROM client_report_logs crl
JOIN users ru ON ru.id = crl.reporting_user_id
JOIN users tu ON tu.id = crl.target_user_id
WHERE (crl.expires_at IS NULL OR crl.expires_at > NOW())
ORDER BY crl.created_at DESC;
