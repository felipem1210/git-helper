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
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func GitAdd(t string, r []string) {
	for _, repo := range r {
		color.Green("Adding all content for repository %s", repo)
		getLocalRepo(t, repo)
		out := exec.Command("/bin/sh", "-c", "git add -A")
		f, err := pty.Start(out)
		if err != nil {
			CheckIfError(err)
		}
		io.Copy(os.Stdout, f)
	}
}

func GitCheckout(t string, r []string, b string) {
	for _, repo := range r {
		color.Green("Repo: %s, checkout to branch %s", repo, b)
		getLocalRepo(t, repo)
		cmd := fmt.Sprintf("git checkout %s", b)
		out := exec.Command("/bin/sh", "-c", cmd)
		f, err := pty.Start(out)
		if err != nil {
			CheckIfError(err)
		}
		io.Copy(os.Stdout, f)
	}
}

func GitClone(target string, auth string, r []string, urls []string) {
	createFolder(target)
	for i := 0; i < len(r); i++ {
		color.Green("Cloning repository %s, url: %s", r[i], urls[i])
		cloneDir := getEnvValue("WORKING_DIR") + "/" + target + "/" + r[i]
		_, err := git.PlainClone(cloneDir, false, getCloneOptions(auth, urls[i]))
		if strings.Contains(fmt.Sprint(err), "already up-to-date") {
			fmt.Println(err)
		} else if err != nil {
			CheckIfError(err)
		}
	}
}

func GitCommit(t string, r []string, message string) {
	for _, repo := range r {
		color.Green("Commiting the changes in %s repo", repo)
		r, _ := getLocalRepo(t, repo)
		w, _ := getGitWorktree(r)
		w.Commit(message, &git.CommitOptions{
			All: true,
		})
	}
}

func GitCreateBranch(t string, r []string, b string) {
	for _, repo := range r {
		color.Green("Creating branch %s in %s repo", b, repo)
		r, _ := getLocalRepo(t, repo)
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
			CheckIfError(err)
		}
	}
}

func GitFetch(t string, r []string) {
	for _, repo := range r {
		detected := detectIfSshOrHttps(t, repo)
		color.Green("Fetching repo %s", repo)
		r, _ := getLocalRepo(t, repo)
		// First try to checkout branch
		err := r.Fetch(getFetchOptions(detected))
		if strings.Contains(fmt.Sprint(err), "already up-to-date") {
			fmt.Println(err)
		} else if err != nil {
			CheckIfError(err)
		}
	}
}

func GitPull(t string, r []string) {
	for _, repo := range r {
		detected := detectIfSshOrHttps(t, repo)
		r, _ := getLocalRepo(t, repo)
		w, _ := getGitWorktree(r)
		b, _ := getCurrentBranch(r)
		color.Green("Pulling changes for branch %s in %s repo", b, repo)
		err := w.Pull(getPullOptions(detected))
		if strings.Contains(fmt.Sprint(err), "already up-to-date") {
			fmt.Println(err)
		} else if err != nil {
			CheckIfError(err)
		}
	}
}

func GitPush(t string, r []string) {
	for _, repo := range r {
		detected := detectIfSshOrHttps(t, repo)
		r, _ := getLocalRepo(t, repo)
		b, _ := getCurrentBranch(r)
		color.Green("Pushing changes for branch %s in %s repo", b, repo)
		err := r.Push(getPushOptions(detected))
		if strings.Contains(fmt.Sprint(err), "already up-to-date") {
			fmt.Println(err)
		} else if err != nil {
			CheckIfError(err)
		}
	}
}

func GitRebase(t string, r []string, bb string) {
	for _, repo := range r {
		r, _ := getLocalRepo(t, repo)
		b, _ := getCurrentBranch(r)
		color.Green("Rebasing branch %s on branch %s in repo %s", bb, b, repo)
		cmd := fmt.Sprintf("git rebase %s", bb)
		out := exec.Command("/bin/sh", "-c", cmd)
		f, err := pty.Start(out)
		CheckIfError(err)
		io.Copy(os.Stdout, f)
	}
}

func GitReset(t string, r []string) {
	for _, repo := range r {
		r, _ := getLocalRepo(t, repo)
		w, _ := getGitWorktree(r)
		color.Green("Reset changes to HEAD in %s repo", repo)
		err := w.Reset(&git.ResetOptions{
			Mode: git.HardReset,
		})
		CheckIfError(err)
	}
}

func GitStatus(t string, r []string) {
	for _, repo := range r {
		getLocalRepo(t, repo)
		color.Green("Status of repo %s", repo)
		cmd := "git status"
		out := exec.Command("/bin/sh", "-c", cmd)
		f, err := pty.Start(out)
		CheckIfError(err)
		io.Copy(os.Stdout, f)
	}
}

// ###################################
// ## Helpers functions to get data ##
// ###################################

func getLocalRepo(target string, name string) (*git.Repository, error) {
	wd := getEnvValue("WORKING_DIR") + "/" + target + "/" + name
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

func getHttpsAuth() *http.BasicAuth {
	return &http.BasicAuth{
		Username: getEnvValue("GIT_ACCESS_USER"),
		Password: getEnvValue("GIT_ACCESS_TOKEN"),
	}
}

func getSshAuth() *ssh.PublicKeys {
	sshPath := getEnvValue("HOME") + "/.ssh/id_rsa"
	sshKey, _ := ioutil.ReadFile(sshPath)
	publicKey, error := ssh.NewPublicKeys("git", []byte(sshKey), "")
	CheckIfError(error)
	return publicKey
}

func detectIfSshOrHttps(target string, name string) string {
	var detected string
	wd := getEnvValue("WORKING_DIR") + "/" + target + "/" + name
	os.Chdir(wd)
	cmd := "git remote show origin"
	out, err := exec.Command("/bin/sh", "-c", cmd).Output()
	CheckIfError(err)
	if strings.Contains(string(out), "git@") {
		detected = "ssh"
	} else if strings.Contains(string(out), "https://") {
		detected = "https"
	}
	return detected
}

func getCloneOptions(auth string, url string) *git.CloneOptions {
	var cloneOptions *git.CloneOptions
	if auth == "ssh" {
		cloneOptions = &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
			Auth:     getSshAuth(),
		}
	} else if auth == "https" {
		cloneOptions = &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
			Auth:     getHttpsAuth(),
		}
	}
	return cloneOptions
}

func getPullOptions(auth string) *git.PullOptions {
	var pullOptions *git.PullOptions
	if auth == "ssh" {
		pullOptions = &git.PullOptions{
			RemoteName:   "origin",
			SingleBranch: true,
			Auth:         getSshAuth(),
			Force:        true,
		}

	} else if auth == "https" {
		pullOptions = &git.PullOptions{
			RemoteName:   "origin",
			SingleBranch: true,
			Auth:         getHttpsAuth(),
			Force:        true,
		}
	}
	return pullOptions
}

func getPushOptions(auth string) *git.PushOptions {
	var pushOptions *git.PushOptions
	if auth == "ssh" {
		pushOptions = &git.PushOptions{
			RemoteName: "origin",
			Auth:       getSshAuth(),
			Force:      true,
			Progress:   os.Stdout,
		}

	} else if auth == "https" {
		pushOptions = &git.PushOptions{
			RemoteName: "origin",
			Auth:       getHttpsAuth(),
			Force:      true,
			Progress:   os.Stdout,
		}
	}
	return pushOptions
}

func getFetchOptions(auth string) *git.FetchOptions {
	var fetchOptions *git.FetchOptions
	if auth == "ssh" {
		fetchOptions = &git.FetchOptions{
			Auth: getSshAuth(),
		}

	} else if auth == "https" {
		fetchOptions = &git.FetchOptions{
			Auth: getHttpsAuth(),
		}
	}
	return fetchOptions
}
