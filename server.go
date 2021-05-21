package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Rusty-Alucard/sp_api/config"
	"github.com/Rusty-Alucard/sp_api/graph/generated"
	"github.com/Rusty-Alucard/sp_api/graph/resolver"
	"github.com/Rusty-Alucard/sp_api/infrastructure/provider"
)

const defaultPort = "8080"

func main() {

	env := flag.String("e", "local", "")
	flag.Parse()

	// config初期化
	if err := config.Init(*env); err != nil {
		panic(err)
	}

	// db初期化
	if err := provider.Init(); err != nil {
		panic(err)
	}

	port, ret := os.LookupEnv("PORT")
	if ret == false {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		EventUseCase: InitializeEventUseCase(),
		TeamUseCase:  InitializeTeamUseCase(),
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
