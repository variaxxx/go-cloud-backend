# go-cloud-backend

Backend for a cloud file storage service with auth, folders, file uploads, Kafka-based async processing, PostgreSQL persistence, Prometheus metrics, Grafana dashboards, and a simple load generator.

> [Frontend repository](https://github.com/variaxxx/go-cloud-frontend)

## Stack

- Go
- PostgreSQL
- Kafka
- Prometheus
- Grafana

## Services

- `api` - Main HTTP backend
- `worker` - Kafka consumer that processes uploaded files
- `load` - Simple load generator that registers/logs in and uploads `.txt` files
- `postgres` - Main database
- `kafka` - Event broker
- `prometheus` - Metrics collection
- `grafana` - Metrics visualization
- `kafka-ui` - Kafka UI for local inspection

The load service:
- creates or logs in a test user
- gets an access token
- uploads simple `.txt` files in parallel

## Requirements

- Docker and Docker Compose
- Make

Optional for local non-Docker worker runs:
- `pdftoppm` from `poppler-utils`

## Configuration

Create `.env` in the project root:

```bash
cp .env.example .env
```

## Run With Makefile

### Full stack

Start the main stack:

```bash
make up
```

This starts:
- `postgres`
- `kafka`
- `api`
- `worker`
- `prometheus`
- `grafana`

Stop the same stack:

```bash
make down
```

### Migrations

Create migration:

```bash
make migrate-create seq=create_some_table
```

Apply migrations:

```bash
make migrate-up
```

Rollback migrations:

```bash
make migrate-down
```

## Local URLs

- API: [http://localhost:8000](http://localhost:8000)
- Kafka UI: [http://localhost:8080](http://localhost:8080)
- Prometheus: [http://localhost:9090](http://localhost:9090)
- Grafana: [http://localhost:3030](http://localhost:3030)

Grafana default credentials:
- login: `admin`
- password: `GF_ADMIN_PASSWORD` from `.env` or `admin`