package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/graphql-go/graphql"

	us "github.com/hellodhlyn/undersky"
)

var submissionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Submission",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"uuid": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"user": &graphql.Field{
			Type: graphql.NewNonNull(userType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var user us.User
				us.DB.Where(&us.User{ID: p.Source.(*us.Submission).UserID}).First(&user)
				if user.ID == 0 {
					return nil, errors.New("알 수 없는 오류가 발생했습니다.")
				}
				return &user, nil
			},
		},
		"description": &graphql.Field{Type: graphql.String},
		"game": &graphql.Field{
			Type: graphql.NewNonNull(gameType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var game us.Game
				us.DB.Where(&us.Game{GameID: p.Source.(*us.Submission).GameID}).First(&game)
				if game.ID == 0 {
					return nil, errors.New("알 수 없는 오류가 발생했습니다.")
				}
				return &game, nil
			},
		},
		"runtime":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"createdAt": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
		"updatedAt": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
		// 순환참조 문제로, `matchList` 필드는 init() 메소드에 선언됩니다.
	},
})

var submissionQuery = &graphql.Field{
	Type:        submissionType,
	Description: "제출 정보를 조회합니다.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var sub us.Submission
		id, _ := strconv.ParseInt(p.Args["id"].(string), 10, 64)
		us.DB.Where(&us.Submission{ID: id}).First(&sub)
		if sub.ID == 0 {
			return nil, nil
		}
		return &sub, nil
	},
}

var submissionListQuery = &graphql.Field{
	Type:        listTypeOf(submissionType, "SubmissionList"),
	Description: "제출 정보 목록을 조회합니다.",
	Args:        graphql.FieldConfigArgument{},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var subs []us.Submission
		us.DB.Model(&us.Submission{}).Order("id desc").Find(&subs)

		var results []interface{}
		for _, s := range subs {
			sub := s
			results = append(results, &sub)
		}

		return listType{results}, nil
	},
}

var submitSourceMutation = &graphql.Field{
	Type:        graphql.NewNonNull(submissionType),
	Description: "소스코드를 제출합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name: "SubmitSourceInput",
				Fields: graphql.InputObjectConfigFieldMap{
					"gameID":                 &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"runtime":                &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"source":                 &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"description":            &graphql.InputObjectFieldConfig{Type: graphql.String},
					"competitorSubmissionID": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		errFailed := errors.New("제출에 실패했습니다. 잠시 후 다시 시도해주세요.")

		input := p.Args["input"].(map[string]interface{})
		user, ok := p.Context.Value(ctxUser).(*us.User)
		if !ok {
			return nil, errors.New("인증 정보가 올바르지 않습니다.")
		}

		var compSub us.Submission
		compSubID, _ := strconv.ParseInt(input["competitorSubmissionID"].(string), 10, 64)
		us.DB.Where(&us.Submission{ID: compSubID}).First(&compSub)
		if compSub.ID == 0 {
			return nil, errors.New("매칭 상대가 올바르지 않습니다.")
		}

		// 새로운 제출 정보를 생성하고, S3에 소스코드를 업로드합니다.
		desc := ""
		if d, ok := input["description"].(string); ok {
			desc = d
		}
		sub := us.NewSubmission(user.ID, input["gameID"].(string), input["runtime"].(string), desc)
		err := s3client.UploadFromBytes(sub.GetS3Key(), strings.NewReader(input["source"].(string)))
		if err != nil {
			return nil, errFailed
		}
		us.DB.Save(&sub)

		// 매칭을 생성합니다.
		match := us.NewMatch(input["gameID"].(string), user.ID, sub.ID, compSub.UserID, compSub.ID)
		us.DB.Save(&match)

		payload := us.SubmissionPayload{MatchUUID: match.UUID}
		err = submissionTask.Produce(string(payload.ToJSON()))
		if err != nil {
			return nil, errFailed
		}

		return sub, nil
	},
}

func init() {
	submissionType.AddFieldConfig("matchList", matchListQuery)
}
