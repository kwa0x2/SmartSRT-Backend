# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Running the Application
- **Development with hot reload**: `docker compose up --build`
  - Main API server runs on port 9000
  - Consumer service runs separately for background processing
  - Uses Air for hot reloading via `.air.app.toml` and `.air.consumer.toml`

### Building
- **Production build**: `docker build -t myapp .`
- **Cross-platform build**: `docker build --platform=linux/amd64 -t myapp .`

### Services Access
- **API Server**: http://localhost:9000
- **MongoDB Express**: http://localhost:8081 (admin interface)
- **RabbitMQ Management**: http://localhost:15672
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001

## Architecture

### Clean Architecture Pattern
This Go backend follows clean architecture with clear separation of concerns:

- **`cmd/`**: Entry points for main server (`main.go`) and consumer (`consumer/main.go`)
- **`api/`**: HTTP layer with delivery handlers, middleware, and routes
- **`domain/`**: Core business entities and interfaces (repository/usecase contracts)
- **`usecase/`**: Business logic layer implementing domain interfaces
- **`repository/`**: Data access layer with implementations for MongoDB, DynamoDB, S3, etc.
- **`bootstrap/`**: Application initialization and dependency injection

### Key Components

**Dual Service Architecture**:
- **Main API Server**: HTTP REST API using Gin framework
- **Consumer Service**: RabbitMQ worker for asynchronous file processing

**External Dependencies**:
- **MongoDB**: Primary database (replica set with 3 nodes)
- **DynamoDB**: AWS NoSQL database for specific use cases
- **S3**: File storage for audio/video files and generated SRT files
- **Lambda**: AWS Lambda for file conversion processing
- **RabbitMQ**: Message queue for async processing
- **Paddle**: Payment processing SDK

**Monitoring Stack**:
- **Prometheus**: Metrics collection (accessible at `/api/v1/metrics`)
- **Grafana**: Visualization dashboard
- **Sentry**: Error tracking and monitoring

### File Processing Flow
1. Files uploaded via REST API are queued in RabbitMQ
2. Consumer service processes queue messages
3. Files are converted to SRT using AWS Lambda
4. Results stored in S3 and notifications sent via email

### Configuration
- Environment variables managed through `.env` file
- Bootstrap package handles all service initialization
- Structured logging with slog throughout the application

## Code Patterns

### Repository Pattern
All data access follows the generic `BaseRepository[T Entity]` interface defined in `domain/base.go`.

### Error Handling
- Structured error handling with Sentry integration
- Manual error reporting in critical paths (see consumer service)

### Middleware Stack
Standard middleware includes CORS, rate limiting, authentication, and Prometheus metrics.