postgres16bull:
	docker run --name postgres16bull -p 2000:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=femon -d postgres:16.9-bullseye

createdb:
	docker exec -it postgres16bull createdb --username=root --owner=root toribook

dropdb:
	docker exec -it postgres16bull dropdb toribook

migrateup:
	migrate -path db/migration -database "postgresql://root:femon@localhost:2000/toribook?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:femon@localhost:2000/toribook?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:femon@localhost:2000/toribook?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:femon@localhost:2000/toribook?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb dropdb postgres16bull migrateup migratedown sqlc test server migrateup1 migratedown1