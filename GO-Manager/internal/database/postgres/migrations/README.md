# PostgreSQL Migrations

This package provides a simple migration runner for the existing PostgreSQL schema.

## Usage

1. Import package:

```go
import "VSRT-Lang/internal/database/postgres/migrations"
```

2. Run migrations before using repositories:

```go
if err := migrations.Run(db); err != nil {
    panic(err)
}
```

3. Add new migrations to `internal/database/postgres/migrations/migrations.go`.
