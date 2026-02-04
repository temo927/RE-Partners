# Pack Calculator

Order packs calculator application built with Go backend and React frontend.

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

Access the frontend at http://localhost:3000

API endpoints:
- `GET /api/pack-sizes` - Get current pack sizes
- `POST /api/pack-sizes` - Update pack sizes
- `POST /api/calculate` - Calculate optimal pack combination
