include .env

test:
	@echo Running all tests
	go test ./...

schema:
	@echo Dumping schema
	docker exec -ti cards_db pg_dump --username=${POSTGRES_USER} --clean --if-exists --schema-only app > schema.sql

load_schema:
	@echo Loading schema
	docker exec -i cards_db psql -U ${POSTGRES_USER} -d app -t < schema.sql