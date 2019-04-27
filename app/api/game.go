package main

import (
	"github.com/graphql-go/graphql"

	us "github.com/hellodhlyn/undersky"
)

var gameType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Game",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(*us.Game).GameID, nil
			},
		},
		"title":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"summary":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"description": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"createdAt":   &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
		"updatedAt":   &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

var gameQuery = &graphql.Field{
	Type:        gameType,
	Description: "게임 정보를 조회합니다.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var game us.Game
		us.DB.Where(&us.Game{GameID: p.Args["id"].(string)}).First(&game)
		if game.ID == 0 {
			return nil, nil
		}
		return &game, nil
	},
}

var gameListQuery = &graphql.Field{
	Type:        listTypeOf(gameType, "GameList"),
	Description: "게임 목록을 조회합니다.",
	Args:        graphql.FieldConfigArgument{},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var games []us.Game
		us.DB.Model(&us.Game{}).Order("game_id asc").Find(&games)

		var results []interface{}
		for _, m := range games {
			results = append(results, &m)
		}

		return listType{results}, nil
	},
}
