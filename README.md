# Zen Demo - Go (Gin)

> :warning: **SECURITY WARNING**
>
> This is a demonstration application that intentionally contains security vulnerabilities for educational purposes.
> - **DO NOT** run this in production environment
> - **DO NOT** run without proper protection
> - It is strongly recommended to use [Aikido Zen](https://www.aikido.dev/zen) as a security layer

A demo application built with Go and the Gin web framework to demonstrate Aikido security features.

## Setup

1. Clone the repository with submodules:
```bash
git clone --recurse-submodules <repository-url>
```

Or if already cloned:
```bash
git submodule update --init --recursive
```

2. Run the application:
```bash
go run main.go
```

The application will be available at `http://localhost:3000`

## Features

This demo app includes vulnerable endpoints to demonstrate Aikido's security capabilities:

- **Command Injection**: `/api/execute` endpoints
- **SQL Injection**: `/api/pets/:id` and `/api/create` endpoints
- **Path Traversal**: `/api/read` and `/api/read2` endpoints
- **SSRF**: `/api/request` endpoints
- **Rate Limiting**: `/test_ratelimiting_*` endpoints
- **Bot Detection**: `/test_bot_blocking` endpoint
- **User Blocking**: `/test_user_blocking` endpoint

## Project Structure

- `main.go` - Main application with routes and handlers
- `database/` - Database layer with SQLite
- `static/` - Git submodule containing shared UI (HTML, CSS, JS)
