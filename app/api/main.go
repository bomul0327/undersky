package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var schema graphql.Schema

type simpleResponse struct {
	Result string
}

var simpleResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SimpleResponse",
	Fields: graphql.Fields{
		"result": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
	},
})

func main() {
	schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"signInWithGoogle": signInWithGoogleQuery,
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				"registerUserWithGoogle": registerUserWithGoogleMutation,
			},
		}),
	})

	// Set GraphQL endpoint.
	h := handler.New(&handler.Config{
		Schema:     &schema,
		Playground: true,
	})

	server := http.NewServeMux()
	server.Handle("/graphql", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS configurations.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		h.ContextHandler(context.Background(), w, r)
	})))

	fmt.Println("Server listening port 8000...")
	http.ListenAndServe(":8000", handlers.CompressHandler(server))
}
