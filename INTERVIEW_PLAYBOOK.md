# AI-Assisted Interview Playbook & Strategy Guide

This guide is tailored specifically to the evaluation rubric of your interview. It covers **architectural decisions**, **prompting blueprints**, **evaluating AI code**, and **explaining trade-offs**.

---

## 1. Architectural Blueprint (Go + Postgres + React)

A standard, clean architecture that interviewers respect:

```
frontend/ (React + TypeScript + Vite)
   │
   │  HTTP / JSON (REST)
   ▼
backend/ (Go + Gin)
   ├── main.go               # Server bootstrap, router, middleware (CORS, Logger, Recovery)
   ├── db/                   # Connection pool (sql.DB), Ping check, schema migration/init
   ├── models/               # Data structs with JSON & DB tags, validation
   ├── handlers/             # HTTP controller layer (bind JSON, status codes, call store/service)
   └── repository/           # Data access layer (SQL queries using $1 parameters)
         │
         ▼
postgres (Docker, port 5432)
```

### Why this architecture? (How to explain to the interviewer)
- **Separation of Concerns**: HTTP handlers do not write SQL queries. Repository layer isolates database logic. If database driver changes, handlers are untouched.
- **Connection Pooling**: Uses `database/sql` built-in connection pool (`SetMaxOpenConns`, `SetMaxIdleConns`) instead of spawning a new connection per request.
- **Security**: 100% parameterized queries (`$1, $2`) preventing SQL injection.
- **Fail-Fast Typing**: Go structs and TypeScript interfaces ensure end-to-end type safety between backend and frontend.

---

## 2. Effective AI Prompting Strategy (for Cursor / Claude Code / Copilot)

Do **NOT** ask the AI to "build everything at once". Interviewers look for how you break down business problems into modular technical prompts.

### Prompt Phase 1: Database Schema & Connection
> *"Generate a PostgreSQL schema for a [Task/Resource] entity with id (UUID or auto-increment), title (varchar 255 NOT NULL), description (text), status (varchar 50 NOT NULL), created_at (TIMESTAMPTZ DEFAULT NOW()), and updated_at (TIMESTAMPTZ DEFAULT NOW()). Include an index on status. Then write a Go function in `db/db.go` to initialize the database connection pool using `github.com/lib/pq` with safe pool limits (`SetMaxOpenConns(25)`, `SetMaxIdleConns(10)`)."*

### Prompt Phase 2: CRUD Repository Layer
> *"Write a Go repository interface and struct for the [Task] entity implementing: `Create`, `GetAll(status string)`, `GetByID(id int)`, `Update(id int, task)`, and `Delete(id int)`. Use parameterized queries (`$1`, `$2`), ensure `defer rows.Close()` is always called, and handle `sql.ErrNoRows` cleanly."*

### Prompt Phase 3: Gin HTTP Handlers & Routes
> *"Write RESTful Gin handlers for `/api/tasks`: POST (201 Created), GET (200 OK), GET /:id (200 / 404), PUT /:id (200 / 404), and DELETE /:id (204 No Content). Validate incoming JSON payloads using struct binding tags and return structured JSON errors: `{\"error\": \"...\"}`."*

### Prompt Phase 4: Frontend UI (React + Tailwind / Clean CSS)
> *"Create a React TypeScript component in `src/App.tsx` that connects to `http://localhost:8080/api/tasks`. Implement:
> 1. A table/card list displaying tasks with status badges.
> 2. A form/modal to create and edit tasks.
> 3. Delete action with confirmation.
> 4. Search/filter bar by status.
> 5. Handle loading states and display user-friendly error banners if the API fails."*

---

## 3. How to Evaluate & Critique AI-Generated Code

Interviewers specifically evaluate: *"Ability to evaluate the Code written by the agents and knowing what’s correct."*

Here are the **top 6 mistakes AI assistants make** in Go + Postgres + React that you should spot and fix live:

| Issue | What AI writes | What YOU should correct / point out |
| :--- | :--- | :--- |
| **Leaking DB Connections** | `rows, err := db.Query(...)` without closing rows | Always add `defer rows.Close()` immediately after checking `err`, and call `rows.Err()` after iteration. |
| **SQL Injection** | `fmt.Sprintf("SELECT * FROM tasks WHERE status = '%s'", status)` | Replace with parameterized query: `db.Query("SELECT ... WHERE status = $1", status)`. |
| **Missing Pool Limits** | `sql.Open("postgres", connStr)` and nothing else | Add `db.SetMaxOpenConns(25)`, `db.SetMaxIdleConns(10)`, `db.SetConnMaxLifetime(5 * time.Minute)`. |
| **Silent Errors** | `_ = db.Exec(...)` or ignoring returned errors | Always check `if err != nil` and return appropriate HTTP 500 or 400. |
| **CORS Failure** | AI forgets CORS middleware on Gin router | Add `cors.Default()` or configure explicit allowed origins for `http://localhost:5173`. |
| **Stale State in React** | Directly mutating state array: `tasks.push(newTask)` | Use immutable state updates: `setTasks(prev => [...prev, newTask])`. |

---

## 4. Technical Trade-Offs to Articulate

### PostgreSQL vs. Redis
- **PostgreSQL**: Chosen as the primary store for ACID transactions, structured schema validation, relational integrity, and rich querying/filtering.
- **Redis**: Ideal for ephemeral cache, rate-limiting, or session tokens. Not suitable as a primary system of record for complex relational CRUD queries.

### Go Gin vs. Standard Library (`net/http`)
- **Gin**: Provides high-performance radix-tree routing, flexible middleware chaining (CORS, recovery, logging), and effortless JSON binding, drastically reducing boilerplate during a timed interview.
- **Trade-off**: Standard library (`net/http` with Go 1.22+ routing) has zero external dependencies, but requires writing more manual JSON parsing and middleware wrappers.

### REST vs. GraphQL
- **REST**: Lightweight, native HTTP caching, clear CRUD semantics (GET/POST/PUT/DELETE status codes), zero client library complexity for a focused interview task.

---

## 5. Quick Run & Teardown Cheat Sheet

```bash
# 1. Start PostgreSQL (in background)
sudo docker compose up -d

# 2. Run Go Backend (runs on :8080)
cd backend && go run main.go

# 3. Run React Frontend (runs on :5173)
cd frontend && npm run dev

# 4. Yeet / Teardown everything
./scripts/yeet.sh
```
