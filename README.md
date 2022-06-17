# GitHelper
A golang cli tool for automating GIT work in multiple repositories

Are you tired of typing the same git commands on multiple repositories of your project?
Are you tired of opening and updating Pull Requests for each repository?

This cli may be your solution to help you save time.

## Installation



## Git integration

You can run local git commands alongside multiple repositories that are inside an specific folder:

* clone (Repos grouped in a Gitlab Group or Github Org).
* checkout
* add all new content
* create new branch locally
* pull
* commit
* reset to HEAD
* rebase
* push
* fetch

## Github integration



## Using the script.

* Make a simbolic lynk of the script to your local `$PATH`
```shell
    sudo ln -s $(pwd)/git-helper /usr/local/bin
```
* Define the envar WORKING_DIR. This must be the folder where you have all the repositories you want to manage.
```shell
    cd <your_working_dir>
    export WORKING_DIR=$(pwd)
```

## Options

* `--action` -- Specify git command to run

* Used in clone action:

  * `--private-token` -- Is Your [Gitlab/GitHub private token]
  * `--git-provider` -- `gitlab` or `github`
  * `--git-domain` -- The domain of the gitlab server (Gitlab)
  * `--group-id` -- The ID of the Group of repositories you want to clone (Gitlab)
  * `--org-name` -- The name of the GitHub organization

* `--message` -- Message for a commit 
* `--branch` -- The branch where you want to work 

--branch option is not mandatory, if not specified the action will be taken in the current branch of each repository

* `--base-branch` -- The base branch when you want to make a rebase
* `--extra-options` -- Any extra options for git command

### Examples for using differents actions:

* Git Clone repositories of group - **Gitlab**:
```shell
  git-helper --action=clone --git-provider=gitlab --private-token="YOUR_GITLAB_TOKEN" --git-domain="YOUR_GITLAB_DOMAIN" --group-id="YOUR_GROUP_ID"  
```

* Git Clone repositories of organization - **Github**:
```shell
  git-helper --action=clone --git-provider=github --private-token="YOUR_GITHUB_TOKEN"  --org-name="YOUR_ORG_NAME"
```

* Create a commit:
```shell
  git-helper --action=commit --message="<your_message>" 
```

* Create and checkout to new branch (locally):
```shell
  git-helper --action=create-branch --branch=<your_branch>
```

* Delete a branch (locally):
```shell
  git-helper --action=delete --branch=<your_branch>
```

* Checkout to existing branch (without pull):
```shell
    git-helper --action=checkout --branch=<your_branch>
```

* Pull a branch:
```shell
  git-helper --action=pull --branch=<your_branch>
```

* Reset all content to HEAD (equivalent to git reset --hard HEAD):
```shell
  git-helper --action=reset   
```

* Add all changes (equivalent to git add -A):
```shell
  git-helper --action=add
```

* Rebase a branch with an origin branch:
```shell
  git-helper --action=rebase --base-branch=<origin_branch> --branch=<your_branch>  
```

* Push all changes to upstream:
```shell
  git-helper --action=push 
```

* Fetch all branches from upstream:
```shell
  git-helper --action=fetch
```
