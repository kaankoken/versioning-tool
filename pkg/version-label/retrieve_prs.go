package versionlabel

import (
	"context"
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
var VersionLabelModule = fx.Options(fx.Provide(LisClosedPrs))

func AddToMap(fMap map[int][]ResultStruct, prNumber int, label tuple.T2[string, int]) {
	if val, ok := fMap[prNumber]; ok {
		val := append(val, ResultStruct{PrNumber: prNumber, LabelType: Urgent})
		fMap[prNumber] = val
	} else {
		arr := []ResultStruct{{PrNumber: prNumber, LabelType: Urgent}}
		fMap[prNumber] = arr
	}
}

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
func FilterMergedPRs(client *PrClient, logger *helper.LogHandler, input *pkg.InputStruct, res *map[int][]ResultStruct) {
	var wg sync.WaitGroup
	ch := make(chan ResultStruct)

	for _, val := range *res {
		wg.Add(1)

		if len(val) == 1 {
			go client.AsyncFilterMergedPR(logger, input, &val[0], &wg, ch)
		} else {
			if val[0].LabelType == Urgent {
				go client.AsyncFilterMergedPR(logger, input, &val[1], &wg, ch)
			} else {
				go client.AsyncFilterMergedPR(logger, input, &val[0], &wg, ch)
			}
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var versionTag tuple.T2[string, int]
	for i := range ch {
		if versionTag.Len() == 0 {
			versionTag = i.LabelType
		} else {
			if i.LabelType.V2 < versionTag.V2 {
				versionTag = i.LabelType
			}
		}
	}

	// TODO: generate the new tag
	// CreateNewTag()
}

/*
LisClosedPrs -> Getting list of PRs that closed according {Base} branch

[client] -> Function is directly connected to {*github.Client}

[logger] -> takes logger as an argument to log crash
[input] -> takes input as an argument that contains {repository name}, {repository owner} & {personal key}

[return] -> returns either successful filtered {map of ResultStruct} or {error}
*/
func LisClosedPrs(client *PrClient, logger *helper.LogHandler, input *pkg.InputStruct) (*map[int][]ResultStruct, error) {
	filteredPRs := map[int][]ResultStruct{}

	// Gets only closed branches according to given base
	ctx := context.Background()
	options := github.PullRequestListOptions{Base: input.Base, State: "closed", ListOptions: github.ListOptions{PerPage: 100}}

	result, _, err := client.Client.PullRequests.List(ctx, input.Owner, input.Repo, &options)
	logger.Error(err)

	urgentTaggedPrNumber := -100
	// Filtering PRs that contains versioning labels
	for _, v := range result {
		for _, l := range v.Labels {
			// TODO: rethink urgent trigger for {githubAction}
			if strings.Contains(strings.ToLower(l.GetName()), strings.ToLower(Urgent.V1)) {
				urgentTaggedPrNumber = v.GetNumber()
				AddToMap(filteredPRs, urgentTaggedPrNumber, Urgent)
				continue
			}

			if strings.Contains(strings.ToLower(l.GetName()), strings.ToLower(Major.V1)) {
				AddToMap(filteredPRs, v.GetNumber(), Major)
				break
			}

			if strings.Contains(strings.ToLower(l.GetName()), strings.ToLower(Minor.V1)) {
				AddToMap(filteredPRs, v.GetNumber(), Minor)
				break
			}

			if strings.Contains(strings.ToLower(l.GetName()), strings.ToLower(Patch.V1)) {
				AddToMap(filteredPRs, v.GetNumber(), Patch)

				break
			}
		}
	}

	if urgentTaggedPrNumber == -100 {
		return &filteredPRs, nil
	}

	return &map[int][]ResultStruct{urgentTaggedPrNumber: filteredPRs[urgentTaggedPrNumber]}, nil
}
