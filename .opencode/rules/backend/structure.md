# Backend Structure Convention

## Please understand that this is for apps/backend/ ONLY

### Naming Convention

1) File names use small capital letters and words are separated via -
2) Folder names use PascalCase

### Root Files and Folders Organisation

1) Ask before changing the files and folders structure at the root of backend

### Child Level Files and Folders Organisation

1) All codes are organised into their specific folders with their specific folder names
2) Do not use general naming for files and folders, such as "util"
3) Each feature domain follows: api/<domain>/{typing,validator,util,view}
4) database/sqlc/ is sqlc-generated (DO NOT EDIT)
5) systemconfig/ holds all configuration types (backend, frontend, database, cache, email, session, ai, oauth)
6) tool/ holds external service clients (Valkey, Mailgun, OAuth2)
7) env/ holds env loading helpers
8) sql/ holds raw SQL for sqlc codegen
9) migrations/ holds Goose migration files
10) cmd/docs/ generates Swagger docs
11) docs/ holds generated Swagger JSON/YAML
12) test/ holds integration tests
