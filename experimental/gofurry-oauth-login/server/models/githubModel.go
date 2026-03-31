package models

import "github.com/GoFurry/gofurry-oauth-login/api/proto/githuboauth"

type GithubOAuthServer struct {
	githuboauth.UnimplementedGithubOAuthServiceServer
	clientID     string // GitHub Client ID
	clientSecret string // GitHub Client Secret
}
