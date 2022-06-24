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
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func GitAdd(r []string) {
	for _, repo := range r {
		color.Green("Adding all content for repository %s", repo)
		getLocalRepo(repo)
		out := exec.Command("/bin/sh", "-c", "git add -A")
		f, err := pty.Start(out)
		if err != nil {
			CheckIfError(err)
		}
		io.Copy(os.Stdout, f)
	}
}

func GitCheckout(r []string, b string) {
	for _, repo := range r {
		color.Green("Repo: %s, checkout to branch %s", repo, b)
		getLocalRepo(repo)
		cmd := fmt.Sprintf("git checkout %s", b)
		out := exec.Command("/bin/sh", "-c", cmd)
		f, err := pty.Start(out)
		if err != nil {
			CheckIfError(err)
		}
		io.Copy(os.Stdout, f)
	}
}

func GitClone(r []string, urls []string) {
	for i := 0; i < len(r); i++ {
		color.Green("Cloning repository %s, url: %s", r[i], urls[i])
		_, err := git.PlainClone(getEnvValue("WORKING_DIR")+"/"+r[i], false, &git.CloneOptions{
			URL:      urls[i],
			Progress: os.Stdout,
			Auth:     getAuthOptions(),
		})
		if strings.Contains(fmt.Sprint(err), "already up-to-date") {
			fmt.Println(err)
		} else if err != nil {
			CheckIfError(err)
		}
		if err != nil {
			CheckIfError(err)
		}
	}
}

func GitCommit(r []string, message string) {
	for _, repo := range r {
		color.Green("Commiting the changes in %s repo", repo)
		r, _ := getLocalRepo(repo)
		w, _ := getGitWorktree(r)
		w.Commit(message, &git.CommitOptions{
			All: true,
		})
	}
}

func GitCreateBranch(r []string, b string) {
	for _, repo := range r {
		color.Green("Creating branch %s in %s repo", b, repo)
		r, _ := getLocalRepo(repo)
		w, _ := getGitWorktree(r)
		branch := fmt.Sprintf("refs/heads/%s", b)
		bRef := plumbing.ReferenceName(branch)
		ref, _ := r.Reference(bRef, true)
		if ref == nil {
			err := w.Checkout(&git.CheckoutOptions{
				Create: true,
				Force:  false,
				Branch: bRef,
			})
			if err != nil {
				CheckIfError(err)
			}
		}
	}
}

func GitFetch(r []string) {
	for _, repo := range r {
		color.Green("Fetching repo %s", repo)
		r, _ := getLocalRepo(repo)
		// First try to checkout branch
		err := r.Fetch(&git.FetchOptions{
			Auth: getAuthOptions(),
		})
		if strings.Contains(fmt.Sprint(err), "already up-to-date") {
			fmt.Println(err)
		} else if err != nil {
			CheckIfError(err)
		}
	}
}

func GitPull(r []string) {
	for _, repo := range r {
		r, _ := getLocalRepo(repo)
		w, _ := getGitWorktree(r)
		b, _ := getCurrentBranch(r)
		color.Green("Pulling changes for branch %s in %s repo", b, repo)
		err := w.Pull(&git.PullOptions{
			RemoteName:   "origin",
			SingleBranch: true,
			Auth:         getAuthOptions(),
			Force:        true,
		})
		if strings.Contains(fmt.Sprint(err), "already up-to-date") {
			fmt.Println(err)
		} else if err != nil {
			CheckIfError(err)
		}
	}
}

func GitPush(r []string) {
	for _, repo := range r {
		r, _ := getLocalRepo(repo)
		b, _ := getCurrentBranch(r)
		color.Green("Pushing changes for branch %s in %s repo", b, repo)
		err := r.Push(&git.PushOptions{
			RemoteName: "origin",
			Auth:       getAuthOptions(),
			Force:      true,
			Progress:   os.Stdout,
		})
		if strings.Contains(fmt.Sprint(err), "already up-to-date") {
			fmt.Println(err)
		} else if err != nil {
			CheckIfError(err)
		}
	}
}

func GitRebase(r []string, bb string) {
	for _, repo := range r {
		getLocalRepo(repo)
		r, _ := getLocalRepo(repo)
		b, _ := getCurrentBranch(r)
		color.Green("Rebasing branch %s on branch %s in repo %s", bb, b, repo)
		cmd := fmt.Sprintf("git rebase %s", bb)
		out := exec.Command("/bin/sh", "-c", cmd)
		f, err := pty.Start(out)
		if err != nil {
			CheckIfError(err)
		}
		io.Copy(os.Stdout, f)
	}
}

func GitReset(r []string) {
	for _, repo := range r {
		r, _ := getLocalRepo(repo)
		w, _ := getGitWorktree(r)
		color.Green("Reset changes to HEAD in %s repo", repo)
		err := w.Reset(&git.ResetOptions{
			Mode: git.HardReset,
		})
		if err != nil {
			CheckIfError(err)
		}
	}
}

func GitStatus(r []string) {
	for _, repo := range r {
		getLocalRepo(repo)
		color.Green("Status of repo %s", repo)
		cmd := "git status"
		out := exec.Command("/bin/sh", "-c", cmd)
		f, err := pty.Start(out)
		if err != nil {
			CheckIfError(err)
		}
		io.Copy(os.Stdout, f)
	}
}

// ###################################
// ## Helpers functions to get data ##
// ###################################

func getLocalRepo(name string) (*git.Repository, error) {
	wd := getEnvValue("WORKING_DIR") + "/" + name
	os.Chdir(wd)
	r, err := git.PlainOpen(wd)
	return r, err
}

func getGitWorktree(r *git.Repository) (*git.Worktree, error) {
	w, err := r.Worktree()
	CheckIfError(err)
	return w, err
}

func getCurrentBranch(r *git.Repository) (string, error) {
	h, err := r.Head()
	branch := h.Name().Short()
	CheckIfError(err)
	return branch, err
}

func getAuthOptions() *http.BasicAuth {
	return &http.BasicAuth{
		Username: getEnvValue("GIT_ACCESS_USER"),
		Password: getEnvValue("GIT_ACCESS_TOKEN"),
	}
}
