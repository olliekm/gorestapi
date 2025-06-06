build:
	@go build -o bin/gorestapi cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/gorestapi

migrate:
	@migrate create ext- sql -dir $(filter-out $@, $(MAKECMDGOALS))

