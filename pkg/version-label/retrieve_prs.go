package versionlabel

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/barweiss/go-tuple"
	"github.com/google/go-github/v47/github"
	"go.uber.org/fx"

	"github.com/kaankoken/helper/pkg/helper"
	"github.com/kaankoken/versioning-tool/pkg"
)

// ResultStruct -> Dependency Injection Data Model for xxx Module
type ResultStruct struct {
	PrNumber  int
	LabelType tuple.T2[string, int]
}

var (
	// Major -> data type for major label
	Major tuple.T2[string, int] = tuple.New2("major", 0)

	// Minor -> data type for minor label
	Minor tuple.T2[string, int] = tuple.New2("minor", 1)

	// Patch -> data type for patch label
	Patch tuple.T2[string, int] = tuple.New2("patch", 2)

	// Urgent -> data type for urgent label
	Urgent tuple.T2[string, int] = tuple.New2("urgent", 0)
)

// VersionLabelModule -> Dependency Injection for VersionLabelModule module
var VersionLabelModule = fx.Options(
	fx.Provide(LisClosedPrs),
)

/*
AsyncFilterMergedPR -> Async concurrent handler to filter merged PRs

[logger] -> takes logger as an argument to log crash
[input] -> takes input as an argument that contains {repository name}, {repository owner} & {personal key}
[data] -> current PR object that in-progress
[wg] & [ch] -> concurrent job handlers
*/
func (client PrClient) AsyncFilterMergedPR(logger *helper.LogHandler, input *pkg.InputStruct, data *ResultStruct, wg *sync.WaitGroup, ch chan<- ResultStruct) {
	ctx := context.Background()
	result, _, err := client.Client.PullRequests.IsMerged(ctx, input.Owner, input.Repo, data.PrNumber)
	logger.Error(err)

	if result {
		ch <- ResultStruct{PrNumber: data.PrNumber, LabelType: data.LabelType}
	}
	defer wg.Done()
}

/*
FilterMergedPRs -> Main looper for list of filtered closed PRs

[logger] -> takes logger as an argument to log crash
[input] -> takes input as an argument that contains {repository name}, {repository owner} & {personal key}
[res] -> filtered PRs which are {"closed"}
*/
func FilterMergedPRs(client *PrClient, logger *helper.LogHandler, input *pkg.InputStruct, res *[]ResultStruct) {
	var wg sync.WaitGroup
	ch := make(chan ResultStruct)

	for _, data := range *res {
		wg.Add(1)
		go client.AsyncFilterMergedPR(logger, input, &data, &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := range ch {
		fmt.Println(i.LabelType.V1)
	}
}

/*
LisClosedPrs -> Getting list of PRs that closed according {Base} branch

[client] -> Function is directly connected to {*github.Client}

[logger] -> takes logger as an argument to log crash
[input] -> takes input as an argument that contains {repository name}, {repository owner} & {personal key}

[return] -> returns either successful filtered {array of ResultStruct} or {error}
*/
func LisClosedPrs(client *PrClient, logger *helper.LogHandler, input *pkg.InputStruct) (*[]ResultStruct, error) {
	filteredPRs := []ResultStruct{}

	// Gets only closed branches according to given base
	ctx := context.Background()
	options := github.PullRequestListOptions{Base: input.Base, State: "closed", ListOptions: github.ListOptions{PerPage: 100}}

	result, _, err := client.Client.PullRequests.List(ctx, input.Owner, input.Repo, &options)
	logger.Error(err)

	// Filtering PRs that contains versioning labels
	for _, v := range result {
		for _, l := range v.Labels {
			// TODO: rethink urgent trigger for {githubAction}
			//if strings.Contains(strings.ToLower(l.GetName()), strings.ToLower(Urgent.V1)) {
			//	filteredPRs = append(filteredPRs, ResultStruct{PrNumber: v.GetNumber(), LabelType: Urgent})
			//	break
			//}

			if strings.Contains(strings.ToLower(l.GetName()), strings.ToLower(Major.V1)) {
				filteredPRs = append(filteredPRs, ResultStruct{PrNumber: v.GetNumber(), LabelType: Major})
				break
			}

			if strings.Contains(strings.ToLower(l.GetName()), strings.ToLower(Minor.V1)) {
				filteredPRs = append(filteredPRs, ResultStruct{PrNumber: v.GetNumber(), LabelType: Minor})
				break
			}

			if strings.Contains(strings.ToLower(l.GetName()), strings.ToLower(Patch.V1)) {
				filteredPRs = append(filteredPRs, ResultStruct{PrNumber: v.GetNumber(), LabelType: Patch})
				break
			}
		}
	}

	return &filteredPRs, nil
}
