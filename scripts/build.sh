#!/usr/bin/env bash
set -e

# try and use the correct MD5 lib (depending on user OS darwin/linux)
MD5=$(which md5 || echo "$(which md5sum) | cut -f 1" )

# remove any previous builds that may have failed
[ -e "./build" ] && \
  echo "Cleaning up old builds..." && \
  rm -rf "./build"

# build nanobox
echo "Building nanobox"
gox -osarch "darwin/amd64 linux/amd64 windows/amd64" -output="./build/v1/{{.OS}}/{{.Arch}}/nanobox"
