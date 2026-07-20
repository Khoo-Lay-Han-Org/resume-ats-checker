-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_type AS ENUM ('super-admin', 'admin', 'client');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password BYTEA NOT NULL,
    displayname VARCHAR(255) NOT NULL UNIQUE,
    user_type user_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    banned_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_users_public_id ON users(public_id);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_updated_at ON users(updated_at);
CREATE INDEX idx_users_banned_at ON users(banned_at);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_expires_at ON users(expires_at);

CREATE TABLE showcase_records (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255)[],
    email VARCHAR(255)[],
    phone_number VARCHAR(255)[],
    address VARCHAR(255)[],
    social_media VARCHAR(255)[],
    job_experience JSONB,
    education JSONB,
    skill VARCHAR(255)[],
    certificate JSONB,
    language VARCHAR(255)[],
    project JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_showcase_records_public_id ON showcase_records(public_id);
CREATE INDEX idx_showcase_records_user_id ON showcase_records(user_id);

CREATE TABLE portfolios (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    template_id INT NOT NULL,
    detail JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_portfolios_public_id ON portfolios(public_id);
CREATE INDEX idx_portfolios_user_id ON portfolios(user_id);
CREATE INDEX idx_portfolios_template_id ON portfolios(template_id);

CREATE TABLE resumes (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    template_id INT NOT NULL,
    detail JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_resumes_public_id ON resumes(public_id);
CREATE INDEX idx_resumes_user_id ON resumes(user_id);
CREATE INDEX idx_resumes_template_id ON resumes(template_id);

CREATE TABLE ats (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    score INT NOT NULL,
    reasoning VARCHAR(355) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_ats_public_id ON ats(public_id);
CREATE INDEX idx_ats_user_id ON ats(user_id);

CREATE TABLE jwt_keys (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id INT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    key VARCHAR(355) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_jwt_keys_public_id ON jwt_keys(public_id);
CREATE INDEX idx_jwt_keys_user_id ON jwt_keys(user_id);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id INT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    session_key VARCHAR(355) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP + INTERVAL '3 days')
);

CREATE INDEX idx_sessions_public_id ON sessions(public_id);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);

CREATE TABLE client_audit_logs (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(355) NOT NULL,
    message VARCHAR(355) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_client_audit_logs_public_id ON client_audit_logs(public_id);
CREATE INDEX idx_client_audit_logs_created_at ON client_audit_logs(created_at);

CREATE TABLE admin_audit_logs (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(355) NOT NULL,
    message VARCHAR(355) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_admin_audit_logs_public_id ON admin_audit_logs(public_id);
CREATE INDEX idx_admin_audit_logs_created_at ON admin_audit_logs(created_at);

CREATE TABLE client_report_logs (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    reporting_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(355) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_client_report_logs_public_id ON client_report_logs(public_id);
CREATE INDEX idx_client_report_logs_reporting_user_id ON client_report_logs(reporting_user_id);
CREATE INDEX idx_client_report_logs_target_user_id ON client_report_logs(target_user_id);

CREATE TABLE error_logs (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type INT NOT NULL,
    message VARCHAR(355) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_error_logs_public_id ON error_logs(public_id);
CREATE INDEX idx_error_logs_created_at ON error_logs(created_at);

CREATE TABLE client_support_messaging (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    type VARCHAR(355) NOT NULL,
    content JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_client_support_messaging_public_id ON client_support_messaging(public_id);
CREATE INDEX idx_client_support_messaging_created_at ON client_support_messaging(created_at);

-- +goose Down
DROP TABLE IF EXISTS client_support_messaging CASCADE;
DROP TABLE IF EXISTS error_logs CASCADE;
DROP TABLE IF EXISTS client_report_logs CASCADE;
DROP TABLE IF EXISTS admin_audit_logs CASCADE;
DROP TABLE IF EXISTS client_audit_logs CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS jwt_keys CASCADE;
DROP TABLE IF EXISTS ats CASCADE;
DROP TABLE IF EXISTS resumes CASCADE;
DROP TABLE IF EXISTS portfolios CASCADE;
DROP TABLE IF EXISTS showcase_records CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP TYPE IF EXISTS user_type;
