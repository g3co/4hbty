# Rakia Technical Assessment

This repository contains solutions for both the logical test and technical test components of the Rakia technical assessment.

## Project Structure

```
.
├── logical_test/     # Solution for the message decoding problem
├── pkg/             # REST API implementation
│   ├── models/      # Data models
│   ├── rest/        # HTTP handlers and routing
│   ├── service/     # Business logic
│   └── store/       # Data storage
├── blog_data.json   # Sample blog data
├── go.mod          # Go module file
└── main.go         # Main application entry point
```

## Part 1: Message Decoding Problem

### Problem Description
Given a string of digits, determine how many ways it can be decoded using the following mapping:
- 'A' -> 1
- 'B' -> 2
- ...
- 'Z' -> 26

### Examples
- Input: "12"
  - Possible decodings: "AB" (1,2) and "L" (12)
  - Output: 2

- Input: "226"
  - Possible decodings: "BZ" (2,26), "VF" (22,6), and "BBF" (2,2,6)
  - Output: 3

- Input: "0"
  - No valid decodings
  - Output: 0

### Running the Logical Test
```bash
cd logical_test
go test -v
```

## Part 2: Blog Platform REST API

### Features
- CRUD operations for blog posts
- In-memory data storage
- RESTful API design
- Error handling
- Unit tests
- Graceful shutdown
- Preloaded sample data

### API Endpoints

1. Get All Posts
```bash
GET /posts
```

2. Get Post by ID
```bash
GET /posts/{id}
```

3. Create New Post
```bash
POST /posts
Content-Type: application/json

{
    "title": "New Post",
    "content": "Post content",
    "author": "Author Name"
}
```

4. Update Post
```bash
PUT /posts/{id}
Content-Type: application/json

{
    "title": "Updated Title",
    "content": "Updated content",
    "author": "Author Name"
}
```

5. Delete Post
```bash
DELETE /posts/{id}
```

### Running the REST API

1. Install dependencies:
```bash
go mod tidy
```

2. Start the server:
```bash
go run main.go
```

The server will start on port 8088 and automatically load the sample blog data from `blog_data.json`.

### Testing the API

You can test the API using curl or any HTTP client. Here are some example requests:

```bash
# Get all posts
curl http://localhost:8088/posts

# Get a specific post
curl http://localhost:8088/posts/1

# Create a new post
curl -X POST http://localhost:8088/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"New Post","content":"New Content","author":"New Author"}'

# Update a post
curl -X PUT http://localhost:8088/posts/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","content":"Updated Content","author":"Updated Author"}'

# Delete a post
curl -X DELETE http://localhost:8088/posts/1
```

### Running Tests

To run all tests:
```bash
go test ./...
```

## Requirements

- Go 1.21 or later
- curl (for testing the API)

## Dependencies

- github.com/gorilla/mux - HTTP router
- github.com/google/uuid - UUID generation
- golang.org/x/sync/errgroup - Error handling and goroutine management

## Error Handling

The API includes error handling for common scenarios:
- Invalid request payload
- Post not found
- Validation errors
- Server errors

## Graceful Shutdown

The server implements graceful shutdown, handling SIGTERM and SIGINT signals to ensure all requests are completed before shutting down. 