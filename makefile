.PHONY: test
test:
	go test ./...


.PHONY: generate
generate:
	@go run scripts/generate/main.go
