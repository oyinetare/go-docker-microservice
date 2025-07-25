# go-docker-microservice

Toy microservice project inspired by the article ["Learn Docker by Building a Microservice"](https://dwmkerr.com/learn-docker-by-building-a-microservice/) by Dave Kerr that manages a directory of email addresses to phone numbers and built with Go and Docker.

This project demonstrates modern microservice architecture patterns, containerization with Docker, and best practices for building scalable applications.

## 📋 Table of Contents

- [Features](#-features)
- [Prerequisites](#-prerequisites)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [API Documentation](#-api-documentation)
- [Testing](#-testing)
- [Development](#-development)
- [Docker Commands](#-docker-commands)
- [Architecture & Design Patterns](#-architecture--design-patterns)
- [Note](#-note)
- [Acknowledgments](#-acknowledgments)

## ✨ Features

- RESTful API built with Go for managing user directory
- MySQL database with Docker containerization
- Multi-stage Docker builds for minimal image size
- Docker Compose for orchestration
- Comprehensive unit and integration tests built with testify and go-sqlmock
- Environment-based configuration
- Layered architecture with clear separation of concerns
- Repository pattern for data access
- Dependency injection for testability
- Gorilla Mux for routing

## 📦 Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.21+ (for local development)

## 📁 Project Structure

```
go-docker-microservice/
├── integration-tests/  # End-to-end tests
├── test-database/      # MySQL Docker setup
├── users-service/      # Go microservice
├── docker-compose.yml  # Service orchestration
└── README.md           # Project documentation
```

## 🚀 Getting Started

### Quick Start with Docker Compose

1. Clone the repository:

```bash
git clone https://github.com/oyinetare/go-docker-microservice.git
cd go-docker-microservice
```

2. Build and run with Docker Compose to start all services:

```bash
docker-compose up --build
```

_you might need to run_ `docker-compose up` _again to start the user-service. This happens because of race condition issues._

3. The API will be available at `http://localhost:8123`

4. Test the endpoints:

```bash
# Get all users
curl http://localhost:8123/users

# Search for a specific user
curl http://localhost:8123/search?email=bart@thesimpsons.com
```

5. Stop the services:

```bash
docker-compose down
```

### Local Development

1. Start the database:

```bash
docker-compose up db
```

2. Install dependencies:

```bash
cd users-service
go mod download
```

3. Run the service:

```bash
go run main.go
```

4. Format your code:

```bash
# Format current directory
go fmt ./...

# Or use gofmt directly
gofmt -w .
```

5. Lint your code (optional):

```bash
# Install golangci-lint (recommended)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Or use go vet (built-in)
go vet ./...
```

## 📖 API Documentation

### Endpoints

#### Get All Users

```http
GET /users
```

**Response:**

```json
[
  { "email": "homer@thesimpsons.com", "phoneNumber": "+1 888 123 1111" },
  { "email": "marge@thesimpsons.com", "phoneNumber": "+1 888 123 1112" },
  { "email": "maggie@thesimpsons.com", "phoneNumber": "+1 888 123 1113" },
  { "email": "lisa@thesimpsons.com", "phoneNumber": "+1 888 123 1114" },
  { "email": "bart@thesimpsons.com", "phoneNumber": "+1 888 123 1115" }
]
```

#### Search User by Email

```http
GET /search?email=bart@thesimpsons.com
```

**Response (200 OK):**

```json
{ "email": "bart@thesimpsons.com", "phoneNumber": "+1 888 123 1115" }
```

**Response (404 Not Found):**

```
User not found
```

### Status Codes

- `200 OK` - Successful request
- `400 Bad Request` - Missing or invalid parameters
- `404 Not Found` - User not found
- `500 Internal Server Error` - Server error

## 🧪 Testing

### Run Unit Tests

```bash
cd users-service

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run with verbose output
go test -v ./...

# Run with race detection
go test -race ./...
```

### Run Integration Tests

```bash
# Start services
docker-compose up -d

# Run integration tests
cd integration-tests
go test -v

# Cleanup
docker-compose down
```

### Test Coverage Goals

- Config package: ~100%
- Repository package: ~85%
- API package: ~90%
- Server package: ~85%

## 💻 Development

### Environment Variables

| Variable        | Description  | Default     |
| --------------- | ------------ | ----------- |
| `PORT`          | Service port | `8123`      |
| `DATABASE_HOST` | MySQL host   | `127.0.0.1` |

### Building the Docker Image

```bash
# Build users service
docker build -t users-service ./users-service

# Build database
docker build -t test-database ./test-database
```

### Running Individual Containers

```bash
# Run database
docker run --name db -p 3306:3306 test-database

# Run users service (linked to database)
docker run --name users -p 8123:8123 --link db:db -e DATABASE_HOST=db users-service
```

## 🐳 Docker Commands

### Docker Compose Commands

```bash
# Start all services
docker-compose up

# Start in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Rebuild and start
docker-compose up --build

# Scale service
docker-compose up --scale users-service=3
```

### Useful Docker Commands

```bash
# List running containers
docker ps

# View container logs
docker logs <container_name>

# Execute command in container
docker exec -it <container_name> /bin/sh

# Inspect container
docker inspect <container_name>

# View resource usage
docker stats
```

## 🏗 Architecture & Design Patterns

```
┌─────────────┐     ┌─────────────┐
│   Client    │────▶│Users Service│
└─────────────┘     │   (Go API)  │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │    MySQL    │
                    │   Database  │
                    └─────────────┘
```

### Layered Architecture

```
API Layer (HTTP Handlers)
    ↓
Business Logic Layer (Server)
    ↓
Data Access Layer (Repository)
    ↓
Database Layer (MySQL)
```

### Design Patterns

This project implements several design patterns:

1. **Repository Pattern** - Abstracts database access
2. **Dependency Injection** - Improves testability and flexibility
3. **Factory Pattern** - Creates configuration objects
4. **Layered Architecture** - Separates concerns
5. **Multi-stage Docker Builds** - Optimizes image size

## 📝 Note

This is an educational project designed for learning Docker and microservices concepts. For production use, consider adding:

- Authentication and authorization
- Input validation and sanitization
- Structured logging and monitoring
- Health checks and circuit breakers
- Database connection pooling
- Proper error handling and recovery

## 🙏 Acknowledgments

- Inspired by the article ["Learn Docker by Building a Microservice" by Dave Kerr](https://dwmkerr.com/learn-docker-by-building-a-microservice/)
- [Microservices architecture style by Microsoft](https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/microservices)
- [Design Patterns Reference - refactoring.guru](https://refactoring.guru/design-patterns) - For understanding the design patterns used in this project
