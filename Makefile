include .env
export

infra-up:
	@docker compose \
		-f ./deploy/docker-compose.yml \
		up -d \
		postgres 

infra-down:
	@docker compose down postgres


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