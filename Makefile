.PHONY: setup build up down test migrate-up migrate-down

setup:
	@echo "Setting up project..."

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

test:
	@echo "Running tests..."

migrate-up:
	@echo "Running migrations..."

migrate-down:
	@echo "Rolling back migrations..."
