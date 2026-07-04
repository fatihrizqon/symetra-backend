# go-fiber-service

A production-oriented backend service built with **Golang** and **Fiber**, emphasizing maintainability, clear boundaries, and long-term scalability. This project demonstrates authentication using **JWT**, a structured **Clean Architecture** approach, and API documentation generated automatically using **Swagger (Swag)**.

---

## ✨ Key Capabilities

- RESTful API built with **Fiber v3**
- **JWT-based authentication** (access token)
- **Dedicated Role-Based Access Control (RBAC)** local middleware (`internal/rbac`)
- Clean Architecture (separation of concerns)
- Request validation using `go-playground/validator`
- PostgreSQL integration using **GORM**
- Centralized and structured logging via a dedicated `logger` module using **Logrus**
- Environment configuration via `.env`
- **Swagger API documentation** (auto-generated)
- Test-ready setup using `stretchr/testify`

---

## 🏗️ Architecture Overview

This project applies **Clean Architecture principles** to enforce clear boundaries between business logic and infrastructure concerns:

- Frameworks and third-party libraries live at the outer layer
- Business rules remain independent and testable
- Infrastructure concerns such as **logging**, database access, and HTTP handling are isolated

Core layers are separated into:

- **Handlers / Controllers** – HTTP request handling and response mapping
- **Use Cases / Services** – business logic orchestration
- **Repositories** – data access abstraction
- **Entities / Domain Models** – core business definitions

Cross-cutting concerns like logging and configuration are centralized (e.g. `logger.go`) to avoid duplication and tight coupling, improving long-term maintainability and operational visibility.

---

## 🛠️ Tech Stack

- **Go** 1.26
- **Fiber** v3
- **GORM** (PostgreSQL)
- **JWT** (`golang-jwt/jwt`)
- **Swagger / OpenAPI** (`swaggo/swag`)
- **Validator** (`go-playground/validator`)
- **Logrus**

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/fatihrizqon/gofiber-microservice.git
cd go-fiber-service
```

### 2. Setup Environment Variables

Create a `.env` file:

```env
APP_PORT=3000
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_fiber_service

JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Application

```bash
go run main.go
```

The server will start on:

```
http://localhost:3000
```

---

## 🔐 Authentication

Authentication is implemented using **JWT**:

- Users authenticate via a login endpoint
- A signed JWT is returned upon successful authentication
- Protected routes require a valid `Authorization: Bearer <token>` header

JWT handling is isolated within the authentication layer to keep business logic clean.

---

## 📚 API Documentation (Swagger)

This project uses **Swag** to generate Swagger documentation automatically from code annotations.

### Install Swag CLI

Make sure you have Go installed, then run:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generate Swagger Docs

```bash
swag init
```

### Access Swagger UI

Once the application is running:

```
http://localhost:3000/swagger/index.html
```

---

## 🧪 Testing

Tests are structured to support unit and service-level testing.

Run all tests:

```bash
go test ./...
```

---

## 📁 Project Structure (Simplified)

```
├── config/
├── docs/          # Swagger generated files
├── util/
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── entity/
│   ├── middleware/
│   └── rbac/      # Dedicated RBAC local middleware
├── logger/
├── middleware/
├── router/
├── test/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

---

## 🎯 Purpose

This repository demonstrates senior-level backend engineering practices, including:

- Designing a maintainable Go service with clear architectural boundaries
- Implementing stateless JWT authentication suitable for distributed systems
- Centralizing logging through a shared logger to improve observability and debuggability
- Maintaining framework-agnostic business logic
- Applying API documentation and testing best practices

It is intended as a **senior backend portfolio project** and a reference for designing maintainable Go services in real-world environments.

---

## 📄 License

This project is open-source and available under the **MIT License**.

---

**Author**  
Fatih Rizqon


## 📦 Recent Updates (Outbox Pattern & Worker Refactoring)

### 1. Robust Background Jobs (Outbox Pattern)
- **Transactional Enqueueing**: Instead of pushing directly to Redis, asynchronous tasks (like sending emails) are now persisted in the PostgreSQL `redis_jobs` table atomically with business entities (e.g., User Registration).
- **Guaranteed Delivery**: A dedicated background sweeper (`internal/worker/sweeper.go`) picks up `PENDING` jobs and forwards them to Asynq. This guarantees tasks are never lost even if Redis goes down.

### 2. Independent Worker Architecture
- **Worker Main (`cmd/worker/main.go`)**: Separated from the HTTP web server.
- **Worker Bootstrap**: Logic has been encapsulated in `config/worker.go` (`BootstrapWorker`) for cleaner initialization of the DB, Asynq Client, Job Sweeper, and Worker Processors, ensuring architectural consistency.

### 3. Full RBAC Management (CRUD)
- **Role & Permission Entities**: Complete CRUD endpoints have been added for `Roles` and `Permissions` to allow Admins (`rbac.manage`) full control over access rules.
- **Admin Endpoints**: New API routes under `/api/v1/roles` and `/api/v1/permissions` with robust data validation and error handling integrated into the layered architecture.

---

## 📐 Architectural Standardization (Boilerplate Pattern)

This boilerplate is designed to be highly reusable as a foundation for various microservices, including highly complex domains like enterprise accounting systems. It strictly enforces a **Layered Clean Architecture** to maintain separation of concerns:

- **`internal/delivery` (HTTP Layer):** Handles routing, request parsing (JSON), and sending responses. Must contain **no business logic** or direct database queries.
- **`internal/service` (Business Logic Layer):** Contains the core business rules. Services act as orchestrators and can interact with multiple repositories (or other services via Dependency Injection) to handle complex processes.
- **`internal/repository` (Data Layer):** Handles direct database access and complex **ACID Transactions**. Database transactions spanning multiple tables (e.g., Journal Entries and Lines) must be safely abstracted here.
- **`internal/worker` (Background Jobs):** Handles asynchronous processing (e.g., generating heavy PDF reports, sending email notifications) without blocking the primary HTTP event loop.

When extending this boilerplate or instructing Junior Developers / AI Assistants to implement new PRDs, **this pattern is the absolute standard**. Ensure that database operations never leak into the `delivery` layer and that background tasks are routed to the `worker` layer, ensuring the system remains decoupled, scalable, and testable.
