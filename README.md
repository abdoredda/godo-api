# Task Manager API

A REST API for managing personal tasks. It supports user accounts, JWT authentication, and task CRUD operations, with each user's tasks kept private.

## Features

- User registration and login
- JWT-protected task endpoints
- Create, list, retrieve, update, and delete tasks
- Filter tasks by completion status, with pagination
- PostgreSQL storage and password hashing

## Tech stack

Go, PostgreSQL, JWT, Docker Compose

## Run locally with Docker

Requirements: Docker and Docker Compose.

1. Create `.env.docker` in the project root:

   ```env
   DATABASE_URL=postgres://godo:password@db:5432/godo?sslmode=disable
   JWT_SECRET=replace-with-a-long-random-secret
   ```

2. Start the API and database:

   ```sh
   docker compose up --build -d
   ```

3. Create the database tables:

   ```sh
   docker compose exec -T db psql -U godo -d godo < schema.sql
   ```

The API is available at `http://localhost:8080`.

## API

| Method | Endpoint | Description | Auth |
| --- | --- | --- | --- |
| `POST` | `/users` | Register an account | No |
| `POST` | `/login` | Log in and receive a JWT | No |
| `GET` | `/users/{id}` | Get a user | No |
| `GET` | `/tasks` | List tasks (`done`, `page`, `limit` filters) | Bearer token |
| `POST` | `/tasks` | Create a task | Bearer token |
| `GET` | `/tasks/{id}` | Get a task | Bearer token |
| `PUT` | `/tasks/{id}` | Update a task | Bearer token |
| `DELETE` | `/tasks/{id}` | Delete a task | Bearer token |

Send the login token in the `Authorization` header as `Bearer <token>`. Task listing defaults to page 1 with 10 results; the maximum page size is 100.

## Run tests
Repository tests require a PostgreSQL test database. Start and initialize the Compose test database:

```sh
docker compose up -d db_test
docker compose exec -T db_test psql -U godo -d godo_test < schema.sql
```

Set `TEST_DATABASE_URL` to `postgres://godo:password@localhost:5433/godo_test?sslmode=disable` in `.env`, then run:

```sh
go test ./...
```
