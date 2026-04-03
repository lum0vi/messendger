run:
	swag init --generalInfo cmd/main.go --output docs
	go run cmd/main.go

migrateup:
    migrate -path db/migration -database "postgresql://user:secret@postgres:5432/postgres?sslmode=disable" -verbose up

migratedown:
    migrate -path db/migration -database "postgresql://user:secret@postgres:5432/postgres?sslmode=disable" -verbose down
