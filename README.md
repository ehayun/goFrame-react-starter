# Tzlev Application

A modern web application combining GoFrame (Go) backend with an embedded React frontend, integrated with the Wowdash Bootstrap admin template.

## Features

- **Backend**: GoFrame (Go) framework
- **Frontend**: React with Vite
- **Database**: PostgreSQL (via MCP)
- **Cache**: Redis for sessions and caching
- **Authentication**: Google OAuth 2.0
- **Email**: SMTP email service
- **i18n**: Multilingual support (English & Hebrew with RTL)
- **CLI Mode**: Command-line interface for administrative tasks

## Prerequisites

- Go 1.24.2 or higher
- Node.js 16+ and npm
- PostgreSQL 16+
- Redis 7+
- Gmail account for OAuth and SMTP

## Quick Start

### 1. Clone the repository

```bash
git clone <repository-url>
cd tzlev
```

### 2. Configure environment

Copy the example environment file and configure:

```bash
cp .env.example .env
# Edit .env with your actual credentials
```

### 3. Build the application

```bash
./build.sh
```

This will:
- Build template assets from `assets/` to `public/`
- Build the React frontend to `public/`
- Build the Go backend binary to `bin/tzlev`

### 4. Run the application

#### Web Mode (Production)

```bash
./bin/tzlev
```

The server will start on http://localhost:8080

#### Development Mode

```bash
# Terminal 1: Watch and rebuild assets
npm run watch:assets

# Terminal 2: Run backend
go run main.go

# Terminal 3: Run frontend dev server
cd frontend
npm run dev
# Access at http://localhost:3000
```

## CLI Commands

The application supports CLI commands for administrative tasks:

```bash
# Show help
./bin/tzlev help

# Show version
./bin/tzlev version

# Run database migrations
./bin/tzlev migrate --action=up

# Seed database
./bin/tzlev seed
```

## Project Structure

```
tzlev/
├── main.go                 # Application entry point
├── config/                 # Configuration files
├── migrations/             # Database migrations
├── templates/              # Email templates
├── internal/               # Private application code
│   ├── controller/         # HTTP request handlers
│   ├── service/            # Business logic
│   ├── repository/         # Data access layer
│   ├── model/              # Database models
│   ├── middleware/         # HTTP middleware
│   └── cli/                # CLI commands
├── assets/                 # Source template assets
├── frontend/               # React source code
│   └── src/
│       ├── components/     # React components
│       ├── pages/          # Route pages
│       └── locales/        # i18n translations
└── public/                 # Served files (generated)
```

## Development Workflow

### Template Assets

Source files are in `assets/` directory. Always edit files there, not in `public/`:

```bash
# Build assets once
npm run build:assets

# Watch for changes
npm run watch:assets
```

### Frontend Development

```bash
cd frontend
npm run dev
```

This starts Vite dev server on http://localhost:3000 with:
- Hot Module Replacement (HMR)
- Proxy to backend API on http://localhost:8080

### Backend Development

```bash
# Run directly
go run main.go

# Or build and run
go build -o bin/tzlev main.go
./bin/tzlev
```

## Configuration

Configuration is managed via:
1. `config/config.yaml` - Base configuration
2. `.env` - Environment variables (secrets)
3. Environment-specific files (config.development.yaml, etc.)

See `.env.example` for required environment variables.

## Database

### Migrations

Run migrations via CLI:

```bash
# Apply migrations
go run main.go migrate --action=up

# Rollback migrations
go run main.go migrate --action=down
```

### MCP Integration

All database queries use the MCP (Model Context Protocol) for PostgreSQL operations.

## Deployment

### Production Build

```bash
./build.sh
```

### Running in Production

```bash
# Set environment to production
export APP_ENVIRONMENT=production

# Run the binary
./bin/tzlev
```

## License

Copyright © 2024 Tzlev. All rights reserved.
