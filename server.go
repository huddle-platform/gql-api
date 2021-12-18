package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"gitlab.lrz.de/projecthub/gql-api/auth"
	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/resolvers"
	kc "github.com/ory/kratos-client-go"
)

const port = "8080"

func main() {
	kratosConfig := kc.NewConfiguration()
	kratosConfig.Host = "kratos-public"
	kratosConfig.Scheme = "http"
	api := kc.NewAPIClient(kratosConfig)
	resolver, err := resolvers.NewResolver(os.ExpandEnv(os.Getenv("DB_CONNECTION_STRING")))
	if err != nil {
		if os.Getenv("USE_DEMO_RESOLVER") == "true" {
			resolver = resolvers.NewDemoResolver()
		} else {
			log.Fatal(err)
		}
	}
	config := generated.Config{Resolvers: resolver}
	config.Directives.IsLoggedIn = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		_, isLoggedIn := auth.IdentityFromContext(ctx)
		if isLoggedIn {
			return next(ctx)
		} else {
			return nil, fmt.Errorf("authenticate please")
		}
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	http.Handle("/api", playground.Handler("GraphQL playground", "/api/query"))
	corsHandler := cors.New(cors.Options{AllowedOrigins: []string{"http://localhost:3000"}, AllowCredentials: true, OptionsPassthrough: true, AllowedHeaders: []string{"*"}}).Handler(srv)
	//corsHandler := cors.Default().Handler(srv)
	http.HandleFunc("/api/query", func(rw http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("ory_kratos_session")
		if err != nil {
			corsHandler.ServeHTTP(rw, r)
			return
		}
		toSessionReq := api.V0alpha2Api.ToSession(context.Background()).Cookie(cookie.String())
		session, _, err := toSessionReq.Execute()
		if err != nil {
			corsHandler.ServeHTTP(rw, r)
			return
		}

		newContext := auth.NewIdentityContext(r.Context(), &session.Identity)
		newRequest := r.WithContext(newContext)
		corsHandler.ServeHTTP(rw, newRequest)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
