# URL Shortener

# 🔗 URL Shortener Service

A **production-ready backend service** built with **Go**, **PostgreSQL**, and **Redis** that allows users to shorten long URLs and redirect them efficiently.

---

## 🚀 Features

1. Shorten long URLs into unique short codes  
2. Redirect short URLs to original links  
3. Prevent duplicate short URLs  
4. Store mappings in **PostgreSQL**  
5. **Redis** caching layer for faster redirects  
6. **Dockerized** for consistent deployment  
7. Configurable via **`.env`** file

---

## 🧰 Tech Stack

| Component | Technology |
|------------|-------------|
| **Language** | Go 1.24+ |
| **Database** | PostgreSQL 16 |
| **Cache** | Redis 7 |
| **Containerization** | Docker & Docker Compose |
| **Migrations** | Goose |
| **Logger** | Zap |
| **Env Loader** | godotenv |
| **ORM** | GORM |

---

## 📁 Folder Structure

```bash
├── config/ # Environment configuration loader
├── database/ # Database connection setup (PostgreSQL)
├── handlers/ # Gin handlers (controllers)
├── middleware/ # Middleware (rate limiting, context,log)
├── migrations/ # SQL migration files
├── models/ # ORM models
├── repository/ # Data access layer (queries)
├── routes/ # API route definitions
├── service/ # Main service logic layer
├── utils/
│ ├── context/ # Custom context with logger and metadata
│ └── cache/ # Redis helper utilities
├── main.go # Application entry point
├── Dockerfile # Docker image setup
├── docker-compose.yml # Container orchestration for app, db, redis
├── .env # Environment configuration file
└── README.md # Project documentation
``` 

---

## ⚙️ Setup Instructions

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/mohan7-code/url-shortener.git
cd url-shortener
``` 

---

### 2️⃣ Create a `.env` File

Create a `.env` file in the root directory and add the following:

```env
# Server Configuration
SERVER_PORT=8080

# Database Configuration
# Use 'db' as host when running via Docker Compose
# Use 'localhost' when running locally without Docker
DATABASE_URL=postgres://<db_user>:<db_password>@db:5432/<db_name>?sslmode=disable

#Db connections
MAX_DB_CONN=20

# Base Short URL 
BASE_SHORT_URL=https://sho.rt

