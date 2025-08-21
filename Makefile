api:
	go run cmd/api/main.go

tests:
	go test -v ./...

.PHONY: api tests
