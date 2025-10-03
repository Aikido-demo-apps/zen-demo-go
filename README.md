# Zen Demo - Go (Gin)

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

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
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
