#!/usr/bin/env sh

# Wait
/usr/bin/wait

# Generate docs
swag init -g v1.go --dir internal/handler/v1,pkg/response,internal/models --output docs/v1 --ot json

# Run migrations
/usr/bin/migrate -path=/app/migrations -database=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable up

# Start application
ENVIRONMENT=${ENVIRONMENT} /app/backend