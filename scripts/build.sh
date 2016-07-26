#!/usr/bin/env bash
set -e

# remove any previous builds that may have failed
[ -e "./build" ] && \
  echo "Cleaning up old builds..." && \
  rm -rf "./build"

printf "\nBuilding nanobox...\n"

# build nanobox
gox -osarch "darwin/amd64 linux/amd64 windows/amd64" -output="./build/v1/{{.OS}}/{{.Arch}}/nanobox"

printf "\nBuilding nanobox updater...\n"

# change into updater directory and build nanobox updater
cd ./updater && gox -osarch "darwin/amd64 linux/amd64 windows/amd64" -output="../build/v1/{{.OS}}/{{.Arch}}/nanobox-update"
