#!/usr/bin/env bash
set -e

# for versioning
getCurrCommit() {
  echo `git rev-parse --short HEAD | tr -d "[ \r\n\']"`
}

# for versioning
getCurrTag() {
  echo `git describe --always --tags --abbrev=0 | tr -d "[v\r\n]"`
}

# remove any previous builds that may have failed
[ -e "./.build" ] && \
  echo "Cleaning up old builds..." && \
  rm -rf "./.build"

printf "\nBuilding nanobox...\n"

# build nanobox
gox -ldflags "-X main.bugsnagToken=$BUGSNAG_TOKEN
	-X github.com/nanobox-io/nanobox/commands.tag=$(getCurrTag)
	-X github.com/nanobox-io/nanobox/commands.commit=$(getCurrCommit)
	-X github.com/nanobox-io/nanobox/util/mixpanel.token=$MIXPANEL_TOKEN" \
	-osarch "darwin/amd64 linux/amd64 windows/amd64" \
	-output="./.build/v1/{{.OS}}/{{.Arch}}/nanobox"

printf "\nBuilding nanobox updater...\n"

# change into updater directory and build nanobox updater
cd ./updater && gox -osarch "darwin/amd64 linux/amd64 windows/amd64" -output="../.build/v1/{{.OS}}/{{.Arch}}/nanobox-update"
