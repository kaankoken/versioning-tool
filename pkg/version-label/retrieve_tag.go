package versionlabel

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/barweiss/go-tuple"
	"github.com/google/go-github/v47/github"
	"github.com/kaankoken/helper/pkg/helper"
	"github.com/kaankoken/versioning-tool/pkg"
)

/*
GenerateNewTag -> Gets existing version tag & increases according to label type

[tag] -> {existing tag} needed to be bumped
[label] -> current {label} type to {bump} the {existing tag}

[return] -> returns bumped version of {tag}
*/
func GenerateNewTag(tag string, label tuple.T2[string, int], logger *helper.LogHandler) string {
	tagPieces := strings.Split(strings.Replace(tag, "v", "", 1), ".")

	if label.V1 == Major.V1 {
		val, err := strconv.Atoi(tagPieces[0])
		logger.Error(err)

		tagPieces[0] = fmt.Sprintf("%d", val+1)
	}
	if label.V1 == Minor.V1 {
		val, err := strconv.Atoi(tagPieces[1])
		logger.Error(err)

		tagPieces[1] = fmt.Sprintf("%d", val+1)
	}
	if label.V1 == Patch.V1 {
		val, err := strconv.Atoi(tagPieces[2])
		logger.Error(err)

		tagPieces[2] = fmt.Sprintf("%d", val+1)

	}

	return "v" + strings.Join(tagPieces, ".")
}

/*
LatestVersionTag -> Gets latest created {tag}

[client] -> Function is directly connected to {*github.Client}

[logger] -> takes logger as an argument to log crash
[input] -> takes input as an argument that contains {repository name}, {repository owner} & {personal key}

[return] -> returns either existing version {vx.x.x} or {v0.0.0}
*/
func LatestVersionTag(client *PrClient, logger *helper.LogHandler, input *pkg.InputStruct) string {
	ctx := context.Background()
	options := github.ListOptions{PerPage: 100}

	// Get most recent tag created for stable version
	result, _, err := client.Client.Repositories.ListTags(ctx, input.Owner, input.Repo, &options)
	logger.Error(err)

	baseVersion := "v0.0.0"

	if len(result) > 0 {
		baseVersion = result[0].GetName()
	}

	return baseVersion
}
