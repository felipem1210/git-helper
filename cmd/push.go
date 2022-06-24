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

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push changes of a branch to remote",
	Long:  `Push changes of a branch from remote. It will make the Push of current branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		json_file, _ := cmd.Flags().GetString("repo-info-json-file")
		repoNames := githelper.GetFromJsonReturnArray(json_file, "Name")
		githelper.GitPush(repoNames)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
