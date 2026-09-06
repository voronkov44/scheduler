APP_NAME := scheduler

.PHONY: run migrate build test fmt vet tidy infra-up infra-down infra-logs

run:
	go run ./cmd/api -config config.yaml

migrate:
	go run ./cmd/migrate -config config.yaml

build:
	go build -o ./bin/$(APP_NAME) ./cmd/api

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

tidy:
	go mod tidy

infra-up:
	docker compose up -d

infra-down:
	docker compose down

infra-logs:
	docker compose logs -f