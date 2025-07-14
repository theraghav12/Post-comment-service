# Post-Comments API

A RESTful API service for managing posts and comments with user authentication, built with Go, Gin, and PostgreSQL.

## Features

- **User Authentication**
  - JWT-based authentication
  - User registration and login
  - Protected routes

- **Posts**
  - Create, read, update, and delete posts
  - List all posts with pagination
  - Get posts by user

- **Comments**
  - Add comments to posts
  - Rich text support (Markdown)
  - Nested replies

- **Security**
  - Rate limiting
  - Input validation
  - Secure password hashing

- **Developer Experience**
  - Structured logging
  - Configuration via environment variables

## Prerequisites

- Go 1.16+
- PostgreSQL 13+
- Make (optional)

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/post-comments-api.git
cd post-comments-api
```

### 2. Set up environment variables

Create a `.env` file in the root directory:

```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postcomments
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key
```

### 3. Install dependencies

```bash
go mod download
```

### 4. Run database migrations

```bash
make migrate-up
```

### 5. Start the server

```bash
make run
```

The API will be available at `http://localhost:8080`



## Project Structure

```
.
├── cmd/                  # Main application entry point
├── config/              # Configuration management
├── controllers/         # Request handlers
├── middleware/          # Custom middleware
│   ├── auth.go          # Authentication middleware
│   ├── logger.go        # Request logging
│   ├── rate_limit.go    # Rate limiting
│   └── validation.go    # Request validation
├── models/              # Database models
│   ├── comment.go
│   ├── post.go
│   └── user.go
├── pkg/                 # Reusable packages
│   └── markdown/        # Markdown processing
├── routes/              # Route definitions
├── utils/               # Utility functions
│   └── database.go      # Database connection
├── .env.example         # Example environment variables
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
└── README.md           # This file
```

## Environment Variables

| Variable     | Default     | Description                          |
|--------------|-------------|--------------------------------------|
| PORT         | 8080        | Port to run the server on            |
| DB_HOST     | localhost   | PostgreSQL host                      |
| DB_PORT     | 5432        | PostgreSQL port                      |
| DB_USER     | postgres    | PostgreSQL user                      |
| DB_PASSWORD | postgres    | PostgreSQL password                  |
| DB_NAME     | postcomments| Database name                        |
| JWT_SECRET  | -           | Secret key for JWT signing           |

## Testing

Run unit tests:

```bash
make test
```

Run integration tests:

```bash
make test-integration
```

## Linting and Formatting

```bash
make lint      # Run linter
make format    # Format code
```

## Deployment

### Using Docker

1. Build the Docker image:

```bash
docker build -t post-comments-api .
```

2. Run the container:

```bash
docker run -p 8080:8080 --env-file .env post-comments-api
```

### Kubernetes

Example Kubernetes deployment files are provided in the `deploy/` directory.

## Rate Limiting

The API implements rate limiting to prevent abuse:
- 5 requests per second
- Burst of 10 requests
- Per-IP basis

## Error Handling

All error responses follow the same format:

```json
{
  "error": "Error message",
  "details": {
    "field1": "validation error 1",
    "field2": "validation error 2"
  }
}
```

## Contributing

1. Fork the repository
2. Create a new branch
3. Make your changes
4. Run tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