# Redis Configuration
# Use 'redis' for Docker, or 'localhost' for local development
REDIS_URL=redis://redis:6379
```
---

### 3️⃣ Run with Docker Compose

Build and start all services (**App**, **Redis**, **PostgreSQL**):

```bash
docker-compose up --build
``` 

This will:

- Start **PostgreSQL**
- Start **Redis**
- Apply **database migrations** automatically
- Launch the **Go application** on port `8080`

---

## ⚠️ Security Note  

Database credentials are hardcoded in `docker-compose.yml` for simplicity and demo purposes.  
**In production**, always store them securely — for example, using **environment variables**, **Docker secrets**, or a **secret manager** (like AWS Secrets Manager).


## 📡 API Documentation

All APIs are prefixed with `/v1`.

---

### 🔹 1. Shorten a Long URL

**Endpoint:**  
`POST /v1/shorten`

**Description:**  
Takes a long URL and returns a shortened version.  
If the same URL already exists, it returns the same short code (idempotent).

**Request:**
```bash
curl -X POST http://localhost:8080/v1/shorten \
-H "Content-Type: application/json" \
-d '{"original_url": "https://www.example.com/some/very/long/url"}'
```

**Request (with custom alias):**
```bash
curl -X POST http://localhost:8080/v1/shorten \
-H "Content-Type: application/json" \
-d '{"original_url": "https://www.example.com/about","custom_alias": "mybrand"}'
```

**Response:**
```bash
{
    "original_url": "https://www.example.com/some/very/long/url",
    "short_url": "https://sho.rt/Xs50Df1m"
}
```

### 🔹 2. Redirect to Original URL

**Endpoint:**
`GET /v1/:short_code`

**Description:** 
Redirects the user to the original long URL.

**Request:**
```bash
curl -L http://localhost:8080/v1/:short_code
```
**Response:**
```bash
{
    "original_url": "https://www.example.com/some/very/long/url",
    "short_url": "https://sho.rt/Xs50Df1m"
}
```

### 🔹 3. Get All Shortened URLs

**Endpoint:**
`GET /v1/urls?page=1&limit=10`

**Description:** 
Fetches a paginated list of all shortened URLs with metadata like creation date, click count, and last accessed time.

**Request:**
```bash
curl -X GET "http://localhost:8080/v1/urls?page=1&limit=10"
```
**Response:**
```bash
{
    "data": [
        {
            "id": "6de45a29-f9bc-43e3-87f5-fe11dcbcf2fc",
            "short_code": "Uswmtf4a",
            "original_url": "https://github.com",
            "click_count": 7,
            "created_at": "2025-10-30T10:01:56.216281Z",
            "last_accessed_at": "2025-10-30T10:23:16.502812Z"
        },
        {
            "id": "56b36890-88e8-419d-a187-9d0be337519e",
            "short_code": "Xs50Df1m",
            "original_url": "https://www.example.com/some/very/long/url",
            "click_count": 0,
            "created_at": "2025-10-30T09:48:20.60084Z",
            "last_accessed_at": "2025-10-30T09:48:20.599185Z"
        }
    ],
    "total_count": 2,
    "pages": 1
}
```

### 🔹 4. Total Analytics and last accessed time

**Endpoint:**
`GET /v1/analytics/:code`

**Description:** 
It will give the total clics/analytics of a particular url.

**Request:**
```bash
curl -L http://localhost:8080/v1/analytics/:code
```
**Response:**
```bash
{
    "short_code": "Uswmtf4a",
    "original_url": "https://github.com",
    "click_count": 7,
    "last_accessed_at": "2025-10-30T10:23:16.502812Z"
}
```

## 🏗️ Architectural Overview

```text
          ┌────────────────────────────┐
          │        Client / API        │
          │ (Frontend / Postman / CLI) │
          └─────────────┬──────────────┘
                        │
                        ▼
                ┌───────────────────┐
                │      Routes       │
                │ (API Endpoints)   │
                └───────────────────┘
                        │
                        ▼
                ┌───────────────────┐
                │    Middleware     │
                │ (Rate Limit, Ctx) │
                └───────────────────┘
                        │
                        ▼
                ┌───────────────────┐
                │     Handlers      │
                │ (Controllers)     │
                └───────────────────┘
                        │
                        ▼
                ┌───────────────────┐
                │     Service       │
                │ (Business Logic)  │
                └───────────────────┘
                        │
          ┌─────────────┼────────────────┐
          ▼                             ▼
┌────────────────────┐        ┌────────────────────┐
│     Repository     │        │       Cache        │
│ (DB Access Layer)  │◄──────►│   Redis (ShortURL) │
└────────────────────┘        └────────────────────┘
          │
          ▼
┌──────────────────────────────┐
│        PostgreSQL DB         │
│ (URL Mappings + Metadata)    │
│ id | short_code | url | ...  │
└──────────────────────────────┘
```

### ⚙️ Component Responsibilities


| **Layer** | **Purpose** |
|------------|-------------|
| **Routes** | Define REST API endpoints and attach them to handlers. Keeps routing centralized. |
| **Middleware** | Applies **rate limiting**, **logging**, and **context injection** before requests reach handlers. |
| **Handlers (Controller Layer)** | Handle incoming requests, validate data, and call the service layer. |
| **Service Layer** | Core business logic — creates short codes, checks cache, increments clicks, and manages URL lifecycle. |
| **Repository Layer** | Performs all database operations using GORM. |
| **Cache Layer (Redis)** | Caches short → long URL mappings for ultra-fast redirects. |
| **Database (PostgreSQL)** | Persistent data store for URLs, click counts, timestamps. |
| **Config Layer** | Loads environment variables via `.env` using godotenv. |
| **Logger (Zap)** | Structured logging for all layers. |



## 🧩 Design Decisions & Trade-offs


| **Decision** | **Trade-off** |
|---------------|----------------|
| Followed a layered architecture (routes → handlers → service → repository) to maintain clear separation of concerns. | Slightly increases boilerplate, but improves scalability, readability, and testing. |
| Added Redis caching to improve redirect performance and reduce database load. | Requires cache synchronization and adds minor operational complexity. |
| Implemented rate limiting middleware to prevent abuse and ensure fair usage. | Limits reset on restart since it’s in-memory; not distributed. |
| Used structured logging with Zap and request context for observability and traceability. | Slightly increases setup complexity but simplifies debugging in production. |
| Designed URL creation to be idempotent, ensuring the same long URL always maps to a consistent short code. | Requires maintaining consistent hash generation logic and handling collisions. |
| Utilized Docker Compose for a reproducible local setup of app, database, and Redis. | Slightly larger image size and initial setup time. |

