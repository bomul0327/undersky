package main

import (
	"errors"
	"strconv"

	"github.com/graphql-go/graphql"

	us "github.com/hellodhlyn/undersky"
)

var matchType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Match",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"uuid": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"state": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(*us.Match).StateDB, nil
			},
		},
		"game": &graphql.Field{
			Type: graphql.NewNonNull(gameType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var game us.Game
				us.DB.Where(&us.Game{GameID: p.Source.(*us.Match).GameID}).First(&game)
				if game.ID == 0 {
					return nil, errors.New("알 수 없는 오류가 발생했습니다.")
				}
				return &game, nil
			},
		},
		"player": &graphql.Field{
			Type: graphql.NewNonNull(userType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var user us.User
				us.DB.Where(&us.User{ID: p.Source.(*us.Match).PlayerID}).First(&user)
				if user.ID == 0 {
					return nil, errors.New("알 수 없는 오류가 발생했습니다.")
				}
				return &user, nil
			},
		},
		"playerSubmission": &graphql.Field{
			Type: graphql.NewNonNull(submissionType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var sub us.Submission
				us.DB.Where(&us.Submission{ID: p.Source.(*us.Match).PlayerSubmissionID}).First(&sub)
				if sub.ID == 0 {
					return nil, errors.New("알 수 없는 오류가 발생했습니다.")
				}
				return &sub, nil
			},
		},
		"competitor": &graphql.Field{
			Type: graphql.NewNonNull(userType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var user us.User
				us.DB.Where(&us.User{ID: p.Source.(*us.Match).CompetitorID}).First(&user)
				if user.ID == 0 {
					return nil, nil
				}
				return &user, nil
			},
		},
		"competitorSubmission": &graphql.Field{
			Type: graphql.NewNonNull(submissionType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var sub us.Submission
				us.DB.Where(&us.Submission{ID: p.Source.(*us.Match).CompetitorSubmissionID}).First(&sub)
				if sub.ID == 0 {
					return nil, errors.New("알 수 없는 오류가 발생했습니다.")
				}
				return &sub, nil
			},
		},
		"totalRound":    &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"playerWin":     &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"competitorWin": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"result":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"errorMessage":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"createdAt":     &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
		"updatedAt":     &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

var matchQuery = &graphql.Field{
	Type:        matchType,
	Description: "경기 정보를 조회합니다.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var match us.Match
		id, _ := strconv.ParseInt(p.Args["id"].(string), 10, 64)
		us.DB.Where(&us.Match{ID: id}).First(&match)
		if match.ID == 0 {
			return nil, nil
		}
		return &match, nil
	},
}

var matchListQuery = &graphql.Field{
	Type:        listTypeOf(matchType, "MatchList"),
	Description: "경기 정보 목록을 조회합니다.",
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		query := us.DB.Model(&us.Match{})
		switch p.Source.(type) {
		case *us.Submission:
			query = query.Where(&us.Match{PlayerSubmissionID: p.Source.(*us.Submission).ID})
		}

		var matches []us.Match
		query.Order("id desc").Find(&matches)

		var results []interface{}
		for _, m := range matches {
			match := m
			results = append(results, &match)
		}

		return listType{results}, nil
	},
}
