dev:
	go run ./cmd/angp

test:
	go test ./... -v

cover:
	go test ./... -cover

lint:
	golangci-lint run ./...

vuln:
	govulncheck ./...

build:
	go build -o angp ./cmd/angp

release-dry:
	goreleaser release --snapshot --clean
