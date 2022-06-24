# GitHelper
A golang cli tool for automating GIT work in multiple repositories

Are you tired of typing the same git commands on multiple repositories of your project?
Are you tired of opening and updating Pull Requests for each repository?

This cli may be your solution to help you save time.

## Installation

Change the version for the [version](https://github.com/felipem1210/git-helper/tags) you want (withouth initial v)

### Linux amd64

```sh
export GITHELPER_VERSION=0.1.0
curl -L "https://github.com/felipem1210/git-helper/releases/download/v${GITHELPER_VERSION}/git-helper_${GITHELPER_VERSION}_linux_amd64.tar.gz" |tar xzv -C /tmp
sudo mv /tmp/git-helper /usr/local/bin/git-helper
```

### Envars needed

* Define the envar WORKING_DIR. This must be the folder where you have all the repositories you want to manage.

```sh
export WORKING_DIR=$(pwd)
```

Define `GIT_ACCESS_USER` and `GIT_ACCESS_TOKEN` with your Github Username and Token

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

You can run following actions in github:

* Create repositories. 

You need to provide a json file with info of the repos you want to create. Check [this json file](examples/json-files/new_repos.json) with the example of the file. Place the file in your $WORKING_DIR

```sh
cp examples/json-files/new_repos.json $WORKING_DIR
```

After creating the repositories a file `repos_info.json` will be created. don't erase it

* Create/update/merge pull requests

You need to provide a json file with info of the PRs you want to create. Check [this json file](examples/json-files/new_pr.json) with the example of the file. Place the file in your $WORKING_DIR

```sh
cp examples/json-files/new_prs.json $WORKING_DIR
```

After creating the repositories a file `pr_info.json` will be created. don't erase it

If you want to update you can only update these fields: `title, body, state, base, maintainer_can_modify`.

## Gitlab integration

COMING SOON

## Usage

Use the `--help`

```sh
git-helper --help
```