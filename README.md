# ZipURL

ZipURL is a lightweight URL shortener API built with Go, Echo, PostgreSQL, and GORM. It exposes an endpoint to create short links and a public redirect route that sends users to the original URL.

## Features

- Create short URLs from long URLs
- Redirect short tokens to original URLs
- PostgreSQL persistence with GORM
- Goose SQL migrations
- ULID-based IDs and Base62 short tokens
- Structured JSON logging with Zerolog
- Request and correlation ID tracking
- Security headers and configurable CORS
- Graceful server shutdown
- Basic health check endpoint
- English and Bangla API response messages

## Tech Stack

- **Language:** Go
- **HTTP framework:** Echo v5
- **Database:** PostgreSQL
- **ORM:** GORM
- **Migrations:** Goose
- **Configuration:** Viper (`.env` or environment variables)
- **Logging:** Zerolog

## Project Structure

```text
.
├── main.go                         # Application bootstrap
├── makefile                        # Migration helper commands
├── migrations/                     # Goose database migrations
└── cmd/
    ├── config/                     # Env, DB, logger, Echo, validator, server config
    ├── dts/                        # Request DTOs
    ├── handler/                    # HTTP handlers
    ├── middleware/                 # Correlation ID, CORS, security headers
    ├── model/                      # GORM models
    ├── repository/                 # Database access layer
    ├── response/                   # Standard API response objects
    ├── route/                      # Route registration
    └── utils/                      # ULID and Base62 utilities
```

## Requirements

- Go `1.25.11` or compatible version defined in `go.mod`
- PostgreSQL
- Goose CLI for migrations

Install Goose if it is not already available:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Configuration

ZipURL loads configuration from a `.env` file in the project root. If `.env` is not found, it falls back to system environment variables.

Create a `.env` file:

```env
SERVER_PORT=8080
ZIP_URL_BASE_URL=http://localhost:8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=zip_url_db
DB_SSLMODE=disable
TIME_ZONE=Asia/Dhaka

LOG_LEVEL_API=info
SLOW_SQL_THRESHOLD=200

ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
ENABLE_HSTS=false
CONTENT_SECURITY_POLICY=default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'
```

### Environment Variables

| Variable | Description | Example |
| --- | --- | --- |
| `SERVER_PORT` | Port used by the HTTP server | `8080` |
| `ZIP_URL_BASE_URL` | Base URL used when returning generated short links | `http://localhost:8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL user | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `postgres` |
| `DB_NAME` | PostgreSQL database name | `zip_url_db` |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` |
| `TIME_ZONE` | Database timezone | `Asia/Dhaka` |
| `LOG_LEVEL_API` | Log level: `error`, `warn`, `info`, `debug`, or `trace` | `info` |
| `SLOW_SQL_THRESHOLD` | Slow SQL threshold in milliseconds | `200` |
| `ALLOWED_ORIGINS` | Comma-separated CORS origins, or `*` for development | `http://localhost:3000` |
| `ENABLE_HSTS` | Enables `Strict-Transport-Security` header | `false` |
| `CONTENT_SECURITY_POLICY` | Optional custom Content Security Policy | see sample above |

## Database Setup

Create the PostgreSQL database:

```bash
createdb zip_url_db
```

Run migrations:

```bash
make migrate-up
```

Check migration status:

```bash
make migrate-status
```

Rollback the last migration:

```bash
make migrate-down
```

Create a new migration:

```bash
make migrate-create name=create_new_table
```

> Note: the current `makefile` uses a hard-coded local PostgreSQL DSN. Update `DB_DSN` in `makefile` if your database credentials are different.

## Running the Application

Install dependencies:

```bash
go mod tidy
```               

Start the server:

```bash
go run main.go
```

The API will be available at:

```text
http://localhost:8080
```

assuming `SERVER_PORT=8080`.

## API Reference

### Health Check

```http
GET /zip-url/health/live
```

Response:

```json
{
  "status": "alive"
}
```

### Create Short URL

```http
POST /zip-url/api/v1/url/shorten
Content-Type: application/json
```

Request body:

```json
{
  "url": "https://example.com/a/very/long/url"
}
```

Success response:

```json
{
  "code": "GENERIC_SUC200",
  "message": "Request processed successfully",
  "lang": "en",
  "data": {
    "short_url": "http://localhost:8080/1AbCdEf"
  }
}
```

Example with `curl`:

```bash
curl -X POST http://localhost:8080/zip-url/api/v1/url/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com/a/very/long/url"}'
```

### Redirect to Original URL

```http
GET /:token
```

Example:

```bash
curl -I http://localhost:8080/1AbCdEf
```

If the token exists, the service returns a temporary redirect (`307`) to the original URL.

If the token is invalid, the service returns:

```json
{
  "code": "ZIPURL_INVUP401",
  "message": "Invalid short url provided",
  "lang": "en"
}
```

## Response Format

Most API responses follow this structure:

```json
{
  "code": "GENERIC_SUC200",
  "message": "Request processed successfully",
  "lang": "en",
  "data": {}
}
```

Set `Accept-Language: bn` to receive supported Bangla response messages.

## Security and Middleware

ZipURL includes middleware for:

- `X-Correlation-ID` and `X-Request-ID` propagation
- Secure response headers such as `X-Content-Type-Options`, `X-Frame-Options`, and `Content-Security-Policy`
- Optional HSTS via `ENABLE_HSTS=true`
- Configurable CORS via `ALLOWED_ORIGINS`
- Request body limit of 2 MB
- Panic recovery

## Development Notes

- Short tokens are generated from ULID-derived entropy and encoded using Base62.
- URL records are stored in the `urls` table.
- The redirect route is registered at the root path (`/:token`), while API routes are grouped under `/zip-url`.
- Logs include request and correlation IDs when available.

## License

This project is licensed under the [MIT License](LICENSE).
