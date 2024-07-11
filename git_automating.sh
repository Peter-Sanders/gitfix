#!/bin/zsh 

# This is a useful util script only in the world where our git implementation is suboptimal
# It takes in a feature branch and a target branch
# It will checkout the target branch, compare it to the feature, and list all modified files
# that file list will be presented to the user for acceptance
# if the list of files is acceptable, it will create a new branch based on the target
# EX: feature = "dummy" target = "target" will result in a new branch "dummy-target" created off of target
# Then, the script will checkout each of the modified files and stage them on the newly createed branch
# From here, its up to the user to review all modified files for acceptance and do the rest
#
#
# WORKFLOW:
#
# commit to your feature branch 
# run this script at the top level of whatever repo you're working out of 

# Checks if a branch exists locally and returns True/False accordingly
check_branch_existence() {
  branch=$1
  result=$(git branch --list "$branch" | wc -l)
  
  result="${result##*([[:space:]])}"
  result="${result%%*([[:space:]])}"
  if [ "$result" != 1 ]; then
    return 1
  else 
    return 0
  fi
}

# Makes a new feature branch and checks to make sure it exists first
# Will checkout an existing branch and not overwrite it
make_new_feature_branch() {
  env_feature_branch=$1

  if check_branch_existence $env_feature_branch; then
    echo "Branch '$env_feature_branch' already exists."
    git checkout "$env_feature_branch"

  else
    git checkout -b $env_feature_branch
    echo "Created and checked out branch '$env_feature_branch'."
fi
}

# Checks if the files are what we want them to be
check_files() {
    files=($@)

    print_files $files
    print "CHECK THIS LIST OF FILES CAREFULLY!!!!!!\nType 'y' to proceed, 'n' to quit, or 'p' to pick only the files you want included (y/n/p): " 
    read response

    case "$response" in

        [yY])
            make_new_feature_branch $env_feature_branch
            new_files=$files
            ;;
        [pP])
            new_files=pick_files $files
            ;;
        [nN])
            echo "Exiting Script"
            exit 0
            ;;
        *) 
            echo "Invalid Input. Exiting"
            exit 1
            ;;
    esac
    
    return $new_files
}

# Pick from a list of files, check them, and return them back
pick_files() {
    files=($@)

    echo "Pick only the files you want to include."
    new_files=()
    for i in $files; do 
      echo "Include $i ? (y/n)"
      read response 
      case "$repsonse" in 
        [yY])
          new_files+=($i)
          ;;
      esac
    done
    for file in "${new_files[@]}"
      do
        echo "$file"
    done

    final_files=check_files $new_files

    return $final_files
}

# Print out an array of files
print_files() {
    files=($@)

    for file in "${files[@]}"
    do
        echo "$file"
    done
}

# Set branch vars
branch_to_merge_into=$1
feature_branch_prime=$2
env_feature_branch=${feature_branch_prime}-${branch_to_merge_into}

# Make sure the source and destination branches exist
if ! check_branch_existence $branch_to_merge_into; then 
  echo "target branch $branch_to_merge_into doesn't exist, exiting."
  exit 1
fi

if ! check_branch_existence $feature_branch_prime; then 
  echo "source branch $branch_to_merge_into doesn't exist, exiting."
  exit 1
fi

# Checkout target and get it up to date
git checkout $branch_to_merge_into && git pull 

# Get the list of modified files between the source and destination
files=($(git diff $feature_branch_prime --name-only))

# print_files $files

# Get the list of files we want to pull into the new branch
good_files=check_files $files

# Checkout the files we want
for i in $good_files; do
  git checkout $feature_branch_prime $i
done

exit 0

