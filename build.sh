#!/usr/bin/env bash
set -e

# try and use the correct MD5 lib (depending on user OS darwin/linux)
MD5=$(which md5 || echo "$(which md5sum) | cut -f 1" )

# remove any previous builds that may have failed
[ -e "./build" ] && \
  echo "Cleaning up old builds..." && \
  rm -rf "./build"

# build nanobox
echo "Building NANOBOX and uploading it to 's3://tools.nanobox.io/nanobox/v1'"
gox -osarch "darwin/amd64 linux/amd64 windows/amd64" -output="./build/v1/{{.OS}}/{{.Arch}}/nanobox"

# look through each os/arch/file and generate an md5 for each
echo "Generating md5s..."
for os in $(ls ./build/v1); do
  for arch in $(ls ./build/v1/${os}); do
    for file in $(ls ./build/v1/${os}/${arch}); do
      cat "./build/v1/${os}/${arch}/${file}" | ${MD5} >> "./build/v1/${os}/${arch}/${file}.md5"
    done
  done
done

# upload to AWS S3
echo "Uploading builds to S3..."
aws s3 sync ./build/v1/ s3://tools.nanobox.io/nanobox/v1 --grants read=uri=http://acs.amazonaws.com/groups/global/AllUsers --region us-east-1

#
echo "Cleaning up..."

# remove build
[ -e "./build" ] && \
  echo "Removing build files..." && \
  rm -rf "./build"
