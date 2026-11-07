# Go E-commerce API

A REST API for an online store, written in **Go** with **Clean Architecture**, the **Repository pattern** and dependency injection. It handles users, authentication (JWT and OAuth), products, categories, shipping addresses, saved cards and orders, with payments processed through **Stripe**.

> **This is part 1 of 2.** It is the backend. The client is a separate **React** project that consumes this API.

I built it as a showcase project: the goal was to apply layered architecture, decoupled persistence and testing in Go to a realistic domain, not just a CRUD.

---

## Highlights

- **Clean Architecture**: the domain layer has no framework or database dependencies. Handlers depend on services, and services depend on repository *interfaces*.
- **Repository pattern with generics**: a single generic implementation (`MemoryRepo[T]`, `FirestoreRepo[T]`) covers CRUD for every model, and each entity only adds what is specific to it.
- **Swappable persistence**: the same code runs on **Firestore** (production), on **JSON files** (local development) or **in memory** (tests), by changing one line in the composition root.
- **Authentication**: JWT access and refresh tokens, bcrypt-hashed passwords, role-based authorization (`user`, `moderator`, `admin`) and OAuth sign-in with **Google, GitHub and Discord**.
- **Payments with Stripe**: customers are created on demand, cards are saved as payment methods and orders are charged with PaymentIntents.
- **Tests** at every layer: generic repositories, services, utilities and the HTTP routes, using `testify/suite`.
- **Documentation** with Swagger/OpenAPI, generated from annotations on the handlers.

## Tech stack

| Area | Tools |
|---|---|
| Language | Go (generics) |
| HTTP | [Fiber v2](https://github.com/gofiber/fiber) |
| Database | Cloud Firestore (Firebase Admin SDK) |
| Payments | [Stripe](https://github.com/stripe/stripe-go) |
| Auth | JWT (`golang-jwt`), bcrypt, OAuth 2.0 |
| Validation | `go-playground/validator` |
| Testing | `testify` (suites and asserts) |
| Docs | Swagger (`swaggo/swag`) |
| Dev tooling | `air` (hot reload), Make |

## Architecture

```
            HTTP request
                 │
┌────────────────▼─────────────────┐
│  API layer        src/api        │  Routes, middlewares, handlers, DTOs
│  (Fiber)                         │  Parses and validates input, maps to HTTP
└────────────────┬─────────────────┘
                 │ depends on domain.XService
┌────────────────▼─────────────────┐
│  Service layer    src/services   │  Business rules (core) and integrations:
│                                  │  auth, payment, email, upload, validation
└────────────────┬─────────────────┘
                 │ depends on domain.XRepository
┌────────────────▼─────────────────┐
│  Repository layer                │  Firestore │ JSON file │ In-memory
│                src/repositories  │  (same interface, swappable)
└──────────────────────────────────┘

        src/domain: models + interfaces (no dependencies)
```

The **composition root** is [`cmd/main.go`](cmd/main.go). It builds the repositories, injects them into the services, the services into the handlers, and registers the routes. Choosing between Firestore and local persistence happens only there:

```go
if devMode {
    userRepo = user.NewJSONUserRepository("db/users.json")
} else {
    userRepo = user.NewFirestoreUserRepository(client, utils.UserColl)
}
userSvc := core.NewUserService(userRepo) // the service doesn't know which one it got
```

### Project structure

```
.
├── cmd/main.go              # Entry point and dependency wiring
├── config/                  # Config loading (app_config.json, firebase_config.json)
├── docs/                    # Generated Swagger docs
├── db/                      # Local JSON data for dev mode (git-ignored)
└── src/
    ├── domain/              # Models + Service/Repository interfaces
    ├── repositories/
    │   ├── shared/          # Generic MemoryRepo[T] and FirestoreRepo[T]
    │   └── <entity>/        # memory.go (memory + JSON) and firestore.go
    ├── services/
    │   ├── core/            # Business logic: users, products, categories, addresses, orders
    │   ├── auth/            # JWT + OAuth providers (Google, GitHub, Discord)
    │   ├── payment/         # Stripe
    │   ├── email/           # SMTP and HTML templates
    │   ├── upload/          # Firebase Storage
    │   └── validation/      # Struct validation
    ├── api/
    │   ├── handlers/        # One package per resource
    │   ├── middlewares/     # Auth, roles and Stripe customer
    │   ├── dtos/            # Request/response contracts
    │   └── test/            # HTTP route tests
    └── utils/               # Thread-safe SyncMap, helpers, constants
```

### Generic repositories

Most of the CRUD lives in a generic type parameterized by the domain model:

```go
type DomainModel interface {
    User | Address | Category | Product | Order | ExampleModel
    GetStringID() string
    GetCreatedDate() int64
}

type MemoryRepo[T DomainModel] struct {
    Store *utils.SyncMap[uuid.UUID, T]
}
```

Each concrete repository embeds the generic one, so it only implements what is specific to that entity. The in-memory store is a `SyncMap` protected by a `RWMutex`, which can optionally persist to a JSON file. That is how dev mode keeps data between restarts without needing a database.

## Domain models

| Model | Main fields |
|---|---|
| **User** | `username`, `email`, `password` (bcrypt, never returned), `role`, `verified_email`, `image_url`, `customer_id` (Stripe) |
| **Product** | `name`, `description`, `category`, `price`, `discount_rate`, `images_url`, `tags`, `available` |
| **Category** | `name` |
| **Address** | `user_id`, `name`, `line1`, `line2`, `city`, `state`, `country`, `postal_code` |
| **Order** | `user_id`, `address_id`, `products` (`product_id` + `quantity`), `amount`, `status`, `paid`, `payment_id` |

Every model embeds `Model { id (UUID), created_at, updated_at }`.

Order status values: `pending`, `processing`, `shipped`, `completed`, `cancelled`.

## Endpoints

Base path: `/api`. Full interactive documentation is available at `/docs/index.html`.

| Method | Route | Access |
|---|---|---|
| **Auth** | | |
| POST | `/auth/signup` | Public |
| POST | `/auth/signin` | Public |
| GET | `/auth/refresh` | Refresh token |
| GET | `/auth/me` | Authenticated |
| GET | `/auth/signout` | Public |
| GET | `/auth/:provider/url` | Public (`google`, `github`, `discord`) |
| GET | `/auth/:provider/callback` | OAuth callback |
| **Users** | | |
| GET | `/user/all`, `/user/get/:id` | Moderator |
| POST / PUT / DELETE | `/user/create`, `/user/update/:id`, `/user/delete/:id` | Admin |
| **Products** | | |
| GET | `/product/all`, `/product/get/:id` | Public |
| POST / PUT / DELETE | `/product/create`, `/product/update/:id`, `/product/delete/:id` | Admin |
| **Categories** | | |
| GET | `/category/all` | Public |
| POST / DELETE | `/category/create`, `/category/delete/:id` | Moderator |
| **Addresses** | | |
| GET / POST | `/address/list`, `/address/create` | Authenticated |
| PUT / DELETE | `/address/update/:id`, `/address/delete/:id` | Owner |
| **Cards (Stripe)** | | |
| GET / POST / DELETE | `/card/list`, `/card/save`, `/card/remove/:id` | Authenticated |
| **Orders** | | |
| GET | `/order/list` | Authenticated (own orders) |
| GET | `/order/get/:id` | Owner (order with address and products) |
| POST | `/order/new` | Authenticated (charges with Stripe) |

### Authentication flow

1. `signin` returns a short-lived **access token** (5 min) and sets a **refresh token** (5 days) in an `HttpOnly` cookie.
2. Protected routes read the access token from the `X-Access-Token` header. The `AuthRequired` and `RoleRequired(role)` middlewares validate the token and the role.
3. `refresh` issues a new access token from the refresh token.
4. With OAuth, the provider redirects to `/:provider/callback`. The API creates or finds the user and redirects to the client with the session already started.

### Purchase flow

1. The `CustomerIDRequired` middleware gets or creates the user's **Stripe customer** and stores its ID.
2. `POST /order/new` validates the address and the products, and **calculates the total on the server** (discounts included), never trusting a price sent by the client.
3. The order is created, the payment is charged with a **PaymentIntent** using the chosen payment method, and the order is marked as paid.

## Getting started

### Requirements

- Go 1.18+
- A Firebase project with Firestore enabled
- A Stripe account (test mode is enough)
- Optional: [`air`](https://github.com/air-verse/air) for hot reload

### Configuration

Configuration lives in two git-ignored files inside `config/`. Copy the templates and fill them in:

```bash
cp config/app_config.example.json config/app_config.json
cp config/firebase_config.example.json config/firebase_config.json
```

- **`app_config.json`**: port, server and client URLs, JWT secrets, header names, OAuth credentials, SMTP and Stripe keys.
- **`firebase_config.json`**: the Firebase *service account* key (Firebase Console → Project settings → Service accounts → *Generate new private key*).

For the `make` frontend commands, copy `.env.example` to `.env` and set the path and URL of the frontend.

### Running

```bash
make run-backend-dev        # API in dev mode (local JSON data) with hot reload
go run cmd/main.go          # API with Firestore
NODE_ENV=production go run cmd/main.go   # production mode (CORS restricted to client_origin)
```

| Command | Description |
|---|---|
| `make dev` | Backend, frontend and opens the browser |
| `make run-backend-dev` | API with `air` in dev mode |
| `make run-frontend-dev` | Runs the React client (`FRONTEND_DIR`) |
| `make gen-docs` | Regenerates the Swagger docs |
| `make install-swagger` | Installs `swag` v1.8.9 |

Once it is running, the docs are at **http://localhost:9000/docs/index.html**.

## Testing

```bash
go test ./src/services/... ./src/repositories/... ./src/utils/...
```

| Layer | What is tested |
|---|---|
| `repositories/shared` | Generic CRUD of `MemoryRepo[T]` and `FirestoreRepo[T]` |
| `services/core` | Business rules on top of in-memory repositories |
| `utils` | Thread-safe `SyncMap` with file persistence, helpers |
| `api/test` | HTTP routes end to end with `app.Test()` |

Service tests use in-memory repositories, so they are fast and need no external services. Firestore, Stripe, SMTP and HTTP route tests are **integration tests**: they need real credentials in `config/` and talk to those services.

---

Made by [ZaphCode](https://github.com/ZaphCode) (Omar Urquidez)
