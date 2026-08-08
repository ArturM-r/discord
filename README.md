# Discord Clone Backend

A Go backend service for a lightweight Discord-like application. It provides user authentication, server and channel management, member membership control, message handling, and a websocket endpoint for real-time communication.

## Project Structure

- `cmd/main.go` - application entry point, graceful shutdown, database connection, and route initialization.
- `config/` - application configuration, database connection, and migration runner.
- `internal/` - domain modules for users, servers, channels, members, messages, and membership checks.
- `ws/` - websocket hub and real-time message broadcasting.
- `migrations/` - SQL migrations for PostgreSQL schema.

## Features

- User registration and login
- PostgreSQL database connection using `pgxpool`
- Auto-run database migrations with `golang-migrate`
- Server CRUD and user server membership management
- Channel creation and deletion
- Message retrieval per channel
- Websocket endpoint for real-time updates
- Graceful shutdown handling

## Requirements

- Go 1.26
- PostgreSQL database
- `.env` file or environment variables configured

## Environment Variables

The app reads the following environment values:

- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string (present in config, not currently used by routes)
- `HMAC_KEY` - secret used for authentication / token signing

## Running Locally

1. Copy `.env.example` to `.env` and update the values for your environment, or export the variables in your shell.

Example `.env`:

```env
DATABASE_URL=postgres://user:password@localhost:5432/discord_db?sslmode=disable
REDIS_URL=redis://localhost:6379
HMAC_KEY=some-secret-key
```

2. Run the app:

```bash
go run ./cmd
```

3. The server listens on `http://localhost:8080`.

## API Endpoints

### Authentication

- `POST /auth/register` - register a new user
- `POST /auth/login` - login and receive auth credentials

### Servers

- `GET /servers` - list user servers
- `POST /servers` - create a new server
- `GET /servers/{id}` - get server details
- `DELETE /servers/{id}` - delete a server

### Members

- `POST /servers/{id}/members` - add a member to a server
- `DELETE /servers/{id}/members/{uid}` - remove a member from a server

### Channels

- `GET /servers/{id}/channels` - list channels for a server
- `POST /servers/{id}/channels` - create a new channel
- `DELETE /servers/{id}/channels/{cid}` - delete a channel

### Messages

- `GET /channels/{id}/messages` - get messages for a channel

### Websocket

- `GET /ws` - websocket upgrade endpoint for real-time messaging

## Database Migrations

Migrations are stored in `migrations/` and run automatically during startup.

- `000001_create_users.up.sql`
- `000002_create_servers.up.sql`
- `000003_create_channels.up.sql`
- `000004_create_messages.up.sql`
- `000005_create_members.up.sql`

## Notes

- The project uses `github.com/google/uuid` for UUID handling.
- The websocket hub broadcasts messages to connected clients by channel.
- The app currently uses a `http.ServeMux` with path placeholders like `/servers/{id}`; client routing logic may need to parse IDs from request paths.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Open a pull request with a clear summary

## License

This repository does not include a license file. Add one if you want to make the project open source.
