#!/bin/sh

# Clean up any existing container
docker stop db 2>/dev/null
docker rm db 2>/dev/null

# Run the MySQL container with platform specification
# Using mysql:8.0 which has better ARM64 support
echo "Starting DB..."
# For ARM64 native performance
docker run --name db -d \
  -e MYSQL_ROOT_PASSWORD=123 \
  -e MYSQL_DATABASE=users \
  -e MYSQL_USER=users_service \
  -e MYSQL_PASSWORD=123 \
  -p 3306:3306 \
  mysql:8.0

# Wait for the database service to start up.
echo "Waiting for DB to start up..."
until docker exec db mysqladmin --silent -uusers_service -p123 ping; do
  echo "Waiting for database connection..."
  sleep 2
done

# Give MySQL a bit more time to fully initialize
sleep 5

# Run the setup script
echo "Setting up initial data..."
docker exec -i db mysql -uusers_service -p123 users < setup.sql

echo "Database setup complete!"