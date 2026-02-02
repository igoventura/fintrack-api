# Project Structure

This project follows **Clean Architecture** principles to ensure isolation of business logic from external concerns (HTTP, Database, Frameworks).

## Directory Tree

```text
.
├── cmd/
│   └── api/                # Application entry point
│       └── main.go         # Wire up dependencies and start the server
├── domain/                 # (Core) Business entities and repository interfaces
│   ├── account.go          # Struct definitions + Repository Interface
│   ├── user.go
│   └── ...
├── internal/
│   ├── api/                # Transport Layer (Adapters)
│   │   ├── handler/        # HTTP Handlers (controllers)
│   │   ├── middleware/     # Auth, Logging, CORS
│   │   ├── router/         # Route definitions and Redoc registration
│   │   └── dto/            # Data Transfer Objects (Request/Response structs)
│   ├── service/            # Use Cases (Business Logic)
│   │   ├── account_service.go # Orchestrates entities and repos
│   │   └── user_service.go
│   ├── db/                 # Persistence Layer (Adapters)
│   │   └── postgres/       # SQL implementation using pgx
│   └── config/             # Configuration loading (env vars, .yaml)
├── docs/                   # Documentation
│   └── openapi.yaml        # OpenAPI 3.0 specification
├── migrations/             # Database migrations
└── Makefile                # Automation commands
```

---

## Architectural Layers

### 1. Domain Layer (`/domain`)
The center of the application. It contains entities and repository interfaces. 
- **Rule**: This layer must **never** import anything from `internal`.

### 2. Service Layer (`/internal/service`)
Contains the business logic (Use Cases). It acts as an orchestrator between the API layer and the Domain.
- **Dependency**: Depends only on the interfaces defined in `/domain`.

### 3. API Layer (`/internal/api`)
Handles all HTTP-specific logic. 
- **Handlers**: Parse requests, call services, and return responses.
- **DTOs**: Ensure we don't leak internal domain details (like password hashes) to the outside world.
- **Redoc**: Served via the router layer using `go-redoc` to render the `docs/openapi.yaml`.

### 4. Database Layer (`/internal/db`)
Handles the actual persistence.
- **Implementation**: Implements the repository interfaces defined in `/domain`.

---

## Dependency Flow
Dependency should always point **inwards** towards the domain.

`Handler (API)` -> `Service (Business Logic)` -> `Repository (Domain Interface)` <- `Postgres (DB Implementation)`

---

## Documentation (OpenAPI)
Documentation is driven by the `docs/openapi.yaml` file.
- The API serves this file via the `/docs` endpoint using **Redoc**.
- Any change in the DTOs or endpoints should be manually updated in the spec file to keep documentation in sync.
