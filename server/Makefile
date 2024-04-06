postgresSetup:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=franco0808 -d postgres:15-alpine

postgres:
	docker exec -it postgres15 psql

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root nimble-chat

dropdb:
	docker exec -it postgres15 dropdb nimble-chat

migrateup:
	migrate -path db/migrations -database "postgresql://root:franco0808@localhost:5433/nimble-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:franco0808@localhost:5433/nimble-chat?sslmode=disable" -verbose down

.PHONY: postgresSetup postgres createdb dropdb migrateup migratedown
