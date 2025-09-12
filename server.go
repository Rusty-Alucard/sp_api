// Package main provides the entry point for the SP API GraphQL server.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/Rusty-Alucard/sp_api/config"
	"github.com/Rusty-Alucard/sp_api/graph/generated"
	"github.com/Rusty-Alucard/sp_api/graph/resolver"
)

const (
	defaultPort     = "8080"
	readTimeout     = 15 * time.Second
	writeTimeout    = 15 * time.Second
	idleTimeout     = 60 * time.Second
)

func main() {
	env := flag.String("e", "local", "")
	flag.Parse()

	// config初期化
	if err := config.Init(*env); err != nil {
		panic(err)
	}

	port, ret := os.LookupEnv("PORT")
	if !ret {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		EventUseCase: InitializeEventUseCase(),
		TeamUseCase:  InitializeTeamUseCase(),
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// Create HTTP server with timeout configuration
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      nil,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(server.ListenAndServe())
}
