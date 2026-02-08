# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main cmd/api/main.go

# Note: Database connection is handled at runtime via DATABASE_URL (Supabase/Postgres).
# To seed the remote database use SEED_DATA=true at runtime or in a CI job that has DATABASE_URL set.

# Runtime Stage
FROM alpine:latest  

# Install tzdata for correct time handling
RUN apk add --no-cache tzdata

WORKDIR /app

# Copy the Pre-built binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/web ./web
# Database is external (Supabase/Postgres). Ensure DATABASE_URL is provided to the container at runtime.

# Expose port 8080 to the outside world
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release

# Command to run the executable
CMD ["./main"]
