.PHONY: test
test:
	go test ./...


.PHONY: generate
generate:
	@go run scripts/generate/main.go
	@go fmt $${PWD}/internal/... > /dev/null
	@go fmt $${PWD}/pkg/... > /dev/null
