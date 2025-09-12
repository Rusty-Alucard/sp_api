export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(shell printenv PATH)

install-go:
	goenv install -s $$(cat .go-version)

install-go-modules:
	go mod tidy

install-tools:
	mkdir -p bin
	go install \
		github.com/99designs/gqlgen \
		golang.org/x/tools/cmd/goimports \
		honnef.co/go/tools/cmd/staticcheck
	go build -o bin/golangci-lint github.com/golangci/golangci-lint/v2/cmd/golangci-lint

install-all: install-go install-go-modules install-tools

gqlgen:
	gqlgen generate

go-imports:
	goimports -w $$(find . -name "*.go")

staticcheck:
	staticcheck github.com/Rusty-Alucard/sp_api/...

lint:
	@if command -v /opt/homebrew/bin/golangci-lint >/dev/null 2>&1; then \
		/opt/homebrew/bin/golangci-lint run; \
	else \
		golangci-lint run; \
	fi

lint-fix:
	@if command -v /opt/homebrew/bin/golangci-lint >/dev/null 2>&1; then \
		/opt/homebrew/bin/golangci-lint run --fix; \
	else \
		golangci-lint run --fix; \
	fi

test:
	go test ./... -v

test-coverage:
	go test ./... -v -coverprofile cover.txt -coverpkg=./...
	go tool cover -html=cover.txt -o index.html

server:
	go run server.go wire_gen.go
