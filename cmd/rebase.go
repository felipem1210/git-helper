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
	"fmt"

	"github.com/felipem1210/git-helper/githelper"
	"github.com/spf13/cobra"
)

// rebaseCmd represents the rebase command
var rebaseCmd = &cobra.Command{
	Use:   "rebase",
	Short: "Make a rebase from a base branch to the current branch in repository",
	Long: `Make a rebase from a base branch to the current branch in repository. To begin the rebase pass option --base-branch.
If you have conflicts you can pass the argument with "continue" or "abort" options.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rebase called")
		base_branch, _ := cmd.Flags().GetString("base-branch")
		json_file, _ := cmd.Flags().GetString("repo-info-json-file")
		repoNames := githelper.GetFromJsonReturnArray(json_file, "Name")
		githelper.GitRebase(repoNames, base_branch)
	},
}

func init() {
	rootCmd.AddCommand(rebaseCmd)
	rebaseCmd.PersistentFlags().StringP("base-branch", "b", "", "The name of the branch to rebase with. If you have conflicts use this argument with continue or abort.")
	rebaseCmd.MarkPersistentFlagRequired("org")
}
