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
package cmd

import (
	"github.com/fatih/color"
	"github.com/felipem1210/git-helper/githelper"
	"github.com/spf13/cobra"
)

// prCmd represents the createPr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests in Github.",
	Long:  `Create/update/merge pull requests in Github for all repos. You need to pass the organization and a JSON file with the info needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		repos_json_file, _ := cmd.Flags().GetString("repo-info-json-file")
		new_pr_json_file, _ := cmd.Flags().GetString("new-pr-json-file")
		pr_info_json_file, _ := cmd.Flags().GetString("pr-info-json-file")
		reviewers, _ := cmd.Flags().GetStringSlice("reviewers")
		create, _ := cmd.Flags().GetBool("create")
		update, _ := cmd.Flags().GetBool("update")
		merge, _ := cmd.Flags().GetBool("merge")
		target, _ := cmd.Flags().GetString("target")

		repos := githelper.MyRepos{}
		repoNames := githelper.ListDirectories(target)
		org := repos.GithubGetOrg(repos_json_file)

		myPrs := githelper.MyPrs{}
		if create {
			myPrs = myPrs.GithubCreatePr(org, repoNames, new_pr_json_file, reviewers)
			err := githelper.WritePrsToJson(myPrs, pr_info_json_file)
			githelper.CheckIfError(err)
			color.Green("The json file %s with pr info was written sucessfully", pr_info_json_file)
		} else if update {
			myPrs.GithubEditPr(org, repoNames, pr_info_json_file)
		} else if merge {
			myPrs.GithubMergePr(org, repoNames, pr_info_json_file)
		}
	},
}

func init() {
	var reviewers []string
	rootCmd.AddCommand(prCmd)
	prCmd.PersistentFlags().String("new-pr-json-file", "new_pr.json", "The json file needed to create the PR.")
	prCmd.PersistentFlags().String("pr-info-json-file", "pr_info.json", "The json file with info of the PRs created.")
	prCmd.PersistentFlags().BoolP("create", "c", false, "The json file with info of the PRs created.")
	prCmd.PersistentFlags().BoolP("update", "u", false, "The json file with info of the PRs created.")
	prCmd.PersistentFlags().BoolP("merge", "m", false, "The json file with info of the PRs created.")
	prCmd.Flags().StringSliceVarP(&reviewers, "reviewers", "r", []string{}, "List of usernames of reviewers for the pull request")
}
