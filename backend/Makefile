BINARY_API    := bin/api
BINARY_WORKER := bin/worker
MIGRATE       := migrate
DB_URL        := postgres://cc:cc_dev_pw@localhost:5432/callcenter?sslmode=disable
MIGRATIONS    := ./migrations
COMPOSE       := docker compose -f infra/docker-compose.dev.yml

.PHONY: all up down migrate rollback generate build run worker test test-short lint seed smoke frontend-install frontend-dev frontend-build

all: build

## Infrastructure
up:
	$(COMPOSE) up -d
	@echo "Waiting for services to be healthy..."
	@$(COMPOSE) ps

down:
	$(COMPOSE) down

## Migrations
migrate:
	$(MIGRATE) -path $(MIGRATIONS) -database "$(DB_URL)" up

rollback:
	$(MIGRATE) -path $(MIGRATIONS) -database "$(DB_URL)" down 1

migrate-status:
	$(MIGRATE) -path $(MIGRATIONS) -database "$(DB_URL)" version

## Code generation
generate:
	sqlc generate -f sqlc/sqlc.yaml

## Build
build:
	go build -o $(BINARY_API)    ./cmd/api
	go build -o $(BINARY_WORKER) ./cmd/worker

## Run (native, reads .env)
run:
	@export $$(grep -v '^#' .env | xargs) && go run ./cmd/api

worker:
	@export $$(grep -v '^#' .env | xargs) && go run ./cmd/worker

## Test
test:
	go test ./... -race -count=1 -timeout=180s

test-short:
	go test ./... -short -race -count=1 -timeout=60s

## Seed dev data
seed:
	PGPASSWORD=cc_dev_pw psql -h localhost -U cc -d callcenter -f scripts/seed.sql

## Smoke test (requires running server)
smoke:
	@./scripts/smoke.sh

## Lint
lint:
	golangci-lint run ./...

## Frontend
frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build
