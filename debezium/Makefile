SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command

run-server:
ifeq ($(OS),Windows_NT)
	$$env:ENV_PATH='config/.env'; go run cmd/debezium/main.go
else
	ENV_PATH=./config/.env go run cmd/debezium/main.go
endif

run-server-local:
ifeq ($(OS),Windows_NT)
	$$env:ENV_PATH='config/.env.local'; go run cmd/debezium/main.go
else
	ENV_PATH=./config/.env.local go run cmd/debezium/main.go
endif

docker-up:
ifeq ($(OS),Windows_NT)
	docker-compose --env-file config/.env.local up -d
else
	docker-compose --env-file config/.env.local up -d
endif

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f db

migrate-up:
ifeq ($(OS),Windows_NT)
	$$env:ENV_PATH='config/.env.local'; go run cmd/migrate/main.go -command up
else
	ENV_PATH=./config/.env.local go run ./cmd/migrate/main.go -command up
endif

migrate-down:
ifeq ($(OS),Windows_NT)
	$$env:ENV_PATH='config/.env.local'; go run cmd/migrate/main.go -command down
else
	ENV_PATH=./config/.env.local go run ./cmd/migrate/main.go -command down
endif

migrate-version:
ifeq ($(OS),Windows_NT)
	$$env:ENV_PATH='config/.env.local'; go run cmd/migrate/main.go -command version
else
	ENV_PATH=./config/.env.local go run ./cmd/migrate/main.go -command version
endif

local-env:
ifeq ($(OS),Windows_NT)
	docker-compose --env-file config/.env.local up -d
	$$env:ENV_PATH='config/.env.local'; go run cmd/migrate/main.go -command up
	$$env:ENV_PATH='config/.env.local'; go run cmd/debezium/main.go
else
	@echo "Starting PostgreSQL in Docker..." && \
	docker-compose --env-file config/.env.local up -d && \
	echo "Waiting for PostgreSQL to be ready..." && \
	sleep 5 && \
	echo "Running migrations..." && \
	ENV_PATH=./config/.env.local go run ./cmd/migrate/main.go -command up
endif