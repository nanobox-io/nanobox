#!/usr/bin/env bash
set -e

# upload to AWS S3
echo "Uploading builds to S3..."
aws s3 sync ./build/ s3://tools.nanopack.io/mist --grants read=uri=http://acs.amazonaws.com/groups/global/AllUsers --region us-east-1
