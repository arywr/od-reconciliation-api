postgres:
	docker run --name postgres11 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=odPass -d postgres:11-alpine

createdb:
	docker exec -it postgres11 createdb --username=root --owner=root od_reconciliation

dropdb:
	docker exec -it postgres11 dropdb od_reconciliation

migrateup:
	migrate -path db/migration -database "postgresql://root:odPass@localhost:5432/od_reconciliation?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:odPass@localhost:5432/od_reconciliation?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

start:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test start