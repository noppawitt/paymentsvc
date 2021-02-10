include .env
export

test:
	go test ./...

run:
	go run main.go

docker-build:
	docker build -t paymentsvc .
