.PHONY:
.SILENT:

run:
	docker-compose -f docker-compose.yml up --build -d

run-redis:
	docker-compose -f docker-compose.yml up --build -d redis

run-app:
	go run ./cmd/app/main.go

mockgen:
	go generate -x -run=mockgen ./internal/...

test:
	go test ./internal/...
