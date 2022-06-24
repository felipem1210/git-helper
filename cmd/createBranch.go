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
	"github.com/felipem1210/git-helper/githelper"
	"github.com/spf13/cobra"
)

// createBranchCmd represents the createBranch command
var createBranchCmd = &cobra.Command{
	Use:   "create-branch",
	Short: "Create a branch locally and checkout into that branch.",
	Long:  `Create a branch locally and checkout into that branch for each repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		branch, _ := cmd.Flags().GetString("branch")
		json_file, _ := cmd.Flags().GetString("repo-info-json-file")
		myRepos := githelper.MyRepos{}
		repoNames := myRepos.GithubGetRepoNames(json_file)
		githelper.GitCreateBranch(repoNames, branch)
	},
}

func init() {
	rootCmd.AddCommand(createBranchCmd)
	createBranchCmd.PersistentFlags().StringP("branch", "b", "", "The name of the branch to create locally")
	createBranchCmd.MarkPersistentFlagRequired("branch")
}
