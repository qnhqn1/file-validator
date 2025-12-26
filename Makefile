.PHONY: tidy run test

tidy:
	go mod tidy

run:
	go run ./cmd/app

test:
	go test ./...