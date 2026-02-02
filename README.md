# FinTrack Core

FinTrack Core is a Go library providing the foundational domain models and persistence logic for the FinTrack ecosystem. It is designed with modularity and testability in mind, following a domain-driven approach.

## Structure

See [PROJECT_STRUCTURE.md](file:///Users/igoventura/Developer/Personal/fintrack-core/PROJECT_STRUCTURE.md) for details on the codebase organization and architectural layers.

## Infrastructure

The project includes a full Docker setup for development.

```bash
# Start Postgres and automatically run migrations
make compose
```

## Automation

A `Makefile` is provided for common development tasks. It automatically loads variables from a `.env` file if it exists.

```bash
# Infrastructure
make compose        # Start DB and run migrations via Docker

# Migrations (Local)
make migrate        # Run migrations using local Go tools
make rollback       # Rollback the last migration
make new-migration name=xxx  # Create a new migration file

# Maintenance
make tidy           # Run go mod tidy
make test           # Run all tests
```

## Migrations

This project uses [tern](https://github.com/jackc/tern) for database migrations. Configuration is handled via `migrations/tern.conf` which pulls from environment variables.

## Getting Started

### Prerequisites

- Go 1.25.6 or higher
- PostgreSQL (for local development/runtime)

### Installation

```bash
go get github.com/igoventura/fintrack-core
```

### Usage

Initialize the PostgreSQL connection and the repository using environment variables (standard `.env` file supported via `godotenv`):

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/igoventura/fintrack-core/internal/db/postgres"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	ctx := context.Background()
	
	// Initialize Connection Pool
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := postgres.NewDB(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Repository
	repo := postgres.NewAccountRepository(db)

	// Use the repository...
	tenantID := "some-tenant-uuid"
	accounts, err := repo.List(ctx, tenantID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d accounts", len(accounts))
}
```

## Soft Delete Policy

FinTrack Core implements a comprehensive soft delete strategy to maintain data auditability and prevent accidental data loss.

- **Traceable Deletion**: Critical entities (`Accounts`, `Transactions`, `Attachments`) support full traceability by storing both the timestamp (`deactivated_at`) and the ID of the user who performed the action (`deactivated_by`).
- **Standard Soft Delete**: Operational entities (`Users`, `Tenants`, `Tags`, `Categories`) use a standard `deactivated_at` timestamp.
- **Join Table Policy**: Many-to-many associations like `Users <-> Tenants` are soft-deleted via timestamp, while lightweight associations like `Transactions <-> Tags` are hard-deleted if they lack specific audit requirements in the schema.

### Environment Variables

Create a `.env` file in your root directory:

```env
DATABASE_URL=postgres://user:password@localhost:5432/fintrack?sslmode=disable
DB_HOST=localhost
DB_PORT=5432
DB_NAME=fintrack
DB_USER=postgres
DB_PASSWORD=postgres
```

## Testing

The project uses `pgxmock` for unit testing the repository layer without requiring a live database.

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.