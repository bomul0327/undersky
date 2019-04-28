package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"

	us "github.com/hellodhlyn/undersky"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"uuid":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"email":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"username":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"createdAt": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
		"updatedAt": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

type googleIdentity struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Google 서버에 토큰 검증을 요청합니다.
// Source : https://developers.google.com/identity/sign-in/web/backend-auth
func getGoogleIdentity(token string) (*googleIdentity, error) {
	req, _ := http.NewRequest("GET", "https://oauth2.googleapis.com/tokeninfo", nil)
	q := req.URL.Query()
	q.Add("id_token", token)
	req.URL.RawQuery = q.Encode()

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, errors.New("unauthorized")
	}

	defer res.Body.Close()
	var identity googleIdentity
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &identity)

	return &identity, nil
}

// Sign In
type signinOutput struct {
	Registered  bool
	AccessToken *string
	SecretToken *string
	ValidUntil  *time.Time
}

var signInOutputType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SignInOutput",
	Fields: graphql.Fields{
		"registered":  &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
		"accessToken": &graphql.Field{Type: graphql.String},
		"secretToken": &graphql.Field{Type: graphql.String},
		"validUntil":  &graphql.Field{Type: graphql.DateTime},
	},
})

var signInWithGoogleQuery = &graphql.Field{
	Type:        graphql.NewNonNull(signInOutputType),
	Description: "Google 계정으로 로그인합니다.",
	Args: graphql.FieldConfigArgument{
		"token": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		identity, err := getGoogleIdentity(p.Args["token"].(string))
		if err != nil {
			return nil, err
		}

		// 이미 존재하는 계정인지 여부를 체크한 후, 응답을 반환합니다.
		var user us.User
		us.DB.Where(&us.User{Email: identity.Email}).First(&user)
		if user.ID == 0 {
			return &signinOutput{Registered: false}, nil
		}

		cred := user.NewCredential()
		us.DB.Save(cred)

		return &signinOutput{
			Registered:  true,
			AccessToken: &cred.AccessToken,
			SecretToken: &cred.SecretToken,
			ValidUntil:  &cred.ValidUntil,
		}, nil
	},
}

// Sign Up
var registerUserWithGoogleMutation = &graphql.Field{
	Type:        graphql.NewNonNull(simpleResponseType),
	Description: "회원 계정을 등록합니다.",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name: "RegisterUserWithGoogleInput",
				Fields: graphql.InputObjectConfigFieldMap{
					"token":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					"username": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
				},
			}),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		input := p.Args["input"].(map[string]interface{})
		identity, err := getGoogleIdentity(input["token"].(string))
		if err != nil {
			return nil, errors.New("Google 인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
		}

		// 닉네임을 체크합니다.
		username := input["username"].(string)
		if len(username) < 2 || 20 < len(username) {
			return nil, errors.New("닉네임 형식이 일치하지 않습니다.")
		}

		var count int
		us.DB.Model(&us.User{}).Where("username ~* ?", username).Count(&count)
		if count != 0 {
			return nil, errors.New("중복된 닉네임입니다.")
		}

		// 새로운 회원 정보를 생성합니다.
		u, _ := uuid.NewRandom()
		us.DB.Save(&us.User{
			UUID:     u.String(),
			Email:    identity.Email,
			Username: username,
		})

		return &simpleResponse{"registered"}, nil
	},
}

var meQuery = &graphql.Field{
	Type:        graphql.NewNonNull(userType),
	Description: "현재 로그인한 유저 정보를 반환합니다.",
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		user, ok := p.Context.Value(ctxUser).(*us.User)
		if !ok {
			return nil, errors.New("인증 정보가 올바르지 않습니다.")
		}
		return user, nil
	},
}

type authClaims struct {
	AccessToken string `json:"accessToken"`
	jwt.StandardClaims
}

func getUserFromJWTToken(tokenString string) (*us.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authClaims{}, func(t *jwt.Token) (interface{}, error) {
		var cred us.Credential
		errs := us.DB.Where(&us.Credential{AccessToken: t.Claims.(*authClaims).AccessToken}).First(&cred).GetErrors()
		if len(errs) > 0 || cred.ID == 0 || cred.HasExpired() {
			return nil, errors.New("인증에 실패했습니다.")
		}

		return []byte(cred.SecretToken), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*authClaims); ok && token.Valid {
		var cred us.Credential
		us.DB.Where(&us.Credential{AccessToken: claims.AccessToken}).First(&cred)

		var user us.User
		us.DB.Where(&us.User{ID: cred.UserID}).First(&user)
		if user.ID == 0 {
			return nil, errors.New("인증에 실패했습니다.")
		}

		return &user, nil
	}
	return nil, errors.New("인증에 실패했습니다.")
}
