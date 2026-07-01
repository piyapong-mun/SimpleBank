runpostgres:
	docker run --name maindb -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mypassword -e POSTGRES_DB=simplebank -d -p 5558:5432 postgres:alpine
createdb:
	docker exec -it maindb createdb -U root -O root simple_bank
dropdb:
	docker exec -it maindb dropdb -U root simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:mypassword@localhost:5558/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:mypassword@localhost:5558/simple_bank?sslmode=disable" -verbose down
migrateup1:
	migrate -path db/migration -database "postgresql://root:mypassword@localhost:5558/simple_bank?sslmode=disable" -verbose up 1
migratedown1:
	migrate -path db/migration -database "postgresql://root:mypassword@localhost:5558/simple_bank?sslmode=disable" -verbose down 1
generate:
	sqlc generate
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/piyapong-mun/simplebank/db/sqlc Store
test:
	go test -v -cover ./...

.PHONY: createdb dropdb runpostgres migrateup migratedown test generate migrateup1 migratedown1 server mock migrateawsup