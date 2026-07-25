set dotenv-load := true
set positional-arguments := true
set shell := ["bash", "-cu"]

# Database connection — all required, fail fast if missing
db-host := env_var("DATABASE_HOST")
db-user := env_var("DATABASE_USERNAME")
db-name := env_var("DATABASE_NAME")
db-pass := env_var("DATABASE_PASSWORD")
db-port := env_var("DATABASE_PORT")
db-conn-string := "host=" + db-host + " user=" + db-user + " dbname=" + db-name + " password=" + db-pass + " port=" + db-port + " sslmode=require"


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
    goose -dir database/migrations postgres '{{ db-conn-string }}' up

# Reset the database (rollback all migrations)
[working-directory: 'apps/backend']
backend-migrate-reset:
    goose -dir database/migrations postgres '{{ db-conn-string }}' reset

# Check migration status
[working-directory: 'apps/backend']
backend-migrate-status:
    goose -dir database/migrations postgres '{{ db-conn-string }}' status
