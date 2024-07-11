# GITFIX
### *A command line interface that makes git problems a bit easier to handle*

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

EX. Lets say you have dev, uat, and main branches. main has become out of sync with dev, and now any branch thats based on main
will have merge conflicts or pull in extraneous commits when pointed at dev. 

We can solve this by 
* hopping over to dev
* creating a new feature branch based on dev (feature-dev)
* cherry picking files from the source feature branch and pulling them into our new feature-dev

This cli automates all that for us

## Usage

gitfix takes 2 required flag arguments;
* -t (target): the target branch (dev/uat/stage/main)
* -s (source): the source feature branch 

the -h argument displays a help text and running "gitfix" with no flags returns a welcome message


EXAMPLE:
    "gitfix -t dev -s DE-2564"

* check that both dev and DE-2564 exist
* hop over to dev and pull the latest
* run "git diff DE-2564 --name-only" to get a list of files
* present all files to the user who can;
    * press 'y' to accept all files with no revisions
    * press 'p' to loop over files 1 by 1 and choose which ones to keep (not recommended for large diffs)
    * press 'v' to enter vi and choose which files to keep
        * prepend "keep " to any file that should be kept
        * very similar to doing an interactive rebase in git and choosing "pick " or "reword " or "squash " etc. to modify files
        * type ":wq" to exit vi and save your choices
    * press 'r' to reset the file list to the original diff
    * press 'q' to quit
    * entering anything else will respawn the prompt
* Once the final file list is set
    * Create the DE-2564-dev branch based on dev
        * If already exists, it will simply check out the existing branch and not overwrite it
    * Checkout each file and apply it to DE-2564-dev
* Prints a final message to the user to double check the files that have been modified