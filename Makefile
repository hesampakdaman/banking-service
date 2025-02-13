.PHONY: build lint run test docker-run docker-build docker-stop docker-test docker-clean

build:
	go build -o banking-service ./cmd/main.go

lint:
	golangci-lint run

run:
	go run ./cmd/main.go

test:
	go test ./...

docker-build:
	docker build -t banking-service .

docker-run:
	docker compose up --build -d

docker-stop:
	docker compose down

docker-test:
	docker run --rm banking-service go test ./...

docker-clean:
	docker compose down --volumes
