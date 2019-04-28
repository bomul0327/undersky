package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/hellodhlyn/sqstask"

	"github.com/hellodhlyn/undersky/libs/s3"
)

type ctxKey int8

const (
	ctxUser ctxKey = iota
)

var (
	schema         graphql.Schema
	s3client       *s3.S3Client
	submissionTask *sqstask.SQSTask
)

// simpleResponse 는 특별한 반환 타입이 없는 mutation API의 응답으로 사용되는 타입입니다.
type simpleResponse struct {
	Result string
}

var simpleResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SimpleResponse",
	Fields: graphql.Fields{
		"result": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
	},
})

func initS3() {
	cli, err := s3.NewClient("undersky-ai")
	if err != nil {
		panic(err)
	}

	s3client = cli
}

func initSQS() {
	task, err := sqstask.NewSQSTask(&sqstask.Options{
		QueueName:   "undersky-submission",
		AWSRegion:   "ap-northeast-2",
		Consume:     nil,
		HandleError: func(err error) { fmt.Println(err) },
	})
	if err != nil {
		panic(err)
	}

	submissionTask = task
}

func main() {
	initS3()
	initSQS()

	// GraphQL 스키마 설정
	schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"signInWithGoogle": signInWithGoogleQuery,
				"me":               meQuery,
				"game":             gameQuery,
				"gameList":         gameListQuery,
				"match":            matchQuery,
				"matchList":        matchListQuery,
				"submission":       submissionQuery,
				"submissionList":   submissionListQuery,
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				"registerUserWithGoogle": registerUserWithGoogleMutation,
				"submitSource":           submitSourceMutation,
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

		// Make up context.
		ctx := context.Background()
		if authorization := r.Header.Get("Authorization"); strings.HasPrefix(authorization, "Bearer ") {
			user, err := getUserFromJWTToken(strings.Split(authorization, " ")[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, ctxUser, user)
		}

		h.ContextHandler(ctx, w, r)
	})))

	fmt.Println("Server listening port 8000...")
	http.ListenAndServe(":8000", handlers.CompressHandler(server))
}
