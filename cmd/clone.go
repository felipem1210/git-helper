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

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone multiple repositories in a WORKING_DIR",
	Long:  `Clone repos inside a Github Organization or Gitlab Group, providing the envars WORKING_DIR, GIT_ACCESS_TOKEN and GIT_ACCESS_USER`,
	Run: func(cmd *cobra.Command, args []string) {
		var repoNames []string
		var repoUrls []string
		provider, _ := cmd.Flags().GetString("provider")
		json_file, _ := cmd.Flags().GetString("repo-info-json-file")
		org, _ := cmd.Flags().GetString("org")
		regex, _ := cmd.Flags().GetString("regexp")
		target, _ := cmd.Flags().GetString("target")
		auth, _ := cmd.Flags().GetString("auth")
		if provider == "github" && json_file != "" {
			myRepos := githelper.MyRepos{}
			myRepos = myRepos.GetGithubRepositoriesInfo(org)
			err := githelper.WriteReposToJson(myRepos, json_file)
			githelper.CheckIfError(err)
			color.Green("The json file %s with repo info was written sucessfully", json_file)
			repoNames = myRepos.GithubGetRepoNames(json_file, regex)
			repoUrls = myRepos.GithubGetGitUrls(json_file, regex, auth)
		}
		githelper.GitClone(target, auth, repoNames, repoUrls)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.PersistentFlags().StringP("provider", "p", "github", "A provider to choose, options: gitub, gitlab")
	//cloneCmd.MarkPersistentFlagRequired("provider")
	cloneCmd.PersistentFlags().StringP("org", "o", "", "The github Organization to work with")
	cloneCmd.PersistentFlags().StringP("regexp", "r", "", "Regexp to apply to repositories clone process, to clone based on it")
	cloneCmd.PersistentFlags().StringP("auth", "a", "ssh", "Select if you want to clone using HTTPS or SSH.")
	cloneCmd.PersistentFlags().String("repo-info-json-file", "repos_info.json", "The name of the json file with info of the repos of the Github Org. It is read for each git local actions.")
}
