TAG=`git describe --always --tags --abbrev=0 | tr -d "[v\r\n]"`
COMMIT=`git rev-parse --short HEAD | tr -d "[ \r\n\']"`
BUILD_DATE=`date -u +%y%m%dT%H%M`
GITSTATUS='$(shell git status 2> /dev/null | tail -n1)'
DIRTY="$(shell [ $(GITSTATUS) = 'no changes added to commit (use "git add" and/or "git commit -a")' ] && echo -n "*")"
GO_LDFLAGS="-s -X github.com/nanobox-io/nanobox/util/odin.apiKey=$(API_KEY) -X github.com/nanobox-io/nanobox/models.nanoVersion=$(TAG) -X github.com/nanobox-io/nanobox/models.nanoCommit=$(COMMIT)$(DIRTY) -X github.com/nanobox-io/nanobox/models.nanoBuild=$(BUILD_DATE)"

default: build

local: linux windows darwin

clean:
	@echo "Cleaning old builds"
	@rm -rf "./.build"

# go get github.com/mitchellh/gox
build: clean
	@echo "Building nanobox"
	@gox -ldflags=$(GO_LDFLAGS) -osarch "darwin/amd64 linux/amd64 windows/amd64" -output="./.build/v2/{{.OS}}/{{.Arch}}/nanobox"
	@echo -en "Nanobox Version $(TAG)-$(BUILD_DATE) ($(COMMIT))" > ./.build/v2/version
	@echo "Building nanobox-update"
	@cd ./updater && gox -osarch "darwin/amd64 linux/amd64 windows/amd64" -ldflags="-s" -output="../.build/v2/{{.OS}}/{{.Arch}}/nanobox-update"

linux:
	@echo "Building nanobox-linux"
	@GOOS=linux go build -ldflags=$(GO_LDFLAGS) -o nanobox-linux

windows:
	@echo "Building nanobox-windows"
	@GOOS=windows go build -ldflags=$(GO_LDFLAGS) -o nanobox-windows

darwin:
	@echo "Building nanobox-darwin"
	@GOOS=darwin go build -ldflags=$(GO_LDFLAGS) -o nanobox-darwin

# go get github.com/kardianos/govendor
test: 
	@govendor test +local -v


.PHONY: fmt test clean build linux windows darwin
