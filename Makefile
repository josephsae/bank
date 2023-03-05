postgres:
	sudo docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15rc1-alpine
createdb:
	sudo docker exec -it postgres15 createdb --username=root --owner=root bank
dropdb:
	sudo docker exec -it postgres15 dropdb bank	
migrateup:
	sudo $(GOBIN)/migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up
migratedown:
	sudo $(GOBIN)/migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
