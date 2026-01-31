# Deployment Instructions for VM/Production

## Method 1: Auto-Seed (Current Default)
The seeder runs automatically on startup and checks if data exists.

### On VM:
```bash
# Build the application
go build -o main cmd/api/main.go

# Run with auto-seed (controlled by env var)
SEED_DATA=true ./main
```

### First Time Setup:
1. App starts
2. Database creates tables
3. Seeder checks if data exists
4. If empty, inserts sample data
5. If has data, skips

## Method 2: Manual Seed Control

### Skip Seeding:
```bash
# Don't set SEED_DATA or set it to false
./main
# OR
SEED_DATA=false ./main
```

### Force Re-Seed:
```bash
# 1. Stop application
# 2. Remove or reset data in your Postgres DB (use psql, Supabase SQL editor, or a migration tool)
#    Example (psql):
#    psql $DATABASE_URL -c "TRUNCATE TABLE products, categories RESTART IDENTITY CASCADE;"

# 3. Start with SEED_DATA=true (ensure DATABASE_URL is set and points to the target DB)
SEED_DATA=true DATABASE_URL="postgres://user:pass@host:5432/dbname" ./main
```

## Docker Deployment

### Dockerfile already configured
The existing `Dockerfile` works as-is.

### Run with Docker:
```bash
# Build image
docker build -t category-api .

# Run with auto-seed (provide DATABASE_URL)
docker run -e DATABASE_URL="postgres://user:pass@host:5432/dbname" -e SEED_DATA=true -p 8080:8080 category-api

# Run without seed (provide DATABASE_URL)
docker run -e DATABASE_URL="postgres://user:pass@host:5432/dbname" -p 8080:8080 category-api
```

---

## Deploying to Supabase (Managed Postgres) ✅

This project can use Supabase as the managed Postgres database. Below are step-by-step instructions to target Supabase for your database (DB).

### 1) Create a Supabase project
- Go to https://app.supabase.com and create a new project (select your organization, project name, region).
- Choose a strong password for the database.

### 2) Get the Database connection string
- In the Supabase dashboard, open **Settings → Database → Connection string**.
- Copy the **Connection string** (it looks like `postgres://<user>:<pass>@<host>:5432/<db>`).
- Ensure the connection string includes `sslmode=require` (or append `?sslmode=require`), e.g.:
  `postgres://user:pass@db.host.supabase.co:5432/postgres?sslmode=require`

### 3) Provide `DATABASE_URL` to the app
- Locally (bash):
```bash
export DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=require"
export SEED_DATA=true   # optional, to seed sample data
go run cmd/api/main.go
```

- Docker (replace with your Supabase connection string):
```bash
docker run -e DATABASE_URL="postgres://user:pass@db.host.supabase.co:5432/postgres?sslmode=require" -e SEED_DATA=true -p 8080:8080 category-api
```

### 4) Seed and Migrations
- On startup the app will AutoMigrate the schema using GORM. If `SEED_DATA=true` the seeder will insert sample data only if tables are empty.
- To force re-seed from a clean DB, truncate tables via Supabase SQL editor or `psql` and start with `SEED_DATA=true`.

### 5) Security & Best Practices
- Do NOT commit the raw `DATABASE_URL` or credentials to source control; use environment variables or secrets manager (GitHub Actions secrets, Docker secrets, etc.).
- Use the Supabase **Service Role Key** only for trusted server-side operations; for normal DB connections use the DB connection string.
- In production, configure monitoring, backups (Supabase offers automatic backups), and restrict network access when applicable.

---

## Deploying to Zeabur (step-by-step)

You can deploy the app to Zeabur using a Docker image (recommended) or by connecting your GitHub repo. Below are concise steps and helpful tips for seeding and runtime configuration.

### 1) Build & publish Docker image (example using Docker Hub or GitHub Container Registry)
```bash
# Build
docker build -t your-username/category-api:latest .
# Push to registry (Docker Hub example)
docker push your-username/category-api:latest
```

### 2) Create an App on Zeabur
- In Zeabur dashboard, create a new **Service** and choose **Container Image**.
- Set the image to `your-username/category-api:latest` (or a registry URL).
- Set the container port to `8080`.

### 3) Set Environment Variables / Secrets
- Add `DATABASE_URL` with your Supabase/Postgres connection (include `?sslmode=require` for Supabase if needed).
- Add optional envs: `SEED_DATA=true` (for seeding), `SEED_EXIT=true` (if you want a one-off that exits after seeding), `GIN_MODE=release`.
- Use Zeabur secrets/variables UI to store these values securely — do not commit them to the repo.

### 4) Run seeder as a one-off job (recommended safe flow)
- Use Zeabur's **Run a Job** / **One-off** feature to execute the seeder without changing the live service.
- Command to run inside container (example):
```
SEED_DATA=true SEED_EXIT=true ./main
```
- This will connect to the configured `DATABASE_URL`, insert sample data if tables are empty, and exit.

### 5) Health checks & monitoring
- Configure a health check endpoint: `/health` (already implemented). Use it for readiness/liveness checks.
- Monitor logs from Zeabur UI to confirm migrations/seeding succeeded and that the service is healthy.

### Notes & Best Practices
- For production, prefer running the seeder as a controlled job rather than enabling `SEED_DATA` on every startup.
- Use Zeabur secrets for credentials and consider network-level protections for your DB.

---

## CI: GitHub Actions (example)

You can add a workflow to build the binary and optionally run the seeder in a controlled CI or manual workflow. The provided example (`.github/workflows/seed.yml`) is configured as a `workflow_dispatch` and expects a repository secret named `DATABASE_URL` that contains your Supabase/Postgres connection string.

- To add the secret: Settings → Secrets → Actions → New repository secret → Name: `DATABASE_URL`, Value: `postgres://user:pass@host:5432/dbname?sslmode=require`.
- To run the seeder manually: Go to Actions → CI - Build & Seed (Supabase) → Run workflow → set `seed=true`.

---

## Docker Compose: Local Supabase/Postgres + App (example)

A lightweight `docker-compose.supabase.yml` example is included to run a local Postgres (compatible with Supabase usage) and the app.

Usage:

```bash
# Build and run the stack
docker compose -f docker-compose.supabase.yml up --build

# Adminer UI available at: http://localhost:8081
# App available at: http://localhost:8080
```

Notes:
- This compose file uses `postgres:15` for local testing. For a full Supabase environment use the Supabase Docker stack or the managed Supabase service.
- The app service will connect to the Postgres service using `DATABASE_URL=postgres://user:pass@db:5432/db?sslmode=disable` and will seed data if `SEED_DATA=true`.

---

## Quick Troubleshooting
- If the app cannot connect: verify `DATABASE_URL`, check that the password is correct, ensure `sslmode=require` (or `sslmode=disable` for local compose) is present, and confirm Supabase project status.
- Use Supabase SQL editor (Dashboard → SQL) or Adminer (local compose) to inspect tables and run manual SQL commands.

## Production Best Practices

### Option 1: Keep Auto-Seed (Small Apps)
- Good for: Development, staging, small apps
- Seeder checks if data exists before inserting
- No duplicate data

### Option 2: Disable Auto-Seed (Large Apps)
- Good for: Production with real data
- Use migrations or manual SQL for production data
- Set `SEED_DATA=false` or don't set it

### Option 3: Use Migration Tools
For production, consider:
- SQL migration files
- Database backup/restore
- Admin panel for data management
