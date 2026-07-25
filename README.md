# 🔗 ZipURL

> A production-ready URL shortening service built with **Go**, **Echo v5**, **PostgreSQL**, and **GORM**.

ZipURL is a RESTful API that allows users to generate and manage short URLs with support for **JWT authentication**, **URL expiration**, **maximum click limits**, **anonymous link creation**, and **user-specific URL ownership**.

The project follows clean backend architecture with repository patterns, SQL migrations, structured logging, security middleware, and configurable environments.

---

## ✨ Features

### URL Shortening

- Create short URLs from long URLs
- Anonymous URL shortening
- Authenticated URL shortening
- Base62 encoded short tokens
- ULID-based unique identifiers
- Duplicate URL detection

### User Management

- User registration
- Secure user login
- JWT Authentication
- Password hashing using bcrypt
- Protected endpoints

### URL Restrictions (Authenticated Users)

- URL expiration
- Maximum click limit
- User-owned URLs
- Anonymous users cannot create restricted URLs

### Redirect Service

- Fast URL redirection
- Automatic click counting
- Expired URL validation
- Maximum click validation

### Infrastructure

- PostgreSQL
- GORM ORM
- Goose SQL migrations
- Zerolog structured logging
- Request & Correlation IDs
- Graceful shutdown
- Configurable CORS
- Security headers
- English & Bangla responses

---

# 🏗 Architecture

```
                     Client
                        │
                        ▼
                  Echo HTTP Server
                        │
        ┌───────────────┴───────────────┐
        │                               │
        ▼                               ▼
    Middleware                     Route Layer
                                        │
                                        ▼
                                   Handlers
                                        │
                                        ▼
                                 Repository Layer
                                        │
                                        ▼
                                     GORM ORM
                                        │
                                        ▼
                                   PostgreSQL
```

---

# 📁 Project Structure

```
.
├── cmd
│   ├── config
│   ├── dts
│   ├── handler
│   ├── middleware
│   ├── model
│   ├── repository
│   ├── response
│   ├── route
│   └── utils
│
├── migrations
├── main.go
├── Makefile
└── README.md
```

---

# 🚀 Tech Stack

| Technology | Purpose |
|------------|---------|
| Go | Programming Language |
| Echo v5 | HTTP Framework |
| PostgreSQL | Database |
| GORM | ORM |
| Goose | SQL Migration |
| Zerolog | Structured Logging |
| Viper | Configuration |
| JWT | Authentication |
| bcrypt | Password Hashing |

---

# 🚀 Getting Started

## Prerequisites

- Go 1.25+
- PostgreSQL
- Goose CLI

Install Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

---

## Clone Repository

```bash
git clone https://github.com/PritomKarmokar/zipurl.git

cd zipurl
```

---

## Install Dependencies

```bash
go mod tidy
```

---

## Configure Environment

Copy `.env.example` to `.env`

```env
cp .env.example .env
# Update `.env` with DB credentials & secret key
```
---

## Run Database Migrations

```bash
make migrate-up
```

Useful commands

```bash
make migrate-status

make migrate-down

make migrate-create name=create_table
```

---

## Start Server

```bash
go run main.go
```

Server

```
http://localhost:8080
```

---

# 📖 API

## Health Check

```
GET /zip-url/health/live
```

---

## User Registration

```
POST /zip-url/api/v1/user/signup
```

Request

```json
{
  "username": "pritom",
  "first_name": "Pritom",
  "last_name": "Karmokar",
  "email": "pritom@example.com",
  "password": "Password@123"
}
```

---

## User Login

```
POST /zip-url/api/v1/user/login
```

Request

```json
{
  "email": "pritom@example.com",
  "password": "Password@123"
}
```

Response

```json
{
  "access_token": "...",
  "refresh_token": "...",
  "token_type": "Bearer"
}
```

---

## Create Short URL (Anonymous)

```
POST /zip-url/api/v1/url/shorten
```

```json
{
    "url":"https://google.com"
}
```

Response

```json
{
    "short_url":"http://localhost:8080/AbCd123"
}
```

---

## Create Short URL (Authenticated)

Header

```
Authorization: Bearer <access_token>
```

Request

```json
{
    "url":"https://google.com",
    "expiry":"24h",
    "maximum_clicks":100
}
```

Authenticated users can optionally configure:

- URL expiration
- Maximum click count

Anonymous users are not allowed to use these options.

---

## Redirect

```
GET /{token}
```

If valid

```
307 Temporary Redirect
```

The click counter is automatically incremented.

If the URL

- doesn't exist
- has expired
- exceeds maximum clicks

an error response is returned.

---

# 🔐 Authentication

Protected endpoints require

```
Authorization: Bearer <JWT_ACCESS_TOKEN>
```

The API verifies

- JWT signature
- User existence
- User status
- Token validity

---

# 🌍 Localization

Supported languages

```
Accept-Language: en
```

or

```
Accept-Language: bn
```


---

# 📊 URL Lifecycle

```
Client
   │
   ▼
Create URL
   │
   ▼
Store in PostgreSQL
   │
   ▼
Receive Short URL
   │
   ▼
User Opens Link
   │
   ▼
Validate
   │
   ├── Exists?
   ├── Expired?
   ├── Max Clicks Reached?
   │
   ▼
Increment Click Count
   │
   ▼
302 Redirect
```

---

# 🚀 Future Improvements

- URL analytics dashboard
- QR code generation
- Custom aliases
- URL management APIs
- Redis caching
- Rate limiting
- Swagger/OpenAPI documentation
- Docker & Docker Compose support
- GitHub Actions CI/CD

---

# 📄 License

This project is licensed under the MIT License.