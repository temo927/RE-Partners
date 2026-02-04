.PHONY: setup build up down test migrate-up migrate-down clean logs

setup:
	@echo "Setting up project..."
	@echo "Creating .env file from .env.example if it doesn't exist..."
	@test -f backend/.env || cp backend/.env.example backend/.env
	@echo "Setup complete. Please review backend/.env and adjust if needed."

build:
	@echo "Building Docker images..."
	docker-compose build

up:
	@echo "Starting services..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Running migrations..."
	@$(MAKE) migrate-up

down:
	@echo "Stopping services..."
	docker-compose down

test:
	@echo "Running backend tests..."
	cd backend && go test ./... -v

migrate-up:
	@echo "Running migrations..."
	@cat backend/internal/adapters/repository/migrations/000001_create_pack_sizes.up.sql | docker-compose exec -T postgres psql -U packcalc -d packcalc

migrate-down:
	@echo "Rolling back migrations..."
	@cat backend/internal/adapters/repository/migrations/000001_create_pack_sizes.down.sql | docker-compose exec -T postgres psql -U packcalc -d packcalc

clean:
	@echo "Cleaning up..."
	docker-compose down -v
	@echo "Removed volumes and containers"

logs:
	@echo "Showing logs..."
	docker-compose logs -f
