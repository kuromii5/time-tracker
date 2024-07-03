# Time Tracker

Time Tracker is a project designed to manage and track worklogs. It uses PostgreSQL as its database, the Chi router for handling routes, and Go for the server and external API. This README provides detailed instructions on setting up and running the project.

## Configuration

The configuration of the project is done through environment variables. Below is an example `.env` file configuration:

```env
ENV=local
DB_URL=postgres://postgres:admin@localhost:5432/time-tracker?sslmode=disable
EXTERNAL_API_PORT=8081
SERVER_PORT=8080
REQ_TIMEOUT=5s
IDLE_TIMEOUT=60s
```

## Setup and Installation

1. Clone the repository

2. Create a .env file

3. Install Go dependencies

4. Run database migrations:

```bash
go run cmd/migrations/main.go --migrate=up
```

To revert migrations, you can use --migrate=down flag

5. Start the external API

```bash
go run cmd/external_api/main.go
```

6. Start the server

```bash
go run cmd/tracker/main.go
```

## Swagger documentation

Swagger documentation is available at the same host as the server. Once the server is running, you can access the documentation in your browser:

```bash
http://localhost:8080/swagger/index.html
```
