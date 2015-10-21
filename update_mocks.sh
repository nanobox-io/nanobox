#! /bin/bash -e

mkdir -p \
  util/vagrant/mock \
  util/server/mock \
  util/server/mist/mock \
  util/s3/mock \
  util/print/mock \
  util/notify/mock \
  util/file/mock \
  util/file/hosts/mock \
  util/mock \
  config/mock

mockgen github.com/nanobox-io/nanobox/util/vagrant Vagrant > util/vagrant/mock/mock.go
mockgen github.com/nanobox-io/nanobox/util/server Server > util/server/mock/mock.go
mockgen github.com/nanobox-io/nanobox/util/server/mist Mist > util/server/mist/mock/mock.go
mockgen github.com/nanobox-io/nanobox/util/s3 S3 > util/s3/mock/mock.go
mockgen github.com/nanobox-io/nanobox/util/print Print > util/print/mock/mock.go
mockgen github.com/nanobox-io/nanobox/util/notify Notify > util/notify/mock/mock.go
mockgen github.com/nanobox-io/nanobox/util/file File > util/file/mock/mock.go
mockgen github.com/nanobox-io/nanobox/util/file/hosts Host > util/file/hosts/mock/mock.go
mockgen github.com/nanobox-io/nanobox/util Util > util/mock/mock.go
mockgen github.com/nanobox-io/nanobox/config Config > config/mock/mock.go