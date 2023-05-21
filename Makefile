start:
	make start-services db_migrate
	go run cmd/bun/main.go -env=dev api

start-services:
	docker-compose up -d

stop-services:
	docker-compose down

build-image:
	docker build -t prone .

db_reset:
	sudo -u postgres psql -c "DROP DATABASE IF EXISTS prone"
	sudo -u postgres psql -c "CREATE DATABASE prone"

	make db_migrate

db_migrate:
	go run cmd/bun/main.go -env=dev db init
	go run cmd/bun/main.go -env=dev db migrate

test:
	TZ= go test ./org
	TZ= go test ./blog

api_test:
	TZ= go run cmd/bun/main.go -env=test api &
	APIURL=http://localhost:8000/api ./scripts/run-api-tests.sh

default: start
