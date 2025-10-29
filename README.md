# URL Shortener
URL Shortener Service

A production-ready backend service built with Go, PostgreSQL, and Redis that allows users to shorten long URLs and redirect them efficiently.

Features

1. Shorten long URLs into unique short codes
2. Redirect short URLs to original links
3.  Prevent duplicate short URLs
4. Store mappings in PostgreSQL
5. Redis caching layer for faster redirects
6. Dockerized for consistent deployment
7. Configurable via .env file


Tech Stack

Language: Go 1.24+
Database: PostgreSQL 16
Cache: Redis 7
Containerization: Docker & Docker Compose
Migrations: Goose
Logger: Zap
Env Loader: godotenv
ORM: GORM


Folder Structure

├── config/               # Environment configuration loader
├── database/             # Database connection setup (PostgreSQL)
├── handlers/             # Gin handlers (controllers)
├── migrations/           # SQL migration files
├── models/               # ORM models 
├── repository/           # Data access layer (queries)
├── routes/               # API route definitions
├── service/              # Business logic layer
├── utils/
│   ├── context/          # Custom context with logger and metadata
│   └── cache/            # Helper cache utilities (if separate from root cache/)
├── main.go               # Application entry point
├── Dockerfile            # Docker image setup
├── docker-compose.yml    # Container orchestration for app, db, redis
├── .env                  # Environment configuration file
└── README.md             # Project documentation


1. Clone the repository

git clone https://github.com/<your-username>/url-shortener.git
cd url-shortener

2️⃣ Create .env file
SERVER_PORT=8080
DATABASE_URL=postgres://postgres:postgres@db:5432/postgres?sslmode=disable
BASE_SHORT_URL=http://localhost:8080
REDIS_URL=redis://redis:6379

3️⃣ Run with Docker Compose
docker-compose up --build


This will:

Start PostgreSQL

Start Redis

Run migrations automatically

Launch the Go app