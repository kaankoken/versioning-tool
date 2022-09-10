package versionlabel

import (
	"context"

	"github.com/google/go-github/v47/github"
	"github.com/kaankoken/versioning-tool/pkg"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
)

// PrClient -> Dependency Injection Data Model while wrapping {*github.Client} for xxx Module
type PrClient struct{ Client *github.Client }

// GithubClient -> DI for github client initialization
var GithubClient = fx.Options(fx.Provide(CreateGithubClient))

/*
CreateGithubClient -> Creates Github client using {input}

[input] -> Function is directly connected to {InputStruct}

[return] -> returns {*github.Client} wrapped with {PrClient}
*/
func CreateGithubClient(input *pkg.InputStruct) (client *PrClient) {
	ctx := context.Background()

	oauth := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: input.EncodedKey},
	)

	authClient := oauth2.NewClient(ctx, oauth)

	return &PrClient{Client: github.NewClient(authClient)}
}
