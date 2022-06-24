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
	"errors"
	"os"

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
		if provider == "github" && json_file != "" {
			myRepos := githelper.MyRepos{}
			repoNames = myRepos.GithubGetRepoNames(json_file)
			repoUrls = myRepos.GithubGetCloneUrls(json_file)
			if _, err := os.Stat(json_file); errors.Is(err, os.ErrNotExist) {
				myRepos = myRepos.GetGithubRepositoriesInfo(org)
				err := githelper.WriteReposToJson(myRepos, json_file)
				if err != nil {
					githelper.CheckIfError(err)
				} else {
					color.Green("The json file %s with repo info was written sucessfully", json_file)
				}
			}
		}
		githelper.GitClone(repoNames, repoUrls)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.PersistentFlags().StringP("provider", "p", "", "A provider to choose, options: gitub, gitlab")
	cloneCmd.MarkPersistentFlagRequired("provider")
	cloneCmd.PersistentFlags().StringP("org", "o", "", "The github Organization to work with")
}
