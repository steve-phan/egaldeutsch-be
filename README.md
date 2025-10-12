# egalDeutsch Backend

A Go backend service for the egalDeutsch language learning platform, built with Gin and PostgreSQL.

## Features

- RESTful API for articles and quizzes
- PostgreSQL database with migrations
- Docker containerization
- Comprehensive logging
- CORS support
- Health check endpoint

## Tech Stack

- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL
- **ORM**: Raw SQL with `database/sql`
- **Migrations**: golang-migrate
- **Configuration**: Environment variables with godotenv
- **Logging**: Logrus
- **Containerization**: Docker & Docker Compose

## Project Structure

```
egaldeutsch-be/
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Gin middleware
│   ├── models/          # Data models
│   └── server/          # Server setup
├── migrations/          # Database migrations
├── .env.example         # Environment variables template
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile          # Docker build configuration
├── Makefile            # Development commands
├── go.mod              # Go module file
└── README.md           # This file
```

## Quick Start

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

### Development Setup

1. **Clone and setup:**

   ```bash
   git clone <repository-url>
   cd egaldeutsch-be
   make setup
   ```

2. **Start services with Docker:**

   ```bash
   make docker-run
   ```

   Or manually:

   ```bash
   docker-compose up --build
   ```

3. **The API will be available at:** `http://localhost:8080`

### Local Development (without Docker)

1. **Start PostgreSQL:**

   ```bash
   docker run -d \
     --name postgres-dev \
     -e POSTGRES_DB=egaldeutsch \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -p 5432:5432 \
     postgres:15-alpine
   ```

2. **Run migrations:**

   ```bash
   make migrate-up
   ```

3. **Run the application:**
   ```bash
   make run
   ```

## API Endpoints

### Health Check

- `GET /health` - Service health check

### Articles

- `GET /api/v1/articles` - List articles (paginated)
- `GET /api/v1/articles/:id` - Get single article
- `POST /api/v1/articles` - Create new article
- `PUT /api/v1/articles/:id` - Update article
- `DELETE /api/v1/articles/:id` - Delete article

### Query Parameters

- `page` - Page number (default: 1)
- `per_page` - Items per page (default: 10, max: 100)

## Configuration

Copy `.env.example` to `.env` and adjust the values:

```env
SERVER_HOST=localhost
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=egaldeutsch
DB_SSLMODE=disable
```

## Development Commands

```bash
# Build the application
make build

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean

# Docker commands
make docker-build
make docker-run
make docker-stop

# Database migrations
make migrate-up
make migrate-down
make migrate-create name=your_migration_name
```

## Database Schema

### Users

- `id` (UUID) - Primary key
- `name` (VARCHAR) - User name
- `role` (VARCHAR) - 'learner' or 'teacher'
- `created_at`, `updated_at` (TIMESTAMP)

### Articles

- `id` (UUID) - Primary key
- `title` (VARCHAR) - Article title
- `summary` (TEXT) - Article summary
- `content` (TEXT) - Full article content
- `level` (VARCHAR) - Language level (A1, A2, B1, B2, C1)
- `author_id` (UUID) - Foreign key to users
- `created_at`, `updated_at` (TIMESTAMP)

### Quizzes

- `id` (UUID) - Primary key
- `title` (VARCHAR) - Quiz title
- `article_id` (UUID) - Foreign key to articles
- `created_at` (TIMESTAMP)

### Quiz Questions

- `id` (UUID) - Primary key
- `quiz_id` (UUID) - Foreign key to quizzes
- `prompt` (TEXT) - Question text
- `options` (JSONB) - Array of answer options
- `answer` (INTEGER) - Index of correct answer
- `created_at` (TIMESTAMP)

## Testing

Run tests with:

```bash
make test
```

## Deployment

### Docker Production Build

```bash
# Build production image
docker build -t egaldeutsch-be:latest .

# Run with external PostgreSQL
docker run -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_USER=your-db-user \
  -e DB_PASSWORD=your-db-password \
  egaldeutsch-be:latest
```

## Contributing

1. Follow Go best practices and conventions
2. Write tests for new features
3. Update documentation as needed
4. Use meaningful commit messages

## License

[Add your license here]
# egaldeutsch-be
