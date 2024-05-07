install:
	go mod tidy

dev:
	DB_HOST=localhost DB_PORT=5432 DB_USER=ak DB_PASSWORD=12345678 DB_NAME=make_app JWT_SECRET=secret go run app.go

test: test-unit test-integration test-e2e

test-unit:
	go test -tags=unit -v ./...

test-coverage:
	go test -cover -tags=unit ./...

	