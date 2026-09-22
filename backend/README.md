# Go API

This is a Go 1.22 API backed by Neon Lakebase Postgres through GORM.

## Local startup

1. Copy the root `.env.example` to `.env` and set the required values.
2. Set `DATABASE_URL` to Neon's pooled URL (host contains `-pooler`) for the API.
3. Set `DATABASE_URL_UNPOOLED` to Neon's direct URL for migrations.
4. From this directory, run `go mod tidy`, `go run ./cmd/migrate`, and then `go run ./cmd/api`.

The migration command uses the direct connection and GORM's `AutoMigrate` to create or update the `users` table. The API only uses the pooled connection. `PreferSimpleProtocol` is enabled to work safely with PgBouncer transaction pooling.

For the included local PostgreSQL container, run `docker compose up -d postgres` from the repository root, then set both database variables to `postgres://app:app@localhost:5432/app?sslmode=disable` before running the migration and API commands. In Neon, keep the two URLs distinct as described above.

`internal/` keeps configuration, domain models, PostgreSQL repositories, authentication, and HTTP routes separate. The GORM schema registry is [internal/db/schema.go](internal/db/schema.go); add future persistent models there. The GCS adapter remains isolated until file endpoints are added.
