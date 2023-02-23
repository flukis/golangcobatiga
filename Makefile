build:
	@go build -o bin/services

run: build
	./bin/services

test: 
	go test ./..

generatedb:
	sqlc generate

automigrate:
	go run main.go migrate