include .env
export

build: go run cmd/main.go

docker-up:
	docker compose up -d

docker-down:
	docker compose down

migrate-up:
	migrate -path migrations -database ${CONN_DB} up

migrate-down:
	migrate -path migrations -database ${CONN_DB} down

