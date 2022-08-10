postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=odPass -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root od_reconciliation

dropdb:
	docker exec -it postgres12 dropdb od_reconciliation

migrateup:
	migrate -path db/migration -database "postgresql://root:odPass@localhost:5432/od_reconciliation?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:odPass@localhost:5432/od_reconciliation?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test