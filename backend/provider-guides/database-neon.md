# database:neon starter

The starter uses GORM with the PostgreSQL driver. The database layer is at `internal/db`, which keeps GORM details out of HTTP handlers. `internal/db/schema.go` is the single registry of managed schema models.

- `DATABASE_URL`: pooled Neon URL for API traffic; the host includes `-pooler`.
- `DATABASE_URL_UNPOOLED`: direct Neon URL for `go run ./cmd/migrate`; the host has no `-pooler` suffix.

Run migrations with `go run ./cmd/migrate` from `backend`. The API's repository and migration commands use the same GORM domain model, so the schema is maintained in versioned application code rather than an ad hoc script.

Official guidance: https://neon.com/docs/connect/connection-pooling
