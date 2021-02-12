gqlgen:
	go run github.com/99designs/gqlgen generate

server:
	go run server.go wire_gen.go