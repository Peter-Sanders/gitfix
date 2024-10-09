# GITFIX

## *A command line interface that makes git problems a bit easier to handle*

## Installation

* clone the repo
* run "make install" (does require sudo)

  * this will
    * build the binary locally
    * make it executable (chmod and such)
    * move it to /usr/local/bin
    * move the manpage to /opt/homebrew/share/man/man8 (man gitfix)
* test installation by running "gitfix" in the terminal

## Description

This cli is designed to cherry pick files from a source branch and move it into a branch ready to merge into a target.

EX. Lets say you have dev, uat, and main branches. Somehow, code was pushed to main that never made it into dev or uat. So when a feature branch based on main is pointed at dev, extra commits and files are staged that are outside the scope of the original feature branch. This confuses developers at best and at worst creates merge conflicts.

We can solve this by taking the same steps:

* hopping over to dev
* creating a new feature branch based on dev (feature-dev)
* checking the diff between main and the source feature branch
* cherry picking files from the source feature branch and pulling them into our new feature-dev branch

This cli automates all that for us

## Usage

gitfix takes 2 required flag arguments;

* -t (target): the target branch (dev/uat/stage/main)
* -s (source): the source feature branch

And one optional flag argument:

* -d (default_branch): the default branch from which the source branch was branched off of (defaults to main)

The -h argument displays a help text and running "gitfix" with no flags returns a welcome message

EXAMPLE:
    "gitfix -t dev -s DE-2564"

Behind the scenes gitfix will:

* fetch origin
* check that both dev (target), DE-2564 (source), and main (default) exist (without the -d arg default is assumed to be "main")
* hop over to main and pull the latest
* run "git diff DE-2564 --name-status" to get a list of files affected and if they were added/modified/deleted
* present the list of files by action to the user
  * if y then move on
  * if n then exit (TODO: loop this back in with the cherry-picking and vi)
* hop to dev and pull the latest
* Create the DE-2564-dev branch based on dev
  * If already exists, it will simply check out the existing branch and not overwrite it
* Checkout each modified file and apply it to DE-2564-dev
* Delete each file that was deleted by the source on DE-2564-dev
* Prints a final message to the user to double check the files that have been affected before committing and pushing the branch
