package main

import "github.com/graphql-go/graphql"

type listType struct {
	Items []interface{}
}

func listTypeOf(t *graphql.Object, name string) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: name,
		Fields: graphql.Fields{
			"items": &graphql.Field{Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(t)))},
		},
	})
}
