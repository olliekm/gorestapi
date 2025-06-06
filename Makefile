build:
	@go build -o bin/gorestapi cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/gorestapi

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-down:
	@go run cmd/migrate/main.go down

migrate-up:
	@go run cmd/migrate/main.go up
