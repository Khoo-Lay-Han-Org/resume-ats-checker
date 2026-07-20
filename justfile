set dotenv-load := true
set positional-arguments := true
set shell := ["bash", "-cu"]

# Database connection — all required, fail fast if missing
_db-host := env_var("DATABASE_HOST")
_db-user := env_var("DATABASE_USERNAME")
_db-name := env_var("DATABASE_NAME")
_db-pass := env_var("DATABASE_PASSWORD")
_db-port := env_var("DATABASE_PORT")
_db-conn-string := "host=" + _db-host + " user=" + _db-user + " dbname=" + _db-name + " password=" + _db-pass + " port=" + _db-port + " sslmode=require"


# Run the backend server
[working-directory: 'apps/backend']
run-backend:
    go run main.go

# Generate SQL code from sqlc
[working-directory: 'apps/backend']
backend-sqlc-generate:
    sqlc generate

# Run all pending migrations
[working-directory: 'apps/backend']
backend-migrate-up:
    goose -dir migrations postgres '{{ _db-conn-string }}' up

# Reset the database (rollback all migrations)
[working-directory: 'apps/backend']
backend-migrate-reset:
    goose -dir migrations postgres '{{ _db-conn-string }}' reset

# Check migration status
[working-directory: 'apps/backend']
backend-migrate-status:
    goose -dir migrations postgres '{{ _db-conn-string }}' status
