# Trossage Backend

Real-time messenger backend with WebSocket support.

## Features

- Registration and authentication (JWT access/refresh tokens)
- Private chats between users
- Real-time message and event delivery (WebSocket)
- Typing indicator
- User search

## Tech Stack

- **Go 1.25**, Gin, pgx
- **WebSocket** (coder/websocket)
- **JWT** for authentication
- **Argon2id** for password hashing
- **PostgreSQL 16**

## Getting Started

### Install Tools

```bash
# Install Homebrew tools
brew install golangci-lint  # Go linter
brew install sqlfluff       # SQL style linter  
brew install rust           # Rust toolchain (required for squawk)

# Install Rust-based tools
cargo install squawk-cli    # PostgreSQL migration safety linter
```

Go tools (swag, mockery) are installed automatically via `go tool`.

### Run

```bash
# Copy and configure environment variables
cp .env.example .env

# Start with Docker Compose
make docker-up
```

API available at `http://localhost:8080`, health-check at `http://localhost:8081/health`.

## API

Swagger documentation: [`docs/swagger.json`](docs/swagger.json).

### Main Endpoints

| Method | Path                      | Description          |
|--------|---------------------------|----------------------|
| POST   | `/api/auth/register`      | Register             |
| POST   | `/api/auth/login`         | Login                |
| POST   | `/api/auth/refresh`       | Refresh tokens       |
| POST   | `/api/auth/logout`        | Logout               |
| GET    | `/api/users/me`           | Current user         |
| GET    | `/api/users/search`       | Search users         |
| POST   | `/api/chats`              | Create chat          |
| GET    | `/api/chats`              | List chats           |
| POST   | `/api/chats/:id/messages` | Send message         |
| GET    | `/api/chats/:id/messages` | Message history      |
| GET    | `/api/ws`                 | WebSocket connection |

## Project Structure

```
cmd/                 # Entry point
internal/
  application/       # Application lifecycle
  config/            # Configuration
  http/              # HTTP handlers, middleware, DTO
  repository/        # Repository layer
  service/           # Business logic
  websocket/         # WebSocket hub and clients
  worker/            # Background tasks
migrations/          # SQL migrations
docs/                # Swagger documentation
```

## Make Commands

| Command                       | Description                                |
|-------------------------------|--------------------------------------------|
| `make gen`                    | Generate code (swagger, mocks)             |
| `make test`                   | Run tests                                  |
| `make lint-all`               | Run all linters                            |
| `make lint-file FILE=path.go` | Lint single file                           |
| `make docker-up`              | Rebuild and start containers               |
| `make docker-down`            | Stop containers and clean images           |
| `make docker-deep-clean`      | Full Docker cleanup (images, cache, vols)  |

