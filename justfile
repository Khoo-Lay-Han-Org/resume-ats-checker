set dotenv-load := true
set positional-arguments := true
set shell := ["bash", "-cu"]

db-host := env_var("DATABASE_HOST")
db-user := env_var("DATABASE_USERNAME")
db-name := env_var("DATABASE_NAME")
db-pass := env_var("DATABASE_PASSWORD")
db-port := env_var("DATABASE_PORT")
db-conn-string := "host=" + db-host + " user=" + db-user + " dbname=" + db-name + " password=" + db-pass + " port=" + db-port + " sslmode=require"

# BACKEND
[working-directory: 'apps/backend']
run-backend:
    go run main.go

[working-directory: 'apps/backend']
backend-sqlc-generate:
    sqlc generate

[working-directory: 'apps/backend']
backend-migrate-up:
    doppler run -- goose -dir database/migrations postgres '{{ db-conn-string }}' up

[working-directory: 'apps/backend']
backend-migrate-reset:
    doppler run -- goose -dir database/migrations postgres '{{ db-conn-string }}' reset

[working-directory: 'apps/backend']
backend-migrate-status:
    doppler run -- goose -dir database/migrations postgres '{{ db-conn-string }}' status

# AI-ML
[working-directory: 'apps/ai-ml']
ai-ml-lint:
	uv run ruff check --fix .

[working-directory: 'apps/ai-ml']
ai-ml-format:
	uv run black .

[working-directory: 'apps/ai-ml']
ai-ml-security:
	uv run bandit -r . -x ./.venv

