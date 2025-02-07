migrateup:
	migrate -path db/migrations -database "postgres://admin:secret@localhost:5432/restaurant_reservation?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgres://admin:secret@localhost:5432/restaurant_reservation?sslmode=disable" -verbose down

test:
	go test -v -cover ./...