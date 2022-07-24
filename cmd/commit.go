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

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a commit for each repository.",
	Long:  `Create a commit to each repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		target, _ := cmd.Flags().GetString("target")
		repoNames := githelper.ListDirectories(target)
		githelper.GitCommit(target, repoNames, message)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.PersistentFlags().StringP("commit", "m", "", "Message for the commit")
	commitCmd.MarkPersistentFlagRequired("commit")
}
