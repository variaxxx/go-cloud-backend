include .env
export

infra-up:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		up -d \
		postgres \
		kafka

kafka-ui-up:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		up -d \
		kafka-ui

prometheus-up:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		up -d \
		prometheus

api-up:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		up -d \
		api

api-build:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		build \
		api

infra-down:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		stop \
		postgres \
		kafka

kafka-ui-down:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		stop \
		kafka-ui

prometheus-down:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		stop \
		prometheus

api-down:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		stop \
		api


migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Missing required argument 'seq'"; \
		exit 1; \
	fi; \
	docker compose \
		-f ./deploy/docker-compose.yml \
		run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "${seq}"

migrate-up:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${PG_USER}:${PG_PASSWORD}@postgres:5432/${PG_DB}?sslmode=disable \
		up

migrate-down:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${PG_USER}:${PG_PASSWORD}@postgres:5432/${PG_DB}?sslmode=disable \
		down
