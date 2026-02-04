# Pack Calculator

Order packs calculator application built with Go backend and React frontend, using hexagonal architecture.

## Prerequisites

- Docker
- Docker Compose
- Make

## Setup

```bash
make setup
make build
make up
```

## Usage

Access the application at http://localhost (port 80)

The application is served through nginx reverse proxy:
- Frontend: http://localhost/
- Backend API: http://localhost/api
- Health check: http://localhost/health

### API Endpoints

- `GET /api/pack-sizes` - Get current pack sizes
- `POST /api/pack-sizes` - Update pack sizes
- `POST /api/calculate` - Calculate optimal pack combination

## Architecture

- **Backend**: Go 1.25 with hexagonal architecture
- **Frontend**: React + Vite + TypeScript
- **Database**: PostgreSQL 16 (versioned, append-only pack sizes)
- **Cache**: Redis 7
- **Reverse Proxy**: Nginx (serves frontend and proxies API)

## Development

```bash
# Run tests
make test

# View logs
make logs

# View nginx logs
make nginx-logs

# Stop services
make down

# Clean everything
make clean
```

## Deployment

For droplet deployment:
1. Clone repository
2. Run `make up`
3. Access via droplet IP on port 80
4. For HTTPS, add SSL certificates and update nginx config
