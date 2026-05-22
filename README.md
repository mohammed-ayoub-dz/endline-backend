# Endline Backend (Go)

A secure, high‑performance backend for the Endline e‑learning platform. Built with **Go**, **Fiber**, **GORM**, and **PostgreSQL**.

## Features

-  **JWT Authentication** – Access + refresh tokens, stored in HTTP‑only cookies.
-  **High security** – Argon2id password hashing, rate limiting, secure headers.
-  **Video progress tracking** – Save and resume watched seconds per user/video.
-  **User library** – Save favourite videos (YouTube links) to a personal list.
-  **Full‑text search** – Search videos (from `videos.json`) by title/description.
-  **Pomodoro timer state** – Persist work/break durations and remaining time in DB.
-  **PostgreSQL + GORM** – Auto‑migration, clean relationships, indexes.
-  **Modular architecture** – `cmd/`, `internal/`, `pkg/` layout.

## Tech Stack

| Component          | Technology                             |
|--------------------|----------------------------------------|
| Language           | Go 1.21+                               |
| Web framework      | [Fiber](https://gofiber.io/) v2        |
| ORM                | [GORM](https://gorm.io/)               |
| Database           | PostgreSQL 15+                         |
| Auth               | JWT (access + refresh), Argon2id       |
| Security           | Helmet, Rate limiter, CORS             |

## Project Structure

```
.
├── cmd
│   └── server
│       └── main.go                 # entry point
├── internal
│   ├── config
│   │   └── config.go               # environment variables
│   ├── database
│   │   └── db.go                   # GORM connection & migration
│   ├── models
│   │   ├── user.go
│   │   ├── video.go
│   │   ├── subject.go
│   │   ├── refresh_token.go
│   │   ├── user_video_progress.go
│   │   └── saved_video.go
│   ├── handler
│   │   ├── auth.go                 # register, login, refresh, logout
│   │   ├── video.go                # search, save to library, get library
│   │   └── progress.go             # update/get video progress
│   ├── middleware
│   │   ├── auth.go                 # JWT validation
│   │   ├── rate_limit.go           # IP‑based rate limiting
│   │   └── security.go             # secure headers (helmet)
│   └── utils
│       ├── argon2.go               # password hashing
│       └── jwt.go                  # token generation & parsing
├── data
│   └── videos.json                 # searchable video catalogue
├── .env                            # environment variables
├── go.mod
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Make (optional)

### Installation

2. **Install dependencies**

```bash
go mod tidy
```

3. **Set up the database**

Create a PostgreSQL database:

```sql
CREATE DATABASE endline;
```

4. **Environment variables**

Copy the example file and edit:

```bash
cp .env.example .env
```

Minimal `.env` content:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=endline
DB_PORT=5432
DB_SSLMODE=disable
ACCESS_SECRET=your_super_secret_key_min_32_chars
REFRESH_SECRET=another_secret_key_min_32_chars
ACCESS_TTL=15m
REFRESH_TTL=720h
```

5. **Run the server**

```bash
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`.

## API Endpoints

All endpoints are prefixed with `/api`.

### Public routes

| Method | Endpoint           | Description                     |
|--------|--------------------|---------------------------------|
| POST   | `/auth/register`   | Create a new user              |
| POST   | `/auth/login`      | Login, returns access token    |
| POST   | `/auth/refresh`    | Get new access token (uses cookie) |
| POST   | `/auth/logout`     | Invalidate refresh token       |

### Protected routes (require JWT)

Add `Authorization: Bearer <access_token>` header.

| Method | Endpoint                     | Description                           |
|--------|------------------------------|---------------------------------------|
| GET    | `/protected/profile`         | Get current user profile              |
| POST   | `/protected/user/videos`     | Save a video to user's library        |
| GET    | `/protected/user/videos`     | Get all saved videos for the user     |
| POST   | `/protected/progress`        | Update watched time for a video       |
| GET    | `/protected/progress/:videoId` | Get saved progress for a video      |
| GET    | `/protected/videos/search?q=` | Search videos in `videos.json`      |

>  Refresh token is stored in an HTTP‑only cookie. The `/refresh` endpoint reads it automatically.

## Database Schema (selected tables)

### `users`

| Column     | Type      | Description               |
|------------|-----------|---------------------------|
| id         | SERIAL PK |                           |
| username   | TEXT UNIQUE| Login name                |
| password   | TEXT      | Argon2id hash             |
| points     | INT       | Default 0                 |
| level      | TEXT      | Default 'beginner'        |
| created_at | TIMESTAMP |                           |

### `videos` (global catalogue)

| Column      | Type      | Description               |
|-------------|-----------|---------------------------|
| id          | SERIAL PK |                           |
| title       | TEXT      |                           |
| description | TEXT      |                           |
| url         | TEXT UNIQUE| YouTube URL               |
| thumbnail   | TEXT      |                           |
| subject_id  | INT (FK)  |                           |

### `user_videos` (many‑to‑many library)

| Column   | Type      | Description |
|----------|-----------|-------------|
| user_id  | INT (FK)  |             |
| video_id | INT (FK)  |             |
| added_at | TIMESTAMP |             |

### `video_progress`

| Column          | Type      | Description               |
|-----------------|-----------|---------------------------|
| id              | SERIAL PK |                           |
| user_id         | INT (FK)  |                           |
| video_id        | INT (FK)  |                           |
| watched_seconds | INT       |                           |
| completed       | BOOLEAN   |                           |
| updated_at      | TIMESTAMP |                           |

## Environment Variables (full list)

| Variable         | Description                             | Default      |
|------------------|-----------------------------------------|--------------|
| DB_HOST          | PostgreSQL host                         | localhost    |
| DB_USER          | Database user                           | postgres     |
| DB_PASSWORD      | Database password                       | (empty)      |
| DB_NAME          | Database name                           | postgres     |
| DB_PORT          | PostgreSQL port                         | 5432         |
| DB_SSLMODE       | SSL mode (disable / require)            | disable      |
| ACCESS_SECRET    | Secret for access tokens (≥32 chars)    | *required*   |
| REFRESH_SECRET   | Secret for refresh tokens (≥32 chars)   | *required*   |
| ACCESS_TTL       | Access token lifetime (e.g. 15m, 24h)   | 15m          |
| REFRESH_TTL      | Refresh token lifetime (e.g. 720h)      | 720h         |
| PORT             | Server listen port                       | 8080         |

## Development

### Live reload

Use [Air](https://github.com/cosmtrek/air) for hot reloading during development:

```bash
go install github.com/cosmtrek/air@latest
air
```

### Testing

Run all tests:

```bash
go test ./...
```

### Building for production

```bash
go build -o endline-api cmd/server/main.go
```

## Security Notes

- Passwords are hashed with **Argon2id** (memory‑hard, resistant to GPU attacks).
- Access tokens are short‑lived (default 15 min). Refresh tokens are rotated after each use.
- Refresh token stored in **HTTP‑only, Secure, SameSite=Strict** cookie.
- Rate limiting: 10 requests per minute per IP for auth routes.
- Security headers: Helmet middleware disables X‑Powered‑By, enables X‑Frame‑Options, etc.
- CORS restricts to frontend origin(s) with credentials.

