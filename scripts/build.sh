#!/usr/bin/env bash
set -e

# disable cgo for true static binaries (will work on alpine linux)
export CGO_ENABLED=0;

# vfor versioning
getCurrCommit() {
  echo `git rev-parse --short HEAD | tr -d "[ \r\n\']"`
}

getCurrTag() {
  echo `git describe --always --tags --abbrev=0 | tr -d "[v\r\n]"`
}

BUILD_DATE=`date -u +%y%m%dT%H%M`
# ^for versioning

# remove any previous builds that may have failed
[ -e "./.build" ] && \
  echo "Cleaning up old builds..." && \
  rm -rf "./.build"

printf "\nBuilding nanobox...\n"

# build nanobox
gox -ldflags "-s -X github.com/nanobox-io/nanobox/util/odin.apiKey=$API_KEY \
			  -X github.com/nanobox-io/nanobox/models.nanoVersion=$(getCurrTag) \
			  -X github.com/nanobox-io/nanobox/models.nanoCommit=$(getCurrCommit) \
			  -X github.com/nanobox-io/nanobox/models.nanoBuild=$BUILD_DATE" \
			  -osarch "darwin/amd64 linux/amd64 windows/amd64" -output="./.build/v2/{{.OS}}/{{.Arch}}/nanobox"

printf "\nWriting version file...\n"
echo -en "Nanobox Version $(getCurrTag)-$BUILD_DATE ($(getCurrCommit))" > ./.build/v2/version

printf "\nBuilding nanobox updater...\n"

# change into updater directory and build nanobox updater
cd ./updater && gox -osarch "darwin/amd64 linux/amd64 windows/amd64" -ldflags="-s" -output="../.build/v2/{{.OS}}/{{.Arch}}/nanobox-update"

#cd ..

#printf "\nCompacting binaries...\n"
#upx ./.build/v2/*/amd64/nanobox*
