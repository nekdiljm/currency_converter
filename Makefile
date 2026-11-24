include .env
export

docker-up:
	docker compose up -d

docker-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "postgres://currency_converter:123456@localhost:5432/currency_converter?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://currency_converter:123456@localhost:5432/currency_converter?sslmode=disable" down

