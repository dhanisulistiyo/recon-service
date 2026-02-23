# Reconciliation Service

Transaction reconciliation microservice for comparing internal system transactions against multiple bank statement files.

Built with:

- Go
- Echo Framework

---

# Features

- Upload system transaction CSV
- Upload multiple bank CSV files
- Asynchronous reconciliation processing
- Job status tracking endpoint
- Reconciliation summary (matched, unmatched, discrepancies)
- In-memory job repository

---

# Project Structure

```
reconciliation-service/
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── di/
│   ├── domain/
│   ├── interface/http/
│   ├── usecase/
│   ├── infrastructure/
│   └── shared/
│
├── scripts/
│   └── docker.sh
│
├── go.mod
├── .env
└── Dockerfile
```

---

# Environment Configuration

Create `.env` file in root:

```
APPLICATION_NAME=reconciliation
APPLICATION_HTTPPORT=8081
```

---

# Run Locally (Without Docker)

## Install dependencies

```
go mod tidy
```

## Run service

```
go run ./cmd
```

Service will start at:

```
http://localhost:8081
```

---

# Run Using Docker

## Option 1 — Manual Docker Commands

### Build image

```
docker build -t reconciliation-service:latest .
```

### Run container

```
docker run -d \
  --name reconciliation-service \
  -p 8081:8081 \
  --env-file .env \
  reconciliation-service:latest
```

---

## Option 2 — Using Script (Recommended)

Make sure script is executable:

```
chmod +x docker.sh
```

### Build image

```
./docker.sh build
```

### Run container

```
./docker.sh run
```

### Build + Run

```
./docker.sh up
```

---

# API Endpoints

---

## Submit Reconciliation Job

**POST** `/reconcile`

### Form-Data Parameters

| Key          | Type | Description |
|-------------|------|------------|
| system_file | file | System transaction CSV |
| bank_files  | file | Multiple bank CSV files |
| start_date  | text | Format: YYYY-MM-DD |
| end_date    | text | Format: YYYY-MM-DD |

---

### Example Curl

```
curl -X POST http://localhost:8081/reconcile \
  -F "system_file=@system.csv" \
  -F "bank_files=@bca.csv" \
  -F "bank_files=@mandiri.csv" \
  -F "start_date=2024-01-01" \
  -F "end_date=2024-01-31"
```

Response:

```
{
  "job_id": "{{uuid}}",
  "status": "PROCESSING"
}
```

---

## Get Job Status

**GET** `/jobs/:id`

Example:

```
curl http://localhost:8081/jobs/{job_id}
```

Response:

```
{
  "id": "uuid",
  "status": "DONE",
  "result": { ...summary... },
  "error": ""
}
```

---

# Reconciliation Summary Output

Summary contains:

- TotalProcessed
- TotalMatched
- TotalUnmatched
- TotalDiscrepancy
- UnmatchedSystem
- UnmatchedBankByBank (grouped by bank)

---

# Architecture

This project follows Clean Architecture principles:

- **Domain Layer** → Entities & interfaces
- **Usecase Layer** → Business logic
- **Interface Layer** → HTTP handlers
- **Infrastructure Layer** → Repository implementation
- **DI Layer** → Dependency injection wiring

Async job pattern prevents HTTP timeout for large CSV files (500k+ records).

---
