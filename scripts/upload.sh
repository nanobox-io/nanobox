#!/usr/bin/env bash
set -e

# try and use the correct MD5 lib (depending on user OS darwin/linux)
MD5=$(which md5 || echo "$(which md5sum) | cut -f 1" )

echo "Generating md5s..."

# look through each os/arch/file and generate an md5 for each
for os in $(ls ./.build/v1); do
  for arch in $(ls ./.build/v1/${os}); do
    for file in $(ls ./.build/v1/${os}/${arch}); do
      cat "./.build/v1/${os}/${arch}/${file}" | ${MD5} >> "./.build/v1/${os}/${arch}/${file}.md5"
    done
  done
done

# upload to AWS S3
echo "Uploading builds to S3..."
aws s3 sync ./.build/v1/ s3://tools.nanobox.io/nanobox/v1 --grants read=uri=http://acs.amazonaws.com/groups/global/AllUsers --region us-east-1
