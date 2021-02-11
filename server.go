package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Rusty-Alucard/sp_api/graph"
	"github.com/Rusty-Alucard/sp_api/graph/generated"
	"github.com/Rusty-Alucard/sp_api/usecase"
	"github.com/Rusty-Alucard/sp_api/infrastructure/persistence"

)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	
	eventPersistence := persistence.NewEventPersistence()
	eventUseCase := usecase.NewEventUseCase(eventPersistence)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{EventUseCase: eventUseCase}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
