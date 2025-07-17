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
    <img src="https://avatars.githubusercontent.com/u/219566722?s=200&v=4" alt="Logo" width="80" height="80">
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

To get a local copy up and running, follow these steps.

### Prerequisites

- Go 1.23+
- PostgreSQL 15

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/GOeda-Co/backend.git
   cd repeatro
   ```
2. Install dependencies
   ```sh
   go mod tidy
   ```
3. Create the database
   ```sql
   CREATE DATABASE repeatro;
   ```
4. Configure your environment (see `config.example.toml` in the root)
5. Start the backend
   ```sh
   go run backend/card/cmd/card/main.go
   # or for other services, adjust the path accordingly
   ```
7. (Optional) Use [air](https://github.com/cosmtrek/air) for auto server restart

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
