#!/bin/sh

# Git init, add and commit project
echo "> Initializing git-remote"

if [ -d ".git" ]; then
    echo "Warning: Existing git remote was found."
    exit
fi

git init . >/dev/null
added_remote=0
remote_name=""

# Ask if user wants to add a git-remote
while true; do
    echo -n "> Add git-remote? [y/N] "
    read yn

    # Evaluate input
    if [ ! -n "$yn" ] || [ "$yn" = "N" ] || [ "$yn" = "n" ]; then
        break
    elif [ "$yn" = "Y" ] || [ "$yn" = "y" ]; then
        # Read and validate remote-name
        echo -n "> Remote-name: "
        read remote_name

        if [ -z "$remote_name" -o "$remote_name" = " " ]; then
            echo "> Remote-name can't be empty"
            return 4
        fi

        # Read and validate remote-url
        echo -n "> Remote-URL: "
        read remote_url

        if [ -z "$remote_url" -o "$remote_url" = " " ]; then
            echo "> Remote-URL can't be empty"
            return 5
        fi

        # Append .git if missing
        if [[ ! $remote_url =~ ^.*\.git$ ]]; then
            remote_url="${remote_url}.git"
        fi

        # Add remote
        git remote add "$remote_name" "$remote_url" >/dev/null
        added_remote=1
        break
    else
        echo "> Invalid input"
        continue
    fi
done

# Git add and commit
git add . >/dev/null
git commit -m "Create project" >/dev/null
git tag -a v0.1.0 -m "Create project" >/dev/null
git checkout --quiet -b develop >/dev/null

# Push if remote was added
if [ $added_remote -eq 1 ]; then
    git push --quiet "$remote_name" master >/dev/null
    git push --quiet -u "$remote_name" develop >/dev/null
fi
