dev:
	go run ./cmd/angp

test:
	go test ./... -v

test-unit:
	go test ./... -v -run TestUnit

test-integration:
	go test ./... -v -run TestIntegration

cover:
	go test ./... -cover

lint:
	golangci-lint run ./...

vuln:
	govulncheck ./...

build:
	go build -o angp ./cmd/angp
