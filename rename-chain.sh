#!/bin/bash

# Required tools checking
command -v sed > /dev/null 2>&1 || { echo 'require sed'; exit 1; }
command -v npm > /dev/null 2>&1 || { echo 'require npm'; exit 1; }
command -v go > /dev/null 2>&1 || { echo 'require go 1.20+'; exit 1; }

# Working directory checking
if [ ! -d "./rename_chain" ] || [ ! -f "./go.mod" ]; then
    echo "Please run this script from the root of the repository."
    exit 1
fi

# Work on new branch to be able to compare changes
git branch -D branch-rename-chain
git checkout -b branch-rename-chain

set -eu # Abort execution if any command fails

# Run go tool to rename chain
go run --tags renamechain rename_chain/main.go

# Update dependencies
cd ./contracts
npm i
cd ..
go mod tidy

# Cleanup
rm -rf ./rename_chain
rm -rf ./build
rm -f ./rename-chain.sh

echo "Done."