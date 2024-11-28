DB_URL = postgresql://pasan:12345@localhost:5433/ocr?sslmode=disable

ocr:
	docker run -d --name ocr -p 5433:5432 -e POSTGRES_USER=pasan -e POSTGRES_PASSWORD=12345 postgres:16-alpine

createdb:
	docker exec -it ocr createdb --username=pasan --owner=pasan ocr

dropdb:
	docker exec -it ocr dropdb --username=pasan ocr

migrate:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

dbdocs:
	dbdocs build docs/db.dbml

dbschema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

server:
	go run main.go

redis:
	docker run --name redis -p 6379:6379 -d redis:alpine3.19

.PHONY: postgres createdb migrate dropdb migrateup migratedown sqlc server dbdocs dbschema redis