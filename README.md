# Chirpy

Chirpy is a small Go HTTP API for creating and reading short messages called chirps. It stores users, chirps, and refresh tokens in PostgreSQL and serves a static app from `/app/`.

## Requirements

- Go version compatible with the `go 1.27.2` directive in `go.mod`
- PostgreSQL
- A migration runner that understands the Goose `-- +goose Up` / `-- +goose Down` annotations in `sql/schema`

## Configuration

Create a `.env` file in the project root, or set these environment variables in your shell:

| Variable | Required | Purpose |
| --- | --- | --- |
| `DB_URL` | Yes | PostgreSQL connection string, for example `postgres://user:password@localhost:5432/chirpy?sslmode=disable` |
| `JWT_SECRET` | For login and protected routes | Secret used to sign access tokens |
| `POLKA_KEY` | For the Polka webhook | API key used to authenticate webhook requests |
| `PLATFORM` | No | Set to `dev` to enable the development-only database reset endpoint |

The application loads `.env` when present, but environment variables already provided by the process can also be used. Do not commit secrets; `.env` is ignored by Git.

## Database setup

Create a PostgreSQL database and set `DB_URL` to its connection string. Apply the migrations in `sql/schema` in numeric order (`001` through `005`) with a Goose-compatible migration tool. The project includes the migration files and `sqlc` query definitions; it does not include a migration command or database bootstrap script.

## Run

From the repository root:

```sh
go run .
```

The server listens on port `8080`. Alternatively, build and run it:

```sh
go build -o chirpy .
./chirpy
```

The root directory is the static file server root, mounted at `/app/`.

## API

JSON endpoints use `Content-Type: application/json` for request bodies. Protected user and chirp operations expect `Authorization: Bearer <access-token>`. A login access token lasts one hour; refresh tokens are stored in the database and expire after 60 days.

| Method | Path | Description | Authentication |
| --- | --- | --- | --- |
| `GET` | `/api/healthz` | Readiness check (`200 OK`) | None |
| `POST` | `/api/users` | Create a user with `email` and `password` | None |
| `PUT` | `/api/users` | Update the signed-in user's email and password | Bearer access token |
| `POST` | `/api/login` | Authenticate with `email` and `password`; returns access and refresh tokens | None |
| `POST` | `/api/refresh` | Exchange a refresh token for a new access token; send the refresh token as a Bearer token | Refresh token |
| `POST` | `/api/revoke` | Revoke a refresh token; send it as a Bearer token | Refresh token |
| `POST` | `/api/chirps` | Create a chirp with a `body` (up to 140 characters) | Bearer access token |
| `GET` | `/api/chirps` | List chirps; optional `author_id` UUID and `sort=asc\|desc` query parameters | None |
| `GET` | `/api/chirps/{chirpID}` | Fetch one chirp | None |
| `DELETE` | `/api/chirps/{chirpID}` | Delete a chirp owned by the signed-in user | Bearer access token |
| `POST` | `/api/polka/webhooks` | Process supported `user.upgraded` events | `Authorization: ApiKey <POLKA_KEY>` |
| `GET` | `/admin/metrics` | Show static-file request count | None |
| `POST` | `/admin/reset` | Delete all users (and their dependent data) | Only when `PLATFORM=dev` |

For the Polka webhook, the supported request shape is:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "00000000-0000-0000-0000-000000000000"
  }
}
```

Chirp creation replaces the words `kerfuffle`, `sharbert`, and `fornax` with `****`, case-insensitively.

## Project layout

- `main.go` and `handler_*.go`: HTTP server, routing, and endpoint handlers
- `internal/auth`: password hashing, JWTs, refresh token generation, and authorization-header parsing
- `sql/schema`: PostgreSQL schema migrations
- `sql/queries`: SQL query definitions
- `internal/database`: `sqlc`-generated database code
- `sqlc.yaml`: `sqlc` generation configuration

To regenerate the database package after changing SQL definitions, install `sqlc` and run:

```sh
sqlc generate
```

## Tests

Run the Go test suite with:

```sh
go test ./...
```
