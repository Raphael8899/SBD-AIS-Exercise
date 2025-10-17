#!/bin/sh

# Stop and remove existing containers if they exist
docker stop orderservice 2>/dev/null || true
docker rm orderservice 2>/dev/null || true
docker stop postgres-db 2>/dev/null || true
docker rm postgres-db 2>/dev/null || true

# Create Docker network if it doesn't exist
docker network create order-network 2>/dev/null || true

# Start PostgreSQL database container
echo "Starting PostgreSQL database..."
docker run -d \
  --name postgres-db \
  --network order-network \
  -e POSTGRES_DB=order \
  -e POSTGRES_USER=docker \
  -e POSTGRES_PASSWORD=docker \
  -e POSTGRES_TCP_PORT=5432 \
  -v postgres-data:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:18

# Wait for database to be ready
echo "Waiting for database to start..."
sleep 5

# Build orderservice Docker image
echo "Building orderservice image..."
docker build -t orderservice .

# Run orderservice container
echo "Starting orderservice..."
docker run -d \
  --name orderservice \
  --network order-network \
  -p 3000:3000 \
  --env-file debug.env \
  orderservice

# Show running containers
echo ""
echo "All containers are running:"
docker ps

echo ""
echo "Application is available at: http://localhost:3000"