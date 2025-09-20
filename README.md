# CommentTree

**CommentTree** is a service for managing hierarchical (tree-structured) comments, supporting unlimited nesting, search, sorting, and pagination. It provides both a backend API and a simple web frontend for interacting with comment trees.

---

## Features

* Unlimited nested comments
* Create, view, and delete comments
* Full-text search across comments
* Pagination and sorting support
* Simple web interface for browsing, replying, and searching comments

---

## Project Structure

```
.
├── backend/                 # Backend service
│   ├── cmd/                 # Application entry points
│   ├── config/              # Configuration files
│   ├── internal/            # Internal application packages
│   │   ├── api/             # HTTP handlers, router, server
│   │   ├── config/          # Config parsing logic
│   │   ├── middlewares/     # HTTP middlewares
│   │── model/               # Data models
│   │── repository/          # Database repositories
│   │── service/             # Business logic
│   ├── migrations/          # Database migrations
│   ├── Dockerfile           # Backend Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/                # Frontend application (HTML + JS)
├── .env.example             # Example environment variables
├── docker-compose.yml       # Multi-service Docker setup
├── Makefile                 # Development commands
└── README.md
```

---

## API Routes

The backend exposes the following HTTP routes under `/api/comments`:

| Method | Route               | Description                                                                                                                                                                                                                                |
| ------ | ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| POST   | `/api/comments/`    | Create a new comment. Include `parent` field to reply to another comment.                                                                                                                                                                  |
| GET    | `/api/comments/:id` | Retrieve a comment and its full subtree (nested replies).                                                                                                                                                                                  |
| GET    | `/api/comments/`    | Retrieve a list of comments with optional query parameters: <br> `parent={id}` – fetch children of a comment <br> `search={query}` – full-text search <br> `limit={n}` – number of comments per page <br> `offset={n}` – pagination offset |
| DELETE | `/api/comments/:id` | Delete a comment and all its nested replies.                                                                                                                                                                                               |

---

## Development Commands

| Command            | Description                                     |
| ------------------ | ----------------------------------------------- |
| `make docker-up`   | Build and start all Docker services             |
| `make docker-down` | Stop and remove all Docker services and volumes |

---

## Environment Variables

Copy `.env.example` to `.env` and set values.