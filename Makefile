gqlgen:
	go run github.com/99designs/gqlgen generate

server:
	go run server.go wire_gen.go

staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck github.com/Rusty-Alucard/sp_api/...
