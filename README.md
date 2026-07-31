# cli-auth

A secure interactive command-line authentication system built in Go. It provides user registration, password-based login with bcrypt hashing, optional TOTP two-factor authentication, account lockout, session management, and SQLite persistence — all wrapped in a containerized CLI.

## Features

- Interactive REPL with tab completion
- User registration with username + password
- Secure password storage using bcrypt
- Optional TOTP-based 2FA (RFC 6238, Google Authenticator compatible)
- Account lockout after repeated failed attempts
- Session management with sliding expiration
- SQLite persistence via GORM (file-based, survives container restarts)
- Docker image with persistent data volume

## Architecture

```
CLI
 ↓
Auth Service
 ↓
Repository
 ↓
Database (SQLite + GORM)
```

| Package            | Responsibility                                   |
| ------------------ | ------------------------------------------------ |
| `cmd`              | Application entry point and wiring               |
| `internal/auth`    | Authentication logic: register, login, TOTP 2FA  |
| `internal/session` | Session lifecycle: create, validate, refresh     |
| `internal/cli`     | Interactive REPL, commands, and handlers         |
| `internal/config`  | Environment-based configuration                  |
| `internal/database`| SQLite connection and migrations                 |
| `internal/models`  | GORM models (`User`, `Session`)                  |
| `internal/repository` | Database access isolation                     |
| `internal/utils`   | Shared helpers                                   |

## Requirements

- Go 1.26+
- Docker (optional, for containerized usage)

## Configuration

Configuration is read from a `.env` file (or environment variables):

| Variable             | Default   | Description                        |
| -------------------- | --------- | ---------------------------------- |
| `DATABASE_PATH`      | `data/app.db` | SQLite database file           |
| `SESSION_TIMEOUT`    | `1h`      | Session lifetime (Go duration)     |
| `LOCKOUT_DURATION`   | `15m`     | Account lockout period             |
| `MAX_LOGIN_ATTEMPTS` | `5`       | Failed attempts before lockout     |

Example `.env`:

```env
DATABASE_PATH=data/app.db
SESSION_TIMEOUT=1h
LOCKOUT_DURATION=15m
MAX_LOGIN_ATTEMPTS=5
```

## Running Locally

```bash
go run ./cmd
```

The database file and migrations are created automatically on startup.

## Docker

Build and run the containerized CLI:

```bash
docker compose up -d
docker attach cli-auth-app-1
```

The SQLite database is persisted in `./data` via a Docker volume and survives `docker compose down`, `docker compose up`, and container recreation.

To stop:

```bash
docker compose down
```

## Commands

| Command        | Description                                 |
| -------------- | ------------------------------------------- |
| `help`         | Show available commands                     |
| `register`     | Create a new account                        |
| `login`        | Log in to your account                      |
| `logout`       | End the current session                     |
| `status`       | Show the current session                    |
| `mfa`          | Enable TOTP two-factor authentication       |
| `disable-mfa`  | Disable two-factor authentication           |
| `quit`         | Exit the application                        |

## Usage

### Register a user

```
> register
Username: alice
Password: ********
Confirm password: ********
Registered alice. You can now log in.
```

### Log in

```
> login
Username: alice
Password: ********
Logged in as alice
```

### Enable 2FA

After logging in, run `mfa` and add the account to your authenticator app (e.g. Google Authenticator) using the printed secret or the `otpauth://` URI:

```
> mfa
MFA enabled. Add this account to your authenticator app:
  Secret: YW5IISF23BA6BK22ORLGQ2TB4SH5HG4G
  URI:    otpauth://totp/cli-auth:alice?algorithm=SHA1&digits=6&issuer=cli-auth&period=30&secret=...
```

Subsequent logins will prompt for a time-based code:

```
> login
Username: alice
Password: ********
MFA code: 123456
Logged in as alice
```

### Check session status

```
> status
Logged in as alice
Session expires: 2026-07-31 16:52:37
Last login: 2026-07-31 15:52:37
MFA: enabled
```

## Security

- Passwords are hashed with **bcrypt** (`golang.org/x/crypto/bcrypt`) — never stored in plaintext.
- Failed login attempts increment a per-user counter; reaching `MAX_LOGIN_ATTEMPTS` locks the account for `LOCKOUT_DURATION`.
- TOTP secrets are generated with a cryptographically secure random source and validated per RFC 6238.
- MFA is enforced on every login once enabled.

## Project Layout

```
.
├── cmd/main.go                 # Entry point, dependency wiring
├── internal/
│   ├── auth/                   # Auth service: register, login, lockout, TOTP
│   ├── cli/                    # REPL, commands, handlers, completion
│   ├── config/                 # AppConfig singleton (.env)
│   ├── database/               # SQLite connection + AutoMigrate
│   ├── models/                 # User and Session GORM models
│   ├── repository/             # User + Session repositories
│   ├── session/                # Session manager
│   └── utils/                  # Helpers
├── data/app.db                 # SQLite database (git-tracked for demo)
├── Dockerfile                  # Multi-stage Go build
├── docker-compose.yml          # Service, volume, env
└── .env                        # Configuration
```
