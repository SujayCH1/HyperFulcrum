#!/usr/bin/env bash

set -e

MAIN="main"

BRANCHES=(
    api
    cache
    conn
    fe
    nodem
    parser
    schema
)

echo "Fetching remotes..."
git fetch --all --prune

for branch in "${BRANCHES[@]}"; do
    echo "Updating $branch..."
    git checkout "$branch"
    git pull --rebase
done

echo "Updating $MAIN..."
git checkout "$MAIN"
git pull --rebase

for branch in "${BRANCHES[@]}"; do
    echo "Merging $branch..."
    git merge --no-ff "$branch"
done

echo "Pushing $MAIN..."
git push

for branch in "${BRANCHES[@]}"; do
    echo "Rebasing $branch..."
    git checkout "$branch"
    git rebase "$MAIN"

    echo "Pushing $branch..."
    git push --force-with-lease
done

git checkout "$MAIN"

echo "Done."
