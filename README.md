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
  - **Create posts as an authenticated user or as a guest (no authentication required)**

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
go run main.go
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



## Deployment

To start the server locally, simply use:

```bash
go run main.go
```

## Creating Posts: Authenticated and Guest

You can create posts using either of the following endpoints:

### 1. Authenticated (requires JWT)
- **Endpoint:** `POST /api/posts`
- **Headers:**
  - `Authorization: Bearer <your-jwt-token>`
- **Body:**
```json
{
  "title": "My First Post",
  "content": "Hello **world**!",
}
```

### 2. Guest (no authentication required)
- **Endpoint:** `POST /api/public/posts`
- **Body:**
```json
{
  "title": "Guest Post",
  "content": "Anyone can post!",
  "author": "Anonymous Guest" 
}
```

Both endpoints return the created post, including Markdown and rendered HTML.

## Creating Comments: Authenticated and Guest

You can add comments to posts using any of the following endpoints:

### 1. Authenticated (requires JWT)
- **Endpoint 1:** `POST /api/posts/comments`
  - **Headers:**
    - `Authorization: Bearer <your-jwt-token>`
  - **Body:**
```json
{
  "post_id": 1,
  "content": "Nice post! [Link](https://example.com)",
}
```
- **Endpoint 2:** `POST /api/comments`
  - **Headers:**
    - `Authorization: Bearer <your-jwt-token>`
  - **Body:**
```json
{
  "post_id": 1,
  "content": "Another way to comment as auth user.",
  "author": "raghav"
}
```

### 2. Guest (no authentication required)
- **Endpoint:** `POST /api/public/comments`
- **Body:**
```json
{
  "post_id": 1,
  "content": "Guest comment with *Markdown*.",
  "author": "Anonymous Guest" 
}
```

All endpoints return the created comment, including Markdown and rendered HTML.

## API Testing with Postman

You can test all API endpoints using Postman. Join the shared Postman workspace to access a pre-built folder structure for testing all endpoints:

[Join Postman Team Workspace](https://app.getpostman.com/join-team?invite_code=bbdadb0beac1134d4fa892774d9b78b77abf5bd5dfaf2e4f65c74cf0746437f3&target_code=7301e672379669bf42456c8c8a5d10ba)

This workspace contains organized collections for authentication, posts, and comments APIs, making testing easy and structured.

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

