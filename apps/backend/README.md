# Resuming Backend

> **Resuming** is a full-stack resume builder and portfolio platform with ATS (Applicant Tracking System) scoring, client-support messaging, and admin audit logging — all served through a type-safe Go API.

---

## Table of Contents

1. [Tech Stack](#tech-stack)
2. [Project Structure](#project-structure)
3. [How to Run](#how-to-run)
4. [Architecture Overview](#architecture-overview)
5. [Authentication & Sessions](#authentication--sessions)
6. [Valkey ↔ PostgreSQL Sync](#valkey--postgresql-sync)
7. [Role-Based Access Control](#role-based-access-control)
8. [ATS Scoring Pipeline](#ats-scoring-pipeline)
9. [API Routes](#api-routes)
10. [Environment Variables](#environment-variables)
11. [Database Migrations](#database-migrations)
12. [Testing](#testing)
13. [Key Patterns](#key-patterns)

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.26 |
| HTTP Framework | Echo v4 |
| Database | PostgreSQL (via `pgx/v5`) |
| Query Builder | `sqlc` (type-safe generated queries) |
| Cache / Session Store | Valkey (Redis-compatible, via `valkey-go`) |
| Auth | JWT (`golang-jwt/jwt/v5`), bcrypt, OTP via email |
| Email | Mailgun |
| Validation | `valtra-go` |
| RBAC | Casbin |
| PDF Parsing | `ledongthuc/pdf` |
| Browser Automation | `go-rod` (for ATS web scraping) |
| ID Generation | `segmentio/ksuid`, `google/uuid` |
| API Docs | Swagger (`swaggo/swag`) |
| Migrations | Goose |

---

## Project Structure

```
apps/backend/
├── main.go                    # Entrypoint: connects Valkey, DB, starts scheduler, launches Echo
├── go.mod / go.sum            # Go module: "resuming"
├── sqlc.yaml                  # sqlc config (PostgreSQL engine, pgx/v5 driver)
│
├── api/                       # ── HTTP Layer ──────────────────────────────
│   ├── router.go              # All route definitions & middleware groups
│   ├── wrapper.go             # Orchestrates multi-step handler chains
│   ├── docs.go                # Swagger annotations
│   ├── response-types-doc.go  # Shared response doc types
│   ├── doc-router.go          # Swagger UI route
│   │
│   ├── middleware/             # Echo middleware
│   │   ├── session.go         # Reads "session" cookie, validates JWT in Valkey
│   │   ├── user-lvl.go        # Casbin RBAC: OnlyAdmin, OnlySuperAdmin
│   │   ├── util/
│   │   │   └── session.go     # ExtractSessionCookie, ParseJWT, CheckSession
│   │   └── config/            # Casbin model & policy CSV files
│   │
│   ├── auth/                  # Registration, login, OTP, OAuth (stub)
│   │   ├── typing/            # Request DTOs (Register, Login, OTP)
│   │   ├── validator/         # Input validation (valtra-go)
│   │   ├── util/              # SendOTP, CheckOTP, FindUser, ValidateEmailMX
│   │   └── view/              # Handlers: PrepareRegistration, Register, Login, SetSession
│   │
│   ├── showcaserecord/        # User showcase data (name, email, skills, etc.)
│   ├── portfolio/             # Portfolio template & content
│   ├── resume/                # Resume template & content
│   ├── ats/                   # ATS scoring engine
│   ├── setting/               # Account settings (username, email, password, delete)
│   ├── client-support/        # Client-to-admin messaging & reporting
│   └── administrator/         # Admin panel (ban, sessions, audit logs, invites)
│
├── database/                  # ── Data Layer ─────────────────────────────
│   ├── config.go              # pgxpool connection, sqlc Queries init
│   ├── column-typing.go       # Custom column type definitions
│   ├── util.go                # Helpers: FindUser, FindUserByPublicId, resolveUserId
│   ├── session-store-typing.go
│   ├── sqlc/                  # sqlc-generated code (DO NOT EDIT)
│   │   ├── db.go, querier.go, models.go
│   │   ├── users.sql.go, sessions.sql.go, jwt_keys.sql.go
│   │   ├── showcase_records.sql.go, portfolios.sql.go, resumes.sql.go
│   │   ├── ats.sql.go
│   │   ├── client_audit_logs.sql.go, admin_audit_logs.sql.go
│   │   ├── client_report_logs.sql.go, client_support_messaging.sql.go
│   │   └── error_logs.sql.go
│   ├── sync-group-cache.go    # Bulk-sync entire tables → Valkey
│   ├── sync-group-database.go # Bulk-sync entire tables → PostgreSQL
│   ├── sync-individual-cache.go    # Per-user sync → Valkey
│   └── sync-individual-database.go # Per-user sync → PostgreSQL
│
├── system-config/             # ── Configuration ──────────────────────────
│   ├── backend.go             # BackendDomain, BackendPort, BackendUri
│   ├── frontend.go            # FrontendDomain, FrontendPort, FrontendUri
│   ├── database.go            # Database DSN construction
│   ├── cache.go               # Valkey connection config
│   ├── email.go               # Mailgun config
│   ├── session.go             # SessionExpiryDuration (3d), OtpExpiryDuration (5m)
│   ├── ai.go                  # AI models service URI
│   ├── oauth.go               # OAuth config (commented out / stub)
│   └── development-stage.go   # ApplicationHosted flag (local vs PaaS)
│
├── scheduler/                 # ── Background Jobs ────────────────────────
│   └── schedule.go            # FirstSync (startup) + FullSync (every 10 min)
│
├── tool/                      # ── External Clients ───────────────────────
│   ├── valkey.go              # Valkey client setup (TLS support)
│   ├── mailgun.go             # SendEmail via Mailgun HTTP API
│   └── oauth2.go              # Google OAuth2 config
│
├── env/                       # ── Env Loading ────────────────────────────
│   └── env.go                 # godotenv.Load("../.env") helper
│
├── sql/                       # Raw SQL for sqlc codegen
│   └── queries/*.sql
│
├── migrations/                # Goose migration files
│   └── 001_initial_schema.sql
│
├── cmd/
│   └── docs/main.go           # Generates Swagger docs via swaggo
│
├── docs/
│   └── swagger.json / .yaml   # Generated API docs
│
└── test/                      # Integration tests (auth, admin, settings, etc.)
```

---

## How to Run

### Prerequisites
- Go 1.26+
- PostgreSQL
- Valkey (or Redis-compatible service like Upstash)
- Mailgun account (for OTP emails)
- `.env` file at the project root (`../.env` relative to backend)

### Install & Start

```bash
# Install dependencies
cd apps/backend
go mod download

# Run the server
go run main.go
```

Server starts on port `5321` by default (or `PORT` env var if set).

### Generate Swagger Docs

```bash
cd cmd/docs
go run main.go
```

### Lint / Format

```bash
bun x ultracite check   # check for issues
bun x ultracite fix     # auto-fix
```

---

## Architecture Overview

### Layered Domain Design

Each feature domain (`auth`, `showcaserecord`, `portfolio`, `resume`, `ats`, `setting`, `client-support`, `administrator`) follows an identical internal structure:

```
api/<domain>/
├── typing/       # Request/response DTOs (plain structs)
├── validator/    # Input validation using valtra-go
├── util/         # Business logic helpers (DB ops, AI calls, etc.)
└── view/         # Echo handler functions (return echo.HandlerFunc)
```

**Flow**: `router.go` → `wrapper.go` → `view handler` → `util` → `database` → `sqlc` → PostgreSQL.

### Wrapper Pattern

`wrapper.go` contains top-level handler functions that **orchestrate multi-step pipelines**. For example, `ATSScoreWebScrapeFlow` calls 11 view functions in sequence:

```go
func ATSScoreWebScrapeFlow(c echo.Context) error {
    ats_view.ExtractResume()(c)
    ats_view.ParseResume()(c)
    ats_view.SectionExistenceCheck()(c)
    // ... 8 more steps ...
    data := c.Get("response_data")
    return c.JSON(http.StatusOK, echo.Map{"message": "...", "data": data})
}
```

Each step reads/writes data to the Echo `context` via `c.Set()`/`c.Get()`. This makes the pipeline explicit and debuggable.

---

## Authentication & Sessions

### Two-Step Flows

Both registration and login use a **two-step OTP (One-Time Password) flow**:

1. **Prepare** (`/prepare-registeration` or `/prepare-login`):
   - Validate input, hash password
   - Store registration data in Valkey (`email:session` hash)
   - Generate 6-digit OTP, store bcrypt-hashed OTP in Valkey (`email:otp`, 5-min TTL)
   - Send OTP via Mailgun email
   - Set `email_for_otp` cookie

2. **Verify** (`/register/:type-of-user` or `/login`):
   - Read `email_for_otp` cookie
   - Validate OTP against bcrypt-hashed value in Valkey
   - On registration: create user in PostgreSQL
   - On login: proceed to session creation

### Session Creation

After successful OTP verification, `SetSession()` generates:

- **`session_key`**: random KSUID (stored in PostgreSQL `sessions` table + Valkey)
- **`signing_key`**: random KSUID (stored in PostgreSQL `jwt_keys` table + Valkey)
- **JWT**: Signed with `signing_key`, claims include `user_public_id` and `exp` (3 days)
- **Session cookie**: `{"public_id": "...", "session_key": "...", "token": "..."}` (HttpOnly)

### Session Validation Middleware

`SessionCheck()` middleware (applied to all `/authed` routes):
1. Reads `session` cookie
2. Parses JWT from cookie, looking up the signing key from Valkey (`public_user_id:jwt_data`)
3. Verifies session exists in Valkey (`public_user_id:session_data`)
4. Sets `public_user_id` on Echo context for downstream handlers

On failure: deletes the session cookie and returns 401.

---

## Valkey ↔ PostgreSQL Sync

This is the most architecturally significant pattern in the backend. Valkey serves as the **hot/cache store**, while PostgreSQL is the **persistent/cold store**.

### Why Two Stores?

- **Valkey**: Fast reads/writes for session data, OTPs, and user data cached at login time
- **PostgreSQL**: Durable, queryable source of truth

### Two Sync Directions

#### 1. Cache → Database (Persistence)

- **Individual sync**: Triggered on login (`SyncIndividualLoginDataToSessionStore`) and after data mutations (e.g., username change). Syncs a single user's data from Valkey to PostgreSQL asynchronously via goroutines.
- **Group sync**: `scheduler.FullSync()` runs every **10 minutes**, bulk-syncing all tables from Valkey to PostgreSQL.

#### 2. Database → Cache (Warm-up)

- **Individual sync**: `SyncIndividualUserDataSessionStore()` and similar functions load from PostgreSQL and write to Valkey (used as fallback when Valkey cache misses).
- **Group sync**: `scheduler.FirstSync()` runs at startup, then `FullSync()` periodically bulk-loads entire tables into Valkey.

### Key Naming Convention in Valkey

```
{public_user_id}:user_data          # Full user profile (for RBAC checks)
{public_user_id}:session_data       # {"session_key": "..."}
{public_user_id}:jwt_data           # JwtKey struct
{public_user_id}:showcaserecord_data
{public_user_id}:portfolio_data
{public_user_id}:resume_data
{public_user_id}:ats_data
{public_user_id}:client_audit_log_data
{public_user_id}:admin_audit_log_data
{public_user_id}:error_log_data
{public_user_id}:client_report_log_data
{public_user_id}:client_support_messages
{public_user_id}:session_id         # Maps public_user_id → session public_id

# Global keys (group sync targets):
error_log_data
client_audit_log_data
admin_audit_log_data
client_configs
admin_configs
client_report_logs
client_support_messages
showcaserecord_data
portfolio_data
resume_data
ats_data
user_data
```

All keys use a TTL matching `SessionExpiryDuration` (3 days), so stale cache naturally expires.

---

## Role-Based Access Control

Uses **Casbin** with two policy files:

| Middleware | Model File | Policy File | Purpose |
|---|---|---|---|
| `OnlyAdmin()` | `api/middleware/config/user-lvl.conf` | `api/middleware/config/user-lvl.csv` | Allows `admin` and `super-admin` |
| `OnlySuperAdmin()` | `api/middleware/config/super-admin-lvl.conf` | `api/middleware/config/super-admin-lvl.csv` | Allows only `super-admin` |

### User Types (PostgreSQL Enum)

```
super-admin  → Full system access (invite admins, remove admins)
admin        → Client management, audit logs, support messaging
client       → Regular users (create resumes, portfolios, ATS scoring)
```

### Middleware Chain

```
SessionCheck() → OnlyAdmin() → OnlySuperAdmin()
```

Route groups are nested in `router.go`:
```go
authed := router.Group("")
authed.Use(middleware.SessionCheck())
{
    is_admin := authed.Group("")
    is_admin.Use(middleware.OnlyAdmin())
    {
        // admin routes ...
        is_super_admin := is_admin.Group("")
        is_super_admin.Use(middleware.OnlySuperAdmin())
        {
            // super-admin routes ...
        }
    }
}
```

---

## ATS Scoring Pipeline

The ATS (Applicant Tracking System) scoring feature evaluates how well a resume matches a job description. Two entry points:

- `POST /ats-score-webscrape` — scrapes job description from a URL
- `POST /ats-score-user-input` — uses user-provided job description text

Both run the same 11-step pipeline (see `wrapper.go:122-204`):

| Step | Handler | Description |
|---|---|---|
| 1 | `ExtractResume` | Pulls resume data from user's showcaserecord |
| 2 | `ParseResume` | Parses structured sections |
| 3 | `SectionExistenceCheck` | Verifies all required sections exist |
| 4 | `FormattingCheck` | Checks formatting quality |
| 5 | `WebScrapeJobDesc` / `UserInputJobDesc` | Gets job description |
| 6 | `ResumeJobTypeCheck` | Infers job type from resume |
| 7 | `JobDescJobTypeCheck` | Infers job type from job description |
| 8 | `JobTypeRelevanceCheck` | Compares job type match |
| 9 | `ResumeSkillsCheck` | Extracts skills from resume |
| 10 | `JobDescSkillsCheck` | Extracts skills from job description |
| 11 | `OverallSkillsCheck` | Compares skill overlap |

Final step `OverallScore` averages all sub-scores into a single integer (0-100).

---

## API Routes

### Public Routes

| Method | Path | Handler |
|---|---|---|
| POST | `/prepare-registeration` | Step 1 of registration (send OTP) |
| POST | `/register/:type-of-user` | Step 2 of registration (verify OTP, create user) |
| POST | `/prepare-login` | Step 1 of login (send OTP) |
| POST | `/login` | Step 2 of login (verify OTP) |
| GET | `/accept-become-admin/:token` | Accept admin invitation |

### Authenticated Routes (require session cookie)

| Method | Path | Description |
|---|---|---|
| POST | `/showcaserecord-add/:type-of-data` | Add showcase record section |
| DELETE | `/showcaserecord-delete` | Delete showcase record |
| PATCH | `/showcaserecord-edit/:type-of-data` | Edit showcase record section |
| GET | `/showcaserecord-retrieve` | Get all showcase records |
| PATCH | `/choose-portfolio-template` | Select portfolio template |
| GET | `/get-portfolio-content` | Get portfolio data |
| PATCH | `/choose-resume-template` | Select resume template |
| GET | `/get-resume-content` | Get resume data |
| POST | `/ats-score-webscrape` | ATS score via URL scraping |
| POST | `/ats-score-user-input` | ATS score via text input |
| POST | `/change-username` | Update username |
| POST | `/change-displayname` | Update display name |
| POST | `/prepare-change-email` | Request email change OTP |
| POST | `/change-email` | Confirm email change |
| POST | `/prepare-change-password` | Request password change OTP |
| POST | `/change-password` | Confirm password change |
| POST | `/prepare-delete-account` | Request account deletion OTP |
| POST | `/delete-account` | Confirm account deletion |
| POST | `/client_comm_to_admin` | Send message to admin |
| POST | `/client_report_other_client` | Report another user |

### Admin Routes (require admin role)

| Method | Path | Description |
|---|---|---|
| POST | `/ban_client` | Ban a client |
| POST | `/remove_individual_session` | Revoke one user's session |
| POST | `/remove_all_session` | Revoke all of a client's sessions |
| GET | `/client_comm_log` | Get support messages |
| POST | `/client_comm_reply_log` | Reply to support message |
| GET | `/get_all_clients` | List all clients |
| GET | `/get_all_admins` | List all admins |
| GET | `/client_audit_logs` | Client audit trail |
| GET | `/admin_audit_logs` | Admin audit trail |
| GET | `/error_audit_logs` | Error logs |

### Super-Admin Routes (require super-admin role)

| Method | Path | Description |
|---|---|---|
| POST | `/remove_admin` | Remove an admin |
| POST | `/invite-become-admin` | Invite user to become admin |

---

## Environment Variables

Required in `.env` at project root:

```env
# ── Application ───────────────────────────────────────────
BACKEND_DOMAIN=localhost          # or your domain
BACKEND_PORT=5321
FRONTEND_DOMAIN=localhost
FRONTEND_PORT=3000
PORT=5321                         # PaaS health check override
DATABASE_TYPE=postgres
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=resuming

# ── Valkey (Redis-compatible cache) ───────────────────────
VALKEY_DOMAIN=localhost
VALKEY_PORT=6379
VALKEY_PASSWORD=
VALKEY_TLS=false                  # true for managed Valkey/Upstash

# ── Email (Mailgun) ───────────────────────────────────────
EMAIL=noreply@yourdomain.com
MAILGUN_DOMAIN=yourdomain.com
MAILGUN_API_KEY=key-xxx

# ── OAuth (optional / stub) ──────────────────────────────
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
```

---

## Database Migrations

Managed by **Goose**. The initial schema (`migrations/001_initial_schema.sql`) creates:

### Tables

| Table | Purpose |
|---|---|
| `users` | All user accounts (clients, admins, super-admins) |
| `showcase_records` | User profile data (name, email, skills, experience, etc.) |
| `portfolios` | Portfolio content per user + template selection |
| `resumes` | Resume content per user + template selection |
| `ats` | ATS scoring results per user |
| `sessions` | Active session keys (1 per user at a time) |
| `jwt_keys` | JWT signing keys (1 per user at a time) |
| `client_audit_logs` | Client action audit trail |
| `admin_audit_logs` | Admin action audit trail |
| `client_report_logs` | Client-to-client reports |
| `error_logs` | System error logs |
| `client_support_messaging` | Support ticket messages |

### Key Design Decisions

- **Soft deletes**: `deleted_at`, `banned_at`, `expires_at` timestamps instead of hard deletes
- **UUID public IDs**: All tables use `public_id` (UUID) for external references, never expose auto-increment `id`
- **Cascade deletes**: `showcase_records`, `portfolios`, `resumes`, `ats` all cascade when a user is deleted
- **JSONB columns**: `detail` fields in portfolios/resumes store flexible template data

---

## Testing

Tests live in `test/` and cover:

- `auth_test.go` — registration, login, session flows
- `administrator_test.go` — admin operations, audit logs
- `client_support_test.go` — messaging, reporting
- `setting_test.go` — account changes
- `validator_test.go` — input validation edge cases

Run tests:

```bash
cd apps/backend
go test ./test/...
```

---

## Key Patterns

### 1. Context-Based Data Passing

Handlers don't return data directly. Instead, they set values on the Echo context:

```go
// In a view handler:
c.Set("response_data", result)
return nil

// In wrapper.go:
data := c.Get("response_data")
return c.JSON(http.StatusOK, echo.Map{"data": data})
```

### 2. Public IDs Over Internal IDs

All external API operations use `public_id` (UUID). Internal auto-increment `id` is never exposed. The `resolveUserId()` helper translates public → private IDs via Valkey cache.

### 3. Async Database Sync

Setting mutations (username, email, etc.) trigger **async Valkey → PostgreSQL sync** in goroutines so the API response isn't blocked:

```go
go func(psid string) {
    if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
        log.Printf("Failed to sync user data: %v", err)
    }
}(public_user_id.(string))
```

### 4. Valkey Cache-Aside with Fallback

Reads check Valkey first. On cache miss (`valkey.IsValkeyNil`), fall back to PostgreSQL and repopulate Valkey:

```go
data, err := tool.Valkey.Do(ctx, ...).ToString()
if err != nil {
    if valkey.IsValkeyNil(err) {
        user, dbErr := database.FindUserByPublicId(public_user_id)
        // sync to Valkey, retry read
    }
}
```

### 5. Generic Validation

The `showcaserecord` domain uses Go generics for type-safe validation:

```go
validated_request, err := util.ValidateData[typing.NameSection](request)
```

### 6. Casbin for RBAC

Role checks are declarative via Casbin model files, not hardcoded `if` statements. To add a new role, update the `.conf` and `.csv` files.

---

## Contributing

1. Read the [architecture overview](#architecture-overview) and [key patterns](#key-patterns)
2. Follow the existing domain structure (`typing` → `validator` → `util` → `view`)
3. Run `bun x ultracite fix` before committing
4. Ensure migrations are added for any schema changes in `migrations/`
5. Update `sql/queries/*.sql` and regenerate sqlc if queries change
