package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"

	us "github.com/hellodhlyn/undersky"
)

var hmacSecret []byte

type authClaims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

func generateAccessToken(userUUID string) (string, error) {
	claim := authClaims{
		userUUID,
		jwt.StandardClaims{
			Issuer:    "UnderskyAPI",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().AddDate(0, 0, 7).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	if len(hmacSecret) == 0 {
		hmacSecret = []byte(os.Getenv("UNDERSKY_SECRET_KEY"))
	}
	return token.SignedString(hmacSecret)
}

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
}

var signInOutputType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SignInOutput",
	Fields: graphql.Fields{
		"registered":  &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
		"accessToken": &graphql.Field{Type: graphql.String},
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

		accessToken, _ := generateAccessToken(user.UUID)
		return &signinOutput{Registered: true, AccessToken: &accessToken}, nil
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
			return nil, err
		}

		u, _ := uuid.NewRandom()
		us.DB.Save(&us.User{
			UUID:  u.String(),
			Email: identity.Email,
		})

		return &simpleResponse{"registered"}, nil
	},
}
