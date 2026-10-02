# Backend Best Coding Practice

## Please understand that this is for apps/backend/ ONLY

1) Ask the user for more clarity and avoid guessing
2) Look at the existing related codes for convention examples
3) Do a study on best practices before coding
4) Always devise a clear plan and ask user for approval of the plan first

## Go Conventions

1) Use context.Context for cancellation and timeouts
2) Handle errors explicitly; do not ignore them
3) Use Echo context for data passing between handlers (c.Set/c.Get)
4) Wrap handler logic in util/ and keep view/ thin
5) Use sqlc-generated queries; do not write raw SQL in Go files
6) Valkey keys follow the pattern: {public_user_id}:<data>_data
7) Async database sync must be done in goroutines to avoid blocking API responses
8) Session validation middleware checks Valkey then PostgreSQL fallback
