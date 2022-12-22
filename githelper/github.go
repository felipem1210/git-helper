/*
Copyright © 2022 Felipe Macias felipem1210@gmail.com

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package githelper

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

type MyRepos []*github.Repository

// Authenticate with Github
func githubInitClient() (*github.Client, context.Context) {
	gh_token := os.Getenv("GIT_ACCESS_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh_token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client, ctx
}

func (myRepos MyRepos) GithubCreateRepos(f string) MyRepos {
	myRepos = myRepos.fromJsontoSliceOfStructs(f)
	var myReposComplete MyRepos
	client, ctx := githubInitClient()
	for _, repo := range myRepos {
		repo_options := &github.Repository{
			AllowForking:        repo.AllowForking,
			Name:                repo.Name,
			DefaultBranch:       repo.DefaultBranch,
			MasterBranch:        repo.MasterBranch,
			Organization:        repo.Organization,
			AllowRebaseMerge:    repo.AllowRebaseMerge,
			AllowSquashMerge:    repo.AllowSquashMerge,
			AllowMergeCommit:    repo.AllowMergeCommit,
			AllowAutoMerge:      repo.AllowAutoMerge,
			AllowUpdateBranch:   repo.AllowUpdateBranch,
			DeleteBranchOnMerge: repo.DeleteBranchOnMerge,
			Private:             repo.Private,
			AutoInit:            repo.AutoInit,
		}
		color.Green("Creating repo %s in organization %s\n", repo.GetName(), *repo.GetOrganization().Name)
		repo_info, _, err := client.Repositories.Create(ctx, *repo.GetOrganization().Name, repo_options)
		if err != nil {
			CheckIfError(err)
		}
		myReposComplete = append(myReposComplete, repo_info)
	}
	return myReposComplete
}

// Assigns a github team to a repository
func (myRepos MyRepos) GithubAssignTeamToRepo(f string, team string) {
	myRepos = myRepos.fromJsontoSliceOfStructs(f)
	client, ctx := githubInitClient()
	for _, repo := range myRepos {
		org := *repo.GetOrganization().Name
		color.Green("Assigning team %s to repo %s\n", team, repo.GetName())
		_, err := client.Teams.AddTeamRepoBySlug(ctx, org, team, org, repo.GetName(), &github.TeamAddTeamRepoOptions{
			Permission: "admin",
		})
		if err != nil {
			CheckIfError(err)
		}
	}
}

func (myRepos MyRepos) GetGithubRepositoriesInfo(org string) MyRepos {
	var myReposComplete MyRepos
	client, ctx := githubInitClient()
	// list all repositories for the authenticated user
	listOptions := &github.ListOptions{
		Page:    0,
		PerPage: 1000,
	}
	repos, _, err := client.Repositories.ListByOrg(ctx, org, &github.RepositoryListByOrgOptions{
		Type:        "all",
		Sort:        "full_name",
		Direction:   "asc",
		ListOptions: *listOptions,
	})
	if err != nil {
		CheckIfError(err)
	}
	myReposComplete = repos
	return myReposComplete
}

func (myRepos MyRepos) GithubGetOrg(f string) string {
	var org string
	myRepos = myRepos.fromJsontoSliceOfStructs(f)
	for _, repo := range myRepos {
		org = repo.GetOrganization().GetLogin()
	}
	return org
}

func (myRepos MyRepos) GithubGetRepoNames(f string, regex string) []string {
	myRepos = myRepos.fromJsontoSliceOfStructs(f)
	var repoNames []string
	for _, repo := range myRepos {
		if regex == "" {
			repoNames = append(repoNames, repo.GetName())
		} else {
			matched, _ := regexp.MatchString(regex, repo.GetName())
			if matched {
				repoNames = append(repoNames, repo.GetName())
			}
		}
	}
	return repoNames
}

func (myRepos MyRepos) GithubGetGitUrls(f string, regex string, auth string) []string {
	myRepos = myRepos.fromJsontoSliceOfStructs(f)
	var repoUrls []string
	var url string
	for _, repo := range myRepos {
		if auth == "ssh" {
			url = repo.GetSSHURL()
		} else if auth == "https" {
			url = repo.GetCloneURL()
		}
		if regex == "" {
			repoUrls = append(repoUrls, url)
		} else {
			matched, _ := regexp.MatchString(regex, url)
			if matched {
				repoUrls = append(repoUrls, url)
			}

		}
	}
	return repoUrls
}

type prCreateInfo github.NewPullRequest
type MyPrs []prInfo
type prInfo struct {
	Name     string `json:"name"`
	Title    string `json:"title,omitempty"`
	PrNumber int    `json:"pr_number,omitempty"`
	Body     string `json:"body,omitempty"`
	State    string `json:"state,omitempty"`
	Base     string `json:"base,omitempty"`
	Head     string `json:"head,omitempty"`
	Url      string `json:"url,omitempty"`
}

func (myPrs MyPrs) GithubCreatePr(org string, repos []string, f string, reviewers []string) MyPrs {
	prCreateInfoPointer := &prCreateInfo{}
	data := prCreateInfoPointer.fromJsontoStruct(f)
	client, ctx := githubInitClient()
	pr_options := &github.NewPullRequest{
		Title:               data.Title,
		Head:                data.Head,
		Base:                data.Base,
		Body:                data.Body,
		MaintainerCanModify: data.MaintainerCanModify,
		Draft:               github.Bool(false),
	}
	for _, repo := range repos {
		color.Green("Creating PR for repo %s\n", repo)
		pr_info, _, err := client.PullRequests.Create(ctx, org, repo, pr_options)
		if err != nil {
			CheckIfError(err)
		}
		fmt.Printf("PR created for repo: %s\n Url: %s\n", repo, pr_info.GetHTMLURL())
		my_pr_info := githubWritePrInfo(repo, pr_info)
		myPrs = append(myPrs, my_pr_info)
		if len(reviewers) != 0 {
			reviewers := &github.ReviewersRequest{
				Reviewers: reviewers,
			}
			_, _, err = client.PullRequests.RequestReviewers(ctx, org, repo, *pr_info.Number, *reviewers)
			if err != nil {
				CheckIfError(err)
			}
		}
	}
	return myPrs
}

func (myPrs MyPrs) GithubEditPr(org string, repos []string, f string) {
	myPrs = myPrs.fromJsontoSliceOfStructs(f)
	client, ctx := githubInitClient()
	for i, pr := range myPrs {
		color.Green("Modifying PR: %v", pr.Url)
		pr_branch := &github.PullRequestBranch{
			Ref: &pr.Base,
		}
		pr_update_options := &github.PullRequest{
			Number: &pr.PrNumber,
			Title:  &pr.Title,
			Body:   &pr.Body,
			State:  &pr.State,
			Base:   pr_branch,
		}
		_, _, err := client.PullRequests.Edit(ctx, org, repos[i], pr.PrNumber, pr_update_options)
		if err != nil {
			CheckIfError(err)
		}
	}
}

func (myPrs MyPrs) GithubMergePr(org string, repos []string, f string) {
	myPrs = myPrs.fromJsontoSliceOfStructs(f)
	client, ctx := githubInitClient()
	for i, pr := range myPrs {
		color.Green("Merging PR: %v", pr.Url)
		_, _, err := client.PullRequests.Merge(ctx, org, repos[i], pr.PrNumber, "merged!", &github.PullRequestOptions{
			DontDefaultIfBlank: false,
		})
		if err != nil {
			CheckIfError(err)
		}
	}
}

func githubWritePrInfo(repo string, pr_info *github.PullRequest) prInfo {
	p := prInfo{}
	//myPrsJson := myReposJson{}
	p.Name = repo
	p.Title = pr_info.GetTitle()
	p.Body = pr_info.GetBody()
	p.Base = pr_info.Base.GetRef()
	p.State = pr_info.GetState()
	p.Head = pr_info.Head.GetRef()
	p.PrNumber = pr_info.GetNumber()
	p.Url = pr_info.GetURL()
	return p
}
