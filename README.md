# SmartSRT Backend

[![Go 1.23](https://img.shields.io/badge/Go-1.23%2B-00ADD8.svg?logo=go&logoColor=white)](https://go.dev/)
[![MongoDB](https://img.shields.io/badge/MongoDB-Replica%20Set-47A248.svg?logo=mongodb&logoColor=white)](https://www.mongodb.com/)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3-FF6600.svg?logo=rabbitmq&logoColor=white)](https://www.rabbitmq.com/)
[![AWS](https://img.shields.io/badge/AWS-Lambda%20%7C%20S3%20%7C%20DynamoDB-FF9900.svg?logo=amazonaws&logoColor=white)](https://aws.amazon.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> A production-ready, AI-powered subtitle generation backend built with Go. Powers SmartSRT — a service that turns audio and video files into accurate `.srt` subtitles.

---

## Overview

SmartSRT Backend is the API and processing service behind SmartSRT. Users upload audio/video files; the backend queues them in RabbitMQ, dispatches transcription jobs to AWS Lambda, stores the resulting subtitle files in S3 and notifies users by email. It exposes a REST API for uploads, subscription management (Paddle) and account features.

The codebase follows **Clean Architecture** with strict separation between `delivery → usecase → repository → domain`, making it testable, extendable and provider-agnostic.

## Features

- **Audio & video transcription** to `.srt` subtitles via AWS Lambda
- **Asynchronous processing pipeline** powered by RabbitMQ
- **Authentication** with JWT plus Google and GitHub OAuth
- **Subscription & billing** integrated with [Paddle](https://www.paddle.com/)
- **Usage tracking & quotas** per user / subscription tier
- **Transactional emails** through Resend with HTML templates
- **SMS / OTP** via Sinch
- **Observability** — Prometheus metrics + Grafana dashboards + Sentry error tracking
- **Containerized stack** — API, consumer, MongoDB replica set, RabbitMQ, monitoring
- **Hot reload** in development with [Air](https://github.com/air-verse/air)

## Tech Stack

| Layer            | Technology                                                                 |
|------------------|----------------------------------------------------------------------------|
| Language         | Go 1.23+                                                                   |
| HTTP Framework   | [Gin](https://github.com/gin-gonic/gin)                                    |
| Database         | MongoDB (3-node replica set)                                               |
| Secondary Store  | AWS DynamoDB                                                               |
| Object Storage   | AWS S3 (uploads + generated `.srt`)                                        |
| Processing       | AWS Lambda                                                                 |
| Message Queue    | RabbitMQ                                                                   |
| Billing          | Paddle SDK                                                                 |
| Email / SMS      | Resend, Sinch                                                              |
| Auth             | JWT, Google OAuth, GitHub OAuth                                            |
| Config           | Viper                                                                      |
| Validation       | go-playground/validator                                                    |
| Monitoring       | Prometheus, Grafana, Sentry                                                |
| Logging          | `log/slog` (structured)                                                    |
| Containerization | Docker & Docker Compose                                                    |

## Architecture

SmartSRT runs as **two cooperating Go services** behind a shared infrastructure stack:

![Architecture](https://smartsrt.s3.eu-west-3.amazonaws.com/assets/architecture.png)

### Clean Architecture Layout

- **`cmd/`** — entry points (`main.go` for API, `consumer/main.go` for the worker)
- **`api/`** — HTTP delivery, middleware and route registration
- **`domain/`** — core entities and repository / usecase interfaces
- **`usecase/`** — business logic
- **`repository/`** — concrete data access (MongoDB, DynamoDB, S3, …)
- **`bootstrap/`** — dependency wiring and service initialization
- **`config/`** — environment-driven configuration (Viper)
- **`utils/`** — shared helpers

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) & Docker Compose
- AWS account with access to Lambda, S3 and DynamoDB
- Paddle account (for billing endpoints)
- Resend & Sinch accounts (for email / SMS)

### 1. Clone the repository

```bash
git clone https://github.com/kwa0x2/SmartSRT-Backend.git
cd SmartSRT-Backend
```

### 2. Configure environment variables

Copy `.env.example` to `.env` and fill in your credentials:

```bash
cp .env.example .env
```

```env
GIN_MODE=debug
APP_ENV=development
SERVER_ADDRESS=:9000

FRONTEND_URL=http://localhost:3000
JWT_SECRET=

# MongoDB (replica set is started by docker compose)
MONGO_URI=mongodb://user:password@mongo_rs0:27017,mongo_rs1:27018,mongo_rs2:27019/smartsrt?replicaSet=rs0&authSource=admin
MONGO_DB_NAME=smartsrt

# OAuth — Google
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:9000/api/v1/auth/google/callback

# OAuth — GitHub
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
GITHUB_REDIRECT_URL=http://localhost:9000/api/v1/auth/github/callback

# AWS
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_REGION=
AWS_S3_BUCKET_NAME=
AWS_LAMBDA_FUNC_NAME=

# SMS (Sinch) & Email (Resend)
SINCH_APP_KEY=
SINCH_APP_SECRET=
RESEND_API_KEY=
NOTIFY_EMAIL=

# Billing (Paddle)
PADDLE_API_KEY=
PADDLE_WEBHOOK_SECRET_KEY=

# Observability
SENTRY_DSN=

# Quotas (minutes / month)
FREE_MONTHLY_LIMIT=600
PRO_MONTHLY_LIMIT=3000
```

### 3. Start the full stack

```bash
docker compose up --build
```

This spins up:

| Service             | URL                              |
|---------------------|----------------------------------|
| API Server          | http://localhost:9000            |
| MongoDB Express     | http://localhost:8081            |
| RabbitMQ Management | http://localhost:15672           |
| Prometheus          | http://localhost:9090            |
| Grafana             | http://localhost:3001            |

Metrics are exposed at `http://localhost:9000/api/v1/metrics`.

### 4. Cross-platform build

```bash
docker build --platform=linux/amd64 -t smartsrt-backend .
docker build --platform=linux/amd64 -f Dockerfile.consumer -t smartsrt-consumer .
```

## API

All endpoints are prefixed with `/api/v1`. Route groups:

| Group              | Purpose                                |
|--------------------|----------------------------------------|
| `/auth`            | Sign-up, login, OAuth, sessions        |
| `/user`            | Profile management                     |
| `/srt`             | Upload media, list and download `.srt` |
| `/subscription`    | Plans and subscription lifecycle       |
| `/paddle`          | Paddle webhooks                        |
| `/usage`           | Per-user usage and quota               |
| `/contact`         | Contact form submissions               |
| `/metrics`         | Prometheus scrape endpoint             |

## Project Structure

```
SmartSRT-Backend/
├── api/                # HTTP delivery, middleware, routes
│   ├── http/delivery/  # Gin handlers
│   ├── middleware/     # Auth, CORS, rate-limit, metrics
│   └── route/          # Route registration per domain
├── bootstrap/          # App initialization & DI
├── cmd/                # Entry points (API + consumer)
├── config/             # Viper-driven configuration
├── domain/             # Entities and contracts
├── email_templates/    # HTML email templates
├── monitoring/         # Prometheus & Grafana configs
├── rabbitmq/           # Queue setup & helpers
├── repository/         # MongoDB / DynamoDB / S3 implementations
├── seeder/             # DB seeders for local development
├── usecase/            # Business logic
├── utils/              # Helpers
├── compose.yaml        # Full Docker Compose stack
├── Dockerfile          # API image
└── Dockerfile.consumer # Consumer image
```

## Related Repositories

- [SmartSRT Frontend](https://github.com/kwa0x2/SmartSRT-Frontend) — web client
- [SmartSRT Lambda](https://github.com/kwa0x2/SmartSRT-Lambda) — transcription Lambda function

## License

Distributed under the MIT License. See [`LICENSE`](LICENSE) for more information.
