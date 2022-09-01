include .env

install: down build build_api up sleep2 load_schema

sleep2:
	@echo "Sleeping for two seconds..."
	sleep 2
	@echo "Yawn... Morning!!"

build_api:
	@echo "Building api container"
	docker compose build api
	@echo "Successfully built"

build:
	@echo Building the app
	env GOOS=linux CGO_ENABLED=0 go build -o ./bin/app_linux ./cmd/api/*.go

up:
	@echo "Starting docker images..."
	docker compose up -d
	@echo "Docker images started!"

stop:
	@echo "Stopping docker compose"
	docker compose stop
	@echo "Docker compose stopped"

down:
	@echo "Removing docker containers"
	docker compose down
	@echo "Docker compose stopped"

run:
	@echo Starting the app
	go build -o ./bin/app ./cmd/api/*.go && ./bin/app

test:
	@echo Running all tests
	APP_ENV=test go test ./...

schema:
	@echo Dumping schema
	docker exec -ti cards_db pg_dump --username=${POSTGRES_USER} --clean --if-exists --schema-only app > schema.sql

load_schema:
	@echo Loading schema
	docker exec -i cards_db psql --username=${POSTGRES_USER} -d app -t < schema.sql