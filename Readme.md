# Medistream   
*A backend system for a modern Health Management System (HMS)*  

---

##  Overview  
Medistream is a backend service for managing healthcare operations such as:  

- User & Role Management (Doctors, Patients, Admins)  
- Medical Records & Prescriptions  
- Vitals & Reports  
- Appointments & Scheduling  
- Authentication & Authorization (JWT-based)  
- Rate Limiting & Request Metrics (Prometheus + Redis)  

Built with **Go + Gin + PostgreSQL + Redis**, Medistream is structured for scalability, modularity, and performance.  

---

## Tech Stack  

- **Language:** Go 1.22+  
- **Framework:** [Gin](https://github.com/gin-gonic/gin)  
- **Database:** PostgreSQL + [GORM](https://gorm.io/)  
- **Cache:** Redis  
- **Metrics:** Prometheus  
- **Auth:** JWT (access & refresh tokens)  
- **Migrations:** Raw SQL with migration scripts using goose

---

## Project Structure  

```bash
.
├── cmd/                # Application entrypoint
├── config/             # Database & Redis configuration
├── controllers/        # Request handlers (appointments, auth, vitals, etc.)
├── metrics/            # Prometheus metrics setup
├── middleware/         # Auth, rate limiting, metrics middleware
├── migrations/         # SQL migration files
├── models/             # GORM models (User, Patient, Doctor, Vitals, etc.)
├── routes/             # Route registration
├── scripts/            # Utility scripts (migrations, etc.)
├── services/           # Business logic (cache, ratelimiter, etc.)
├── tests/              # API tests, factories, helpers
└── utils/              # Helpers (JWT, logger, context utils, etc.)

```
## Getting Started


- Go 1.22+
- PostgreSQL 14+
- Redis 7+
- Goose (use the provided migration script in scripts/migrate.sh)


## Setup

```bash
git clone https://github.com/<your-username>/medistream.git
cd medistream

```

## Install Dependencies
```bash
go mod tidy

```

## Setup environment variables

```bash
DATABASE_URL=postgres://user:password@localhost:5432/medistream?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=supersecret
PORT=8080

```

## Run migrations

```bash
./scripts/migrate.sh up

```


## Run migrations

```bash
./scripts/migrate.sh up
```

## Run the Server

```bash
go run cmd/main.go

```

## Testing 

```bash
go test ./tests/... -v
```

## Metrics and Monitoring

- Exposes Prometheus metrics under /metrics.
- Tracks DB queries,request timings,cache hits and misses

## Authentication
- JWT Access & Refresh tokens
- RBAC (DOCTOR, PATIENT, ADMIN)

## Future Plans
- [ ] Add support for notifications (email/SMS)
- [ ] Expand reporting & analytics
- [ ] Containerize with Docker & Helm (K8s support)
- [ ] CI/CD pipeline integration


## Contributing

- Contributions are welcome! Please fork the repo, create a branch, and submit a PR.