#! /bin/bash -e

mkdir -p util/vagrant/mock util/mock config/mock
mockgen github.com/nanobox-io/nanobox-cli/util/vagrant Vagrant > util/vagrant/mock/mock.go
mockgen github.com/nanobox-io/nanobox-cli/util Util > util/mock/mock.go
mockgen github.com/nanobox-io/nanobox-cli/config Config > config/mock/mock.go