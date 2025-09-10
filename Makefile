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
		honnef.co/go/tools/cmd/staticcheck \
		github.com/golangci/golangci-lint/cmd/golangci-lint

install-all: install-go install-go-modules install-tools

gqlgen:
	gqlgen generate

go-imports:
	goimports -w $$(find . -name "*.go")

staticcheck:
	staticcheck github.com/Rusty-Alucard/sp_api/...

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test ./... -v

test-coverage:
	go test ./... -v -coverprofile cover.txt -coverpkg=./...
	go tool cover -html=cover.txt -o index.html

server:
	go run server.go wire_gen.go
