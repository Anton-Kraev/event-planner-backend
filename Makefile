.PHONY:
.SILENT:

run:
	go run ./cmd/app/main.go

mockgen:
	go generate -x -run=mockgen ./...

test:
	go test ./internal/...
