#!/bin/bash
# TPA System Deployment Script
set -e

echo "=== TPA System Deployment ==="

# Load Docker images
echo "Loading Docker images..."
docker load -i tpa-all-images.tar

# Create .env file if not exists
if [ ! -f .env ]; then
    cat > .env << 'ENVEOF'
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tpa
JWT_SECRET=your-secure-jwt-secret-change-this
ENVEOF
    echo "Created .env file - edit with your settings"
fi

# Start services
echo "Starting services..."
docker-compose up -d

echo "Waiting for services..."
sleep 10

docker-compose ps
echo ""
echo "=== Done! Access: http://localhost:8086/tpa ==="
