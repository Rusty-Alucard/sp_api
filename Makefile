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
		honnef.co/go/tools/cmd/staticcheck

install-all: install-go install-go-modules install-tools

gqlgen:
	gqlgen generate

staticcheck:
	staticcheck github.com/Rusty-Alucard/sp_api/...

server:
	go run server.go wire_gen.go
