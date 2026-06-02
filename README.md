# Expense Tracker API

A personal expense tracking REST API built with **Go** and **Beego v2**.  
Users can register, log in, and manage their daily expenses with filtering, sorting, and spending summaries.  
All data is stored in CSV files.

---

## Tech Stack

| | |
|---|---|
| Language | Go 1.22+ |
| Framework | Beego v2 |
| Storage | CSV files (`encoding/csv`) |
| Docs | Swagger (auto-generated) |

---

## Project Structure

```
expense-tracker-api/
├── conf/
│   └── app.conf              # App configuration
├── controllers/
│   ├── base.go               # Shared response helpers
│   ├── auth.go               # Register, Login
│   ├── expense.go            # Expense CRUD, filters, summary
│   └── health.go             # Health check
├── models/
│   ├── response.go           # Standard response structs
│   ├── user.go               # User structs + request types
│   ├── user_csv.go           # User CSV read/write
│   ├── expense.go            # Expense structs + request types
│   └── expense_csv.go        # Expense CSV read/write
├── services/
│   ├── auth.go               # Registration and login logic
│   ├── auth_test.go          # Auth tests
│   ├── expense.go            # Expense business logic
│   └── expense_test.go       # Expense tests
├── filters/
│   └── auth.go               # X-User-ID auth filter
├── routers/
│   └── router.go             # All route definitions
├── data/                     # CSV files (auto-created at runtime)
├── swagger/                  # Auto-generated API docs
├── main.go
├── go.mod
└── go.sum
```

---

## Setup

### Prerequisites

- Go 1.22 or higher
- Bee CLI tool

```bash
go install github.com/beego/bee/v2@latest
```

Make sure `$GOPATH/bin` is in your `PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Install dependencies

```bash
go mod tidy
```

### Run the server

```bash
bee run
```

Server starts at `http://localhost:8080`

### Generate Swagger docs

```bash
bee generate docs
```

Visit `http://localhost:8080/swagger/` for interactive API documentation.

---

## Running Tests

```bash
go test ./... -cover
```

Run with verbose output:

```bash
go test ./services/... -v -cover
```

---

## API Reference

### Health Check

```bash
curl http://localhost:8080/api/v1/health
```

---

### Auth

#### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123"
  }'
```

Response `201`:
```json
{
  "status": "success",
  "data": null
}
```

#### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
```

Response `200`:
```json
{
  "status": "success",
  "data": {
    "user_id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

---

### Expenses

All expense endpoints require the `X-User-ID` header set to your user ID from login.

#### Create Expense

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "title": "Lunch",
    "amount": 350.50,
    "category": "Food",
    "note": "Team lunch",
    "expense_date": "2025-06-10"
  }'
```

**Valid categories:** `Food`, `Transport`, `Housing`, `Entertainment`, `Shopping`, `Healthcare`, `Education`, `Utilities`, `Other`

#### List Expenses

```bash
# All expenses
curl http://localhost:8080/api/v1/expenses \
  -H "X-User-ID: 1"

# Filter by category
curl "http://localhost:8080/api/v1/expenses?category=Food" \
  -H "X-User-ID: 1"

# Filter by date range
curl "http://localhost:8080/api/v1/expenses?date_from=2025-06-01&date_to=2025-06-30" \
  -H "X-User-ID: 1"

# Sort by amount descending
curl "http://localhost:8080/api/v1/expenses?sort_by=amount&sort_order=desc" \
  -H "X-User-ID: 1"

# Limit results
curl "http://localhost:8080/api/v1/expenses?limit=10" \
  -H "X-User-ID: 1"

# Combined
curl "http://localhost:8080/api/v1/expenses?category=Food&date_from=2025-06-01&sort_by=amount&sort_order=desc&limit=5" \
  -H "X-User-ID: 1"
```

#### Get One Expense

```bash
curl http://localhost:8080/api/v1/expenses/1 \
  -H "X-User-ID: 1"
```

#### Update Expense

```bash
curl -X PUT http://localhost:8080/api/v1/expenses/1 \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "title": "Team Lunch",
    "amount": 400.00
  }'
```

#### Delete Expense

```bash
curl -X DELETE http://localhost:8080/api/v1/expenses/1 \
  -H "X-User-ID: 1"
```

#### Spending Summary

```bash
curl "http://localhost:8080/api/v1/expenses/summary?date_from=2025-06-01&date_to=2025-06-30" \
  -H "X-User-ID: 1"
```

Response:
```json
{
  "status": "success",
  "data": {
    "date_from": "2025-06-01",
    "date_to": "2025-06-30",
    "total_amount": 15400.50,
    "total_count": 3,
    "by_category": [
      { "category": "Housing", "total": 15000.00, "count": 1 },
      { "category": "Food",    "total": 350.50,   "count": 1 },
      { "category": "Transport", "total": 50.00,  "count": 1 }
    ]
  }
}
```

---

## Query Parameters

| Parameter | Type | Description | Example |
|---|---|---|---|
| `category` | string | Filter by category | `?category=Food` |
| `date_from` | YYYY-MM-DD | On or after this date | `?date_from=2025-06-01` |
| `date_to` | YYYY-MM-DD | On or before this date | `?date_to=2025-06-30` |
| `sort_by` | string | `amount` or `expense_date` | `?sort_by=amount` |
| `sort_order` | string | `asc` or `desc` (default: `desc`) | `?sort_order=asc` |
| `limit` | int | Max number of results | `?limit=10` |

---

## Error Responses

All errors follow this format:

```json
{
  "status": "error",
  "message": "Description of the error"
}
```

| Status Code | Meaning |
|---|---|
| `400` | Bad request — invalid input or validation failure |
| `401` | Unauthorized — missing or invalid `X-User-ID` header |
| `404` | Not found — expense does not exist |
| `409` | Conflict — email already registered |
| `500` | Server error |