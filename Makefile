
LINTER_VERSION=v1.42.1

.PHONY: build clean run

# Runs the projects docker containers
run:
	go run simple-kv-store/cli/main.go

# Runs the go testing stage
test:
	go test ./...

# Runs the go testing stage checking for race conditions.
test-race:
	go test ./... -race

get-linter:
	command -v golangci-lint || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh  | sh -s --  -b ${GOPATH}/bin ${LINTER_VERSION}
# Runs the Go linter
lint: get-linter
	golangci-lint run