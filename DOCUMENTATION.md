# Codebase Documentation

This document explains the purpose and functionality of each Go source file in the project. The project follows a **Clean Architecture** (Layered Architecture) pattern.

**Catatan migrasi:** Proyek ini telah dipindahkan dari SQLite ke Supabase (Postgres). Perubahan termasuk penggunaan `DATABASE_URL` untuk koneksi, penggunaan driver `gorm.io/driver/postgres`, dan pengaturan seeding runtime. Lihat `DEPLOYMENT.md` untuk detail.


## 📂 `cmd/api/`

### `main.go`
**Purpose**: The entry point of the application.
- **Responsibilities**:
    - Initializes the database connection (`database.Connect()`).
    - Sets up the dependency injection wiring (Repository -> Service -> Handler).
    - Configures the Gin router.
    - Registers middleware (Rate Limiting) and routes.
    - Starts the HTTP server on port `8080`.

## 📂 `internal/entity/`

### `category.go`
**Purpose**: Defines the domain model.
- **Components**:
    - `Category`: A struct representing the category table in the database.
    - **Tags**: Includes `gorm` tags for database constraints and `json` tags for API serialization.

## 📂 `internal/repository/`

### `category_repository.go`
**Purpose**: Handles direct database interactions.
- **Interface** `CategoryRepository`: Defines the contract for data access.
- **Implementation**: Uses GORM methods (`Create`, `Find`, `Save`, `Delete`) to perform CRUD operations against the Postgres database (Supabase).

## 📂 `internal/service/`

### `category_service.go`
**Purpose**: Contains business logic.
- **Interface** `CategoryService`: Defines the business operations available.
- **Implementation**:
    - Acts as a bridge between the Handler and the Repository.
    - Can handle validation or business rules before calling the repository (though currently simple pass-through for basic CRUD).

## 📂 `internal/handler/`

### `category_handler.go`
**Purpose**: Manages HTTP Requests and Responses.
- **Responsibilities**:
    - Parses incoming JSON bodies.
    - Validates inputs (binding).
    - Calls the Service layer.
    - Returns appropriate HTTP status codes and JSON responses.
    - **Swagger Annotations**: Contains comments used to generate the API documentation (`godoc`).

## 📂 `internal/middleware/`

### `rate_limit.go`
**Purpose**: Protects the API from abuse.
- **Functionality**:
    - Implements a Token Bucket algorithm using `uber-go/ratelimit`.
    - `RateLimitMiddleware`: A Gin middleware that blocks requests if the rate limit (e.g., 10 req/sec) is exceeded.

## 📂 `pkg/database/`

### `db.go`
**Purpose**: Centralized database configuration.
- **Functionality**:
    - `Connect()`: Initializes the connection to Postgres using the `DATABASE_URL` environment variable (e.g., Supabase connection string).
    - Performs **Auto Migration** to create/update database tables based on the Entity struct.
    - Uses the `gorm.io/driver/postgres` driver for GORM.
