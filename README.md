# TicketFlow Backend

Production-ready, idiomatic Go backend project structure for **TicketFlow** — an online event ticketing platform designed with layered architecture, strong concurrency handling, and clean separation of concerns.

---

## 1. Project Overview & Architecture

This project adapts the standard Go project layout defined in [golang-standards/project-layout](https://github.com/golang-standards/project-layout) to a pragmatic layered architecture.

```text
ticketflow/
├── .env.example              # Environment variables template
├── .gitignore                 # Git ignore rules for Go, binaries, and local configs
├── Dockerfile                 # Multi-stage production container build (Alpine runtime)
├── docker-compose.yml         # Local container orchestration (API + PostgreSQL 16)
├── Makefile                   # Development automation commands
├── README.md                  # Project documentation and developer guide
├── go.mod                     # Go module definitions
├── go.sum                     # Checksums of dependencies
│
├── api/                       # API contracts and specifications
│   └── openapi.yaml           # OpenAPI 3.0 specification for endpoints
│
├── cmd/                       # Application entrypoints
│   └── api/
│       └── main.go            # API entrypoint: config loading, logger, graceful shutdown
│
├── internal/                  # Private application code (non-exportable)
│   ├── config/                # Strongly typed environment configuration
│   │   ├── config.go          # Config struct, defaults, and validation rules
│   │   └── config_test.go     # Unit tests for configuration
│   │
│   ├── handler/               # Transport layer: HTTP handlers & response serialization
│   │   ├── health.go          # Health probes (/healthz, /livez)
│   │   └── health_test.go     # Unit tests for handlers
│   │
│   ├── middleware/            # HTTP middleware components
│   │   ├── logger.go          # Structured HTTP request logger using log/slog
│   │   ├── recover.go         # Panic recovery middleware with structured logging
│   │   └── middleware_test.go # Unit tests for middlewares
│   │
│   └── server/                # HTTP server setup and lifecycle management
│       ├── server.go          # Server wrapping http.Server, Gin engine, route registry
│       └── server_test.go     # Server lifecycle & route mapping tests
│
├── migrations/                # Database migration scripts (SQL)
│   └── .gitkeep               # Directory placeholder for future migration files
│
├── scripts/                   # Development and operational scripts
│   └── dev.ps1                # Development server runner for PowerShell
│
└── tests/                     # Black-box integration and end-to-end tests
    └── integration/
        └── health_test.go     # End-to-end HTTP integration test suite
```

---

## 2. Layer Responsibilities & Architectural Blueprint

In accordance with Go best practices and pragmatic layered architecture, directories are only added when they serve an active purpose:

- **`cmd/`**: Entrypoints of the application. `cmd/api/main.go` wires dependencies together and coordinates graceful shutdown on `SIGINT` / `SIGTERM`.
- **`internal/config/`**: Reads and validates environment variables with typed values and sensible defaults.
- **`internal/handler/`**: Transport layer. Decodes HTTP requests, triggers services, and encodes HTTP responses.
- **`internal/middleware/`**: Cross-cutting concerns such as structured request logging (`log/slog`) and panic recovery.
- **`internal/server/`**: Encapsulates `http.Server` and `gin.Engine`, with route registration and graceful termination using timeout contexts.
- **`internal/domain/`** *(Next phase)*: Pure business entities, value objects, domain errors, and repository interfaces. Free of HTTP or DB dependencies.
- **`internal/service/`** *(Next phase)*: Business use case orchestration and transaction management (e.g., atomic ticket reservations).
- **`internal/repository/`** *(Next phase)*: PostgreSQL data access implementations using `pgx` and row-level locking (`SELECT ... FOR UPDATE`).
- **`internal/infrastructure/`** *(Next phase)*: Database connection pool lifecycle (`pgxpool`), external payment gateways, etc.
- **`migrations/`**: Versioned SQL migration files.
- **`api/`**: OpenAPI 3.0 documentation.
- **`tests/`**: Black-box integration tests executing real HTTP requests over the network.

---

## 3. Getting Started

### Prerequisites

- **Go**: `1.22+` (installed: `go version go1.27.4 windows/amd64`)
- **Docker & Docker Compose** (optional, for containerized environment)
- **Make** / `mingw32-make` (optional, for Makefile targets)

### Configuration

Create a `.env` file from `.env.example`:

```bash
cp .env.example .env
```

Default settings:
```env
APP_ENV=development
PORT=8080
LOG_LEVEL=debug
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
SERVER_IDLE_TIMEOUT=60s
SERVER_SHUTDOWN_TIMEOUT=15s
```

---

## 4. Development Commands

### Running Locally

Using Go directly:
```bash
go run ./cmd/api
```

Using PowerShell script:
```powershell
.\scripts\dev.ps1
```

Using Make:
```bash
make run
```

### Building Binary

```bash
# Via Make
make build

# Or via Go directly
go build -o bin/ticketflow-api.exe ./cmd/api
```

### Running Tests

```bash
# Run internal unit tests with race detection
go test -v -race ./internal/...

# Run integration tests
go test -v -race ./tests/integration/...

# Run all tests with code coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out
```

Or using Make:
```bash
make test
make test-integration
make test-coverage
```

### Docker Support

Build and run using Docker Compose:
```bash
docker compose up --build
```

---

## 5. Verified Endpoints

Once the application is running on port `8080`:

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/healthz` | System health probe (returns status, app name, environment, timestamp) |
| `GET` | `/livez` | Lightweight process liveness probe (`{"status":"alive"}`) |
| `GET` | `/api/v1/health` | Versioned API health check |

**Example Response (`/healthz`):**
```json
{
  "status": "ok",
  "application": "ticketflow",
  "environment": "development",
  "timestamp": "2026-11-24T15:30:00Z"
}
```
