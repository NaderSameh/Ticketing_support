postgres:
	docker run --name postgres8 -p 5432:5432 -e POSTGRES_USER=nader -e POSTGRES_PASSWORD=nader123 -e POSTGRES_DB=ticketing_support -d postgres:latest

migrateup:
	migrate -path db/migration -database "postgresql://nader:nader123@localhost:5432/ticketing_support?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://nader:nader123@localhost:5432/ticketing_support?sslmode=disable" -verbose down

sqlc:
	sqlc generate

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/naderSameh/ticketing_support/db/sqlc Store

test:
	go test -cover ./...

swag:
	swag init --parseDependency  --parseInternal -g main.go

.PHONY:
	migrateup migratedown postgres sqlc mock swag