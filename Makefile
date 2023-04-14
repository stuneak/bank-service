DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

createdb:
	docker-compose exec -it postgres createdb --username=root --owner=root simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock: 
	mockgen -destination db/mock/store.go  -package mock_db github.com/stuneak/simplebank/sqlc_internal Store

.PHONY: createdb dropdb migrateup migratedown new_migration sqlc test server mock