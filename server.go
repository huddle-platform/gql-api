package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/resolvers"
)

const port = "8080"

func main() {
	resolver, err := resolvers.NewResolver(os.ExpandEnv(os.Getenv("DB_CONNECTION_STRING")))
	if err != nil {
		if os.Getenv("USE_DEMO_RESOLVER") == "true" {
			resolver = resolvers.NewDemoResolver()
		} else {
			log.Fatal(err)
		}
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/api", playground.Handler("GraphQL playground", "/api/query"))
	corsHandler := cors.Default().Handler(srv)
	http.Handle("/api/query", corsHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
