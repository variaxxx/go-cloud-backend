include .env
export

COMPOSE_FILE=./deploy/docker-compose.yml
DC=docker compose -f $(COMPOSE_FILE)

.PHONY: \
	up down \
	infra-up infra-down \
	api-up api-down api-build api-rebuild \
	worker-up worker-down worker-build worker-rebuild \
	prometheus-up prometheus-down grafana-up grafana-down \
	kafka-ui-up kafka-ui-down \
	logs-api logs-worker logs-prometheus logs-grafana logs-kafka \
	migrate-create migrate-up migrate-down

# Stack
up:
	@$(DC) up -d postgres kafka api worker prometheus grafana

down:
	@$(DC) stop postgres kafka api worker prometheus grafana

# Infra
infra-up:
	@$(DC) up -d postgres kafka

kafka-ui-up:
	@$(DC) up -d kafka-ui

infra-down:
	@$(DC) stop postgres kafka

kafka-ui-down:
	@$(DC) stop kafka-ui

# Observability
prometheus-up:
	@$(DC) up -d prometheus

prometheus-down:
	@$(DC) stop prometheus

grafana-up:
	@$(DC) up -d grafana

grafana-down:
	@$(DC) stop grafana

# API
api-up:
	@$(DC) up -d api

api-build:
	@$(DC) build api

api-rebuild:
	@$(DC) build api
	@$(DC) up -d api

# Worker
worker-up:
	@$(DC) up -d worker

worker-build:
	@$(DC) build worker

worker-rebuild:
	@$(DC) build worker
	@$(DC) up -d worker

# Logs
logs-api:
	@$(DC) logs -f api

logs-worker:
	@$(DC) logs -f worker

logs-prometheus:
	@$(DC) logs -f prometheus

logs-grafana:
	@$(DC) logs -f grafana

logs-kafka:
	@$(DC) logs -f kafka

# Stop services
api-down:
	@$(DC) stop api

worker-down:
	@$(DC) stop worker

# Migrations
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Missing required argument 'seq'"; \
		exit 1; \
	fi; \
	$(DC) run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "${seq}"

migrate-up:
	@$(DC) run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${PG_USER}:${PG_PASSWORD}@postgres:5432/${PG_DB}?sslmode=disable \
		up

migrate-down:
	@$(DC) run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${PG_USER}:${PG_PASSWORD}@postgres:5432/${PG_DB}?sslmode=disable \
		down
