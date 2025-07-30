<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/GOeda-Co/backend">
    <img src="Logo.png" alt="Logo" width="80" height="80">
  </a>
  <h3 align="center">Repeatro – Anki-Style Vocabulary Learning App</h3>
  <p align="center">
    A modern web-based vocabulary learning tool inspired by Anki, built with Go and PostgreSQL.<br />
    <a href="https://github.com/GOeda-Co/frontend"><strong>Frontend repo »</strong></a>
    <br />
    <a href="https://github.com/GOeda-Co/backend/issues">Issues</a>
    &middot;
    <a href="https://github.com/GOeda-Co/backend/pulls">Pull Requests</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#features">Features</a></li>
    <li><a href="#built-with">Built With</a></li>
    <li><a href="#project-architecture-overview">Project Architecture Overview</a></li>
    <li><a href="#database-architecture-overview">Database Architecture Overview</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
    <li><a href="#contributors">Contributors</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

Repeatro is a modern, open-source vocabulary learning app inspired by Anki. It leverages spaced repetition (SM2 algorithm) to help users efficiently retain vocabulary. Organize your words into decks, track your progress, and enjoy a simple, effective learning experience.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Features

- Spaced repetition for efficient vocabulary retention (SM2 algorithm)
- JWT-based user authentication
- Decks to organize vocabulary by topic or language
- RESTful API with [Swaggo][swaggo] auto-generated Swagger docs

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Built With

- [Go](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [Swaggo][swaggo] (API docs)
- [lingua-go][lingua-go] (language detection)
- [JWT](https://jwt.io/) (authentication)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Project Architecture Overview

Repeatro follows a **microservices architecture** pattern, designed for scalability, maintainability, and clear separation of concerns. The system is composed of 5 main services that communicate via gRPC for internal operations and expose a unified REST API through an API gateway.

### Service Structure

```
┌─────────────────┐    HTTP/REST     ┌─────────────────┐
│   Frontend      │ ◄────────────── │  Repeatro       │
│   (Flutter)   │                 │  (API Gateway)  │
└─────────────────┘                 └─────────────────┘
                                             │
                                        gRPC │
                    ┌────────────────────────┼────────────────────────┐
                    │                        │                        │
                    ▼                        ▼                        ▼
            ┌─────────────┐         ┌─────────────┐         ┌─────────────┐
            │    SSO      │         │    Card     │         │    Deck     │
            │  Service    │         │  Service    │         │  Service    │
            │(Auth/Users) │         │(Vocabulary) │         │(Collections)│
            └─────────────┘         └─────────────┘         └─────────────┘
                    │                        │                        │
                    │                        │                        │
                    └────────────────────────┼────────────────────────┘
                                             │
                                             ▼
                                    ┌─────────────┐
                                    │    Stats    │
                                    │  Service    │
                                    │ (Analytics) │
                                    └─────────────┘
```

### Service Responsibilities

- **Repeatro (API Gateway)**: HTTP REST endpoints, request routing, response aggregation, Swagger documentation
- **SSO Service**: User authentication, JWT token management, authorization, admin role management
- **Card Service**: Individual vocabulary card CRUD, SM2 spaced repetition algorithm, card expiration logic
- **Deck Service**: Card collections management, deck organization, card-to-deck relationships
- **Stats Service**: Learning analytics, progress tracking, performance metrics, review history

### Communication Patterns

- **External Communication**: REST API with JSON payloads, CORS-enabled for web clients
- **Internal Communication**: gRPC with Protocol Buffers, type-safe service contracts
- **Authentication**: JWT tokens passed through gRPC metadata for service-to-service auth
- **Service Discovery**: Direct addressing with configurable endpoints (Consul integration planned)

### Technology Stack

Each service follows **Clean Architecture** principles with consistent layering:
- **Presentation Layer**: gRPC controllers, HTTP handlers
- **Business Layer**: Domain services, business logic
- **Data Layer**: GORM repositories, PostgreSQL integration
- **Cross-cutting**: Logging (slog), security, configuration management

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Database Architecture Overview

Repeatro implements a **database-per-service** pattern where each microservice owns its data domain while maintaining logical relationships through application-level coordination.

### Database Design Principles

- **Service Autonomy**: Each service manages its own PostgreSQL schema
- **Eventual Consistency**: Cross-service data consistency through event-driven updates
- **ACID Compliance**: Local transactions within service boundaries
- **UUID Primary Keys**: Distributed system-friendly identifiers

### Data Models & Relationships

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   SSO Service   │    │  Card Service   │    │  Deck Service   │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │    User     │ │    │ │    Card     │ │    │ │    Deck     │ │
│ │ + ID (PK)   │ │    │ │ + CardId    │ │    │ │ + DeckId    │ │
│ │ + Email     │ │    │ │ + CreatedBy │ │    │ │ + CreatedBy │ │
│ │ + PassHash  │ │    │ │ + Word      │ │    │ │ + Name      │ │
│ │ + IsAdmin   │ │    │ │ + Translation│ │   │ │ + Cards[]   │ │
│ └─────────────┘ │    │ │ + Easiness  │ │    │ └─────────────┘ │
│                 │    │ │ + Interval  │ │    │                 │
│ ┌─────────────┐ │    │ │ + DeckID    │ │    └─────────────────┘
│ │    App      │ │    │ │ + ExpiresAt │ │              │
│ │ + ID (PK)   │ │    │ └─────────────┘ │              │
│ │ + Name      │ │    │                 │              │
│ │ + Secret    │ │    └─────────────────┘              │
│ └─────────────┘ │              │                      │
└─────────────────┘              │                      │
                                 │                      │
                                 ▼                      ▼
                        ┌─────────────────┐    ┌─────────────────┐
                        │  Stats Service  │    │ Cross-Service   │
                        │                 │    │ Relationships   │
                        │ ┌─────────────┐ │    │                 │
                        │ │   Review    │ │    │ User 1:N Card   │
                        │ │ + ResultId  │ │    │ User 1:N Deck   │
                        │ │ + UserID    │ │    │ Deck 1:N Card   │
                        │ │ + CardID    │ │    │ Card 1:N Review │
                        │ │ + DeckId    │ │    │                 │
                        │ │ + Grade     │ │    └─────────────────┘
                        │ │ + CreatedAt │ │
                        │ └─────────────┘ │
                        └─────────────────┘
```

### Key Data Features

**Spaced Repetition (SM2 Algorithm)**:
- `Easiness`: Difficulty factor (default 2.5)
- `Interval`: Days until next review
- `RepetitionNumber`: Count of successful reviews
- `ExpiresAt`: Automatic scheduling timestamp

**Multi-tenancy & Security**:
- UUID-based user identification across services
- Cascade deletion for data cleanup
- Row-level security through user ownership

**Performance Optimizations**:
- Indexed foreign keys (`DeckID`, `CreatedBy`)
- Connection pooling (5-20 connections per service)
- PostgreSQL array support for card tags
- Time-based queries for expired cards

### Data Consistency Strategy

**Cross-Service Data Flow**:
1. **User Operations**: SSO → Card/Deck (user validation via gRPC)
2. **Learning Flow**: Card → Stats (review recording)
3. **Deck Management**: Deck ↔ Card (bidirectional updates)

**Consistency Mechanisms**:
- **Synchronous**: Real-time user validation, immediate feedback
- **Asynchronous**: Statistics aggregation, background processing
- **Compensating Actions**: Rollback strategies for cross-service failures

Each service uses GORM's `AutoMigrate` for schema management, with a planned migration to Goose for production-grade database versioning.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Getting Started

This section guides you through setting up Repeatro using Docker for quick deployment and testing.

### Prerequisites

Make sure you have the following installed on your system:

- **Docker 27.3+** - [Download & Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose** - Usually included with Docker Desktop
- **Git** - [Download & Install Git](https://git-scm.com/downloads)

> **Note**: Go and PostgreSQL are not required for Docker setup as they run inside containers.

### Quick Start with Docker

#### 1. Clone the Repository

```bash
git clone https://github.com/GOeda-Co/backend.git
cd backend
```

#### 2. Environment Configuration

Create environment files for all services:

```bash
# Copy the example environment file
cp .env.example .env
```

Edit the `.env` file with your preferred settings:

```properties
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=repeatro

# JWT Secret (generate a secure random string for production)
SECRET=your-super-secret-jwt-key-change-this-in-production

# Service Configuration
CARD_HOST_PORT=50051
CARD_CONTAINER_PORT=50051

DECK_HOST_PORT=50054
DECK_CONTAINER_PORT=50054

REPEATRO_HOST_PORT=8080
REPEATRO_CONTAINER_PORT=8080

SSO_HOST_PORT=44044
SSO_CONTAINER_PORT=44044

STAT_HOST_PORT=50055
STAT_CONTAINER_PORT=50055
```

> **Security Note**: Change the `SECRET` value to a strong, unique string for production deployments.

#### 3. Start All Services

Build and start all microservices with Docker Compose:

```bash
# Build and start all services in detached mode
docker-compose up --build -d

# View logs (optional)
docker-compose logs -f
```

This command will:
- Build Docker images for all microservices
- Start PostgreSQL database
- Launch all services (SSO, Card, Deck, Stats, Repeatro Gateway)
- Set up networking between containers

#### 4. Initialize the Application

Add the application entry to the database (temporary setup step):

```bash
# Connect to the PostgreSQL container
docker exec -it postgres psql -U postgres -d repeatro

# Add the application record
INSERT INTO apps (id, name, secret) VALUES (1, 'repeatro', 'your-super-secret-jwt-key-change-this-in-production');

# Exit PostgreSQL
\q
```

> **Important**: Replace `your-super-secret-jwt-key-change-this-in-production` with the same `SECRET` value from your `.env` file.

#### 5. Verify Installation

Check that all services are running:

```bash
# View running containers
docker-compose ps

# Check service health
curl http://localhost:8080/swagger/index.html
```

You should see:
- All services in "Up" status
- Swagger documentation accessible at `http://localhost:8080/swagger/index.html`

### Service Endpoints

Once running, the following endpoints will be available:

- **Repeatro Gateway (REST API)**: `http://localhost:8080`
- **Swagger Documentation**: `http://localhost:8080/swagger/index.html`
- **SSO Service (gRPC)**: `localhost:44044`
- **Card Service (gRPC)**: `localhost:50051`
- **Deck Service (gRPC)**: `localhost:50054`
- **Stats Service (gRPC)**: `localhost:50055`
- **PostgreSQL Database**: `localhost:5432`

### Common Commands

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (clears database)
docker-compose down -v

# View logs for specific service
docker-compose logs -f repeatro

# Rebuild specific service
docker-compose up --build repeatro

# Access database directly
docker exec -it postgres psql -U postgres -d repeatro
```



<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Development

This section describes how to set up and run the Repeatro project locally for development without Docker.

### Prerequisites

- **Go 1.24+** - [Download & Install Go](https://golang.org/dl/)
- **PostgreSQL 15+** - [Download & Install PostgreSQL](https://www.postgresql.org/download/)
- **Git** - [Download & Install Git](https://git-scm.com/downloads)

### Local Development Setup

#### 1. Clone the Repository

```bash
git clone https://github.com/GOeda-Co/backend.git
cd backend
```

#### 2. Database Setup

Create a PostgreSQL database for the project:

```sql
-- Connect to PostgreSQL as superuser
psql -U postgres

-- Create database
CREATE DATABASE repeatro;

-- Create user (optional, or use existing postgres user)
CREATE USER tomatocoder WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE repeatro TO tomatocoder;
```

#### 3. Environment Configuration

Create a `.env` file in the root directory:

```bash
# Copy example environment file
cp .env.example .env
```

Edit the `.env` file with your local settings:

```properties
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=tomatocoder
DB_PASS=postgres
DB_NAME=repeatro

# JWT Secret (generate a secure random string)
SECRET=your-super-secret-jwt-key-here

# Config Path - use local.yaml for development
CONFIG_PATH=./config/local.yaml

# Service Ports for Local Development
CARD_HOST_PORT=50051
CARD_CONTAINER_PORT=50051

DECK_HOST_PORT=50054
DECK_CONTAINER_PORT=50054

REPEATRO_HOST_PORT=8080
REPEATRO_CONTAINER_PORT=8080

SSO_HOST_PORT=44044
SSO_CONTAINER_PORT=44044

STAT_HOST_PORT=50055
STAT_CONTAINER_PORT=50055
```

#### 4. Local Configuration Files

Each service needs a `local.yaml` configuration file. Create them in each service's config directory:

**SSO Service** (`sso/config/local.yaml`):
```yaml
env: local

connection_string: "host=localhost port=5432 user=tomatocoder password=postgres dbname=repeatro sslmode=disable"

grpc:
  port: 44044
  address: ":44044"
  timeout: 10s

token_ttl: 15m
secret: ${SECRET}
```

**Card Service** (`card/config/local.yaml`):
```yaml
env: local

connection_string: "host=localhost port=5432 user=tomatocoder password=postgres dbname=repeatro sslmode=disable"

grpc:
  port: 50051
  address: ":50051"
  timeout: 10s

clients:
  sso:
    address: ":44044"
    timeout: 5s
    retries_count: 3
  stat:
    address: ":50055"
    timeout: 5s
    retries_count: 3

secret: ${SECRET}
```

**Deck Service** (`deck/config/local.yaml`):
```yaml
env: local

connection_string: "host=localhost port=5432 user=tomatocoder password=postgres dbname=repeatro sslmode=disable"

grpc:
  port: 50054
  address: ":50054"
  timeout: 10s

clients:
  sso:
    address: ":44044"
    timeout: 5s
    retries_count: 3

secret: ${SECRET}
```

**Stats Service** (`stats/config/local.yaml`):
```yaml
env: local

connection_string: "host=localhost port=5432 user=tomatocoder password=postgres dbname=repeatro sslmode=disable"

grpc:
  port: 50055
  address: ":50055"
  timeout: 10s

clients:
  sso:
    address: ":44044"
    timeout: 5s
    retries_count: 3

secret: ${SECRET}
```

**Repeatro Gateway** (`repeatro/config/local.yaml`):
```yaml
env: local

connection_string: "host=localhost port=5432 user=tomatocoder password=postgres dbname=repeatro sslmode=disable"

clients:
  card:
    address: ":50051"
    timeout: 5s
    retries_count: 3
  deck:
    address: ":50054"
    timeout: 5s
    retries_count: 3
  sso:
    address: ":44044"
    timeout: 5s
    retries_count: 3
  stat:
    address: ":50055"
    timeout: 5s
    retries_count: 3

grpc:
  port: 50054
  address: ":50054"
  timeout: 10s

http_server:
  address: "0.0.0.0:8080"
  port: 8080
  timeout: 4s
  idle_timeout: 30s

secret: ${SECRET}
```

#### 5. Install Dependencies

Install Go dependencies for each service:

```bash
# SSO Service
cd sso && go mod tidy && cd ..

# Card Service
cd card && go mod tidy && cd ..

# Deck Service
cd deck && go mod tidy && cd ..

# Stats Service
cd stats && go mod tidy && cd ..

# Repeatro Gateway
cd repeatro && go mod tidy && cd ..
```

#### 6. Running Services

Start each service in separate terminal windows/tabs in the following order:

**Terminal 1 - SSO Service:**
```bash
cd sso
CONFIG_PATH=./config/local.yaml go run cmd/sso/main.go
```

**Terminal 2 - Stats Service:**
```bash
cd stats
CONFIG_PATH=./config/local.yaml go run cmd/stats/main.go
```

**Terminal 3 - Card Service:**
```bash
cd card
CONFIG_PATH=./config/local.yaml go run cmd/card/main.go
```

**Terminal 4 - Deck Service:**
```bash
cd deck
CONFIG_PATH=./config/local.yaml go run cmd/deck/main.go
```

**Terminal 5 - Repeatro Gateway:**
```bash
cd repeatro
CONFIG_PATH=./config/local.yaml go run cmd/repeatro/main.go
```

#### 8. Verify Setup

Once all services are running, you can:

1. **Check service health** by accessing individual gRPC ports
2. **Access Swagger documentation** at `http://localhost:8080/swagger/index.html`
3. **Test API endpoints** using the Swagger UI or curl commands

### Development Tips

- **Hot Reload**: Use tools like [air](https://github.com/cosmtrek/air) for automatic reloading during development
- **Database Migrations**: GORM auto-migration is enabled, so tables will be created automatically
- **Logging**: Set `env: local` in config files for detailed debug logging
- **Service Discovery**: In local development, services communicate via `localhost:port`
- **Config Changes**: Restart services after modifying `local.yaml` files

### Troubleshooting

**Port Already in Use:**
```bash
# Find process using port
lsof -i :8080
# Kill process
kill -9 <PID>
```

**Database Connection Issues:**
- Verify PostgreSQL is running: `brew services start postgresql` (macOS) or `sudo systemctl start postgresql` (Linux)
- Check connection string in `local.yaml` files
- Ensure database `repeatro` exists

**JWT Signature Issues:**
- Ensure all services use the same `SECRET` value in their config files

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

Once running, you can interact with the API (REST/gRPC) for deck and card management, user authentication, and spaced repetition review. See the [Swagger docs](backend/repeatro/docs/swagger.yaml) for API details.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Roadmap

- [ ] Goose migrations
- [ ] Import/export via CSV, JSON
- [ ] Enhance current stats
- [ ] Language detection
- [ ] Simple frontend

See the [open issues][issues-url] for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contact

Maintainers: <br>
[@tomatoCoderq](https://github.com/tomatoCoderq) <br>
[@constable](https://github.com/constable) <br>
[@Kaghorz](https://github.com/tomatoCoderq)<br>

Project Link: [https://github.com/GOeda-Co/backend](https://github.com/GOeda-Co/backend)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Acknowledgments

- [Anki](https://apps.ankiweb.net/) – inspiration
- [Goose][goose]
- [Swaggo][swaggo]
- [lingua-go][lingua-go]
- [Img Shields](https://shields.io)
- [Best-README-Template](https://github.com/othneildrew/Best-README-Template)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contributors

<a href="https://github.com/GOeda-Co/backend/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=GOeda-Co/backend" alt="Contributors" />
</a>

**Project contributors:**
- [@tomatoCoderq](https://github.com/tomatoCoderq)
- [@constabIe](https://github.com/constabIe)
- [@Kaghorz](https://github.com/Kaghorz)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
[contributors-shield]: https://img.shields.io/github/contributors/GOeda-Co/backend.svg?style=for-the-badge
[contributors-url]: https://github.com/GOeda-Co/backend/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/GOeda-Co/backend.svg?style=for-the-badge
[forks-url]: https://github.com/GOeda-Co/backend/network/members
[stars-shield]: https://img.shields.io/github/stars/GOeda-Co/backend.svg?style=for-the-badge
[stars-url]: https://github.com/GOeda-Co/backend/stargazers
[issues-shield]: https://img.shields.io/github/issues/GOeda-Co/backend.svg?style=for-the-badge
[issues-url]: https://github.com/GOeda-Co/backend/issues
[license-shield]: https://img.shields.io/github/license/GOeda-Co/backend.svg?style=for-the-badge
[license-url]: https://github.com/GOeda-Co/backend/blob/main/LICENSE
[goose]: https://github.com/pressly/goose
[swaggo]: https://github.com/swaggo/swag
[lingua-go]: https://github.com/pemistahl/lingua-go
