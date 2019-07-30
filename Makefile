SHELL=/bin/bash

# APP info
APP_NAME := proxy
APP_VERSION := 0.2.0_SNAPSHOT
APP_CONFIG := $(APP_NAME).yml
APP_STATIC_FOLDER := static
APP_UI_FOLDER := ui
APP_PLUGIN_FOLDER := plugin

# Get release version from environment
ifneq "$(VERSION)" ""
   APP_VERSION := $(VERSION)
endif

# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  GOPATH := ~/go
  #$(error Please set the environment variable GOPATH before running `make`)
endif


PATH := $(PATH):$(GOPATH)/bin

# Go environment
CURDIR := $(shell pwd)
OLDGOPATH:= $(GOPATH)
NEWGOPATH:= $(CURDIR):$(CURDIR)/vendor:$(GOPATH)

GO        := GO15VENDOREXPERIMENT="1" go
GOBUILD  := GOPATH=$(NEWGOPATH) CGO_ENABLED=1  $(GO) build -ldflags -s -gcflags "-m -m"
GOBUILDNCGO  := GOPATH=$(NEWGOPATH) CGO_ENABLED=0  $(GO) build -ldflags -s
GOTEST   := GOPATH=$(NEWGOPATH) CGO_ENABLED=1  $(GO) test -ldflags -s

ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"
GO_FILES=$(find . -iname '*.go' | grep -v /vendor/)
PKGS=$(go list ./... | grep -v /vendor/)

FRAMEWORK_FOLDER := $(CURDIR)/../framework/
FRAMEWORK_BRANCH := master
FRAMEWORK_VENDOR_FOLDER := $(CURDIR)/vendor/
FRAMEWORK_VENDOR_BRANCH := master

FRAMEWORK_OFFLINE_BUILD := ""
ifneq "$(OFFLINE_BUILD)" ""
   FRAMEWORK_OFFLINE_BUILD := $(OFFLINE_BUILD)
endif

.PHONY: all build update test clean

default: build

build: config
	@#echo $(GOPATH)
	@echo $(NEWGOPATH)
	$(GOBUILD) -o bin/$(APP_NAME)
	@$(MAKE) restore-generated-file

build-cmd: config
	cd cmd/backup && GOOS=darwin GOARCH=amd64 $(GOBUILDNCGO) -o ../../bin/backup-darwin
	cd cmd/backup && GOOS=linux  GOARCH=amd64 $(GOBUILDNCGO) -o ../../bin/backup-linux64
	cd cmd/backup && GOOS=windows GOARCH=amd64 $(GOBUILDNCGO) -o ../../bin/backup-windows64.exe
	@$(MAKE) restore-generated-file

# used to build the binary for gdb debugging
build-race: clean config update-vfs
	$(GOBUILD) -gcflags "-m -N -l" -race -o bin/$(APP_NAME)
	@$(MAKE) restore-generated-file

tar: build
	cd bin && tar cfz ../bin/$(APP_NAME).tar.gz $(APP_NAME) $(APP_CONFIG)

cross-build: clean config update-vfs
	$(GO) test
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/$(APP_NAME)-windows64.exe
	GOOS=darwin  GOARCH=amd64 $(GOBUILD) -o bin/$(APP_NAME)-darwin64
	GOOS=linux  GOARCH=amd64 $(GOBUILD) -o bin/$(APP_NAME)-linux64
	@$(MAKE) restore-generated-file


build-win:
	CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64     $(GOBUILD) -o bin/$(APP_NAME)-windows64.exe
	CC=i686-w64-mingw32-gcc   CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386         $(GOBUILD) -o bin/$(APP_NAME)-windows32.exe

build-linux:
	GOOS=linux  GOARCH=amd64  $(GOBUILD) -o bin/$(APP_NAME)-linux64
	GOOS=linux  GOARCH=386    $(GOBUILD) -o bin/$(APP_NAME)-linux32
	GOOS=linux  GOARCH=arm   GOARM=5    $(GOBUILD) -o bin/$(APP_NAME)-armv5

build-darwin:
	GOOS=darwin  GOARCH=amd64     $(GOBUILD) -o bin/$(APP_NAME)-darwin64
	GOOS=darwin  GOARCH=386       $(GOBUILD) -o bin/$(APP_NAME)-darwin32

build-bsd:
	GOOS=freebsd  GOARCH=amd64    $(GOBUILD) -o bin/$(APP_NAME)-freebsd64
	GOOS=freebsd  GOARCH=386      $(GOBUILD) -o bin/$(APP_NAME)-freebsd32
	GOOS=netbsd  GOARCH=amd64     $(GOBUILD) -o bin/$(APP_NAME)-netbsd64
	GOOS=netbsd  GOARCH=386       $(GOBUILD) -o bin/$(APP_NAME)-netbsd32
	GOOS=openbsd  GOARCH=amd64    $(GOBUILD) -o bin/$(APP_NAME)-openbsd64
	GOOS=openbsd  GOARCH=386      $(GOBUILD) -o bin/$(APP_NAME)-openbsd32

all: clean config update-vfs cross-build restore-generated-file

all-platform: clean config update-vfs cross-build-all-platform restore-generated-file

cross-build-all-platform: clean config build-bsd build-linux build-darwin build-win  restore-generated-file

format:
	go fmt $$(go list ./... | grep -v /vendor/)

clean_data:
	rm -rif dist
	rm -rif data
	rm -rif log

clean: clean_data
	rm -rif bin
	mkdir bin

init:
	@echo building $(APP_NAME) $(APP_VERSION)
	@if [ ! -d $(FRAMEWORK_FOLDER) ]; then echo "framework does not exist";(cd ../&&git clone -b $(FRAMEWORK_BRANCH) https://github.com/infinitbyte/framework.git) fi
	@if [ ! -d $(FRAMEWORK_VENDOR_FOLDER) ]; then echo "framework vendor does not exist";(git clone  -b $(FRAMEWORK_VENDOR_BRANCH) https://github.com/infinitbyte/framework-vendor.git vendor) fi
	@if [ "" == $(FRAMEWORK_OFFLINE_BUILD) ]; then (cd $(FRAMEWORK_FOLDER) && git pull origin $(FRAMEWORK_BRANCH)); fi;
	@if [ "" == $(FRAMEWORK_OFFLINE_BUILD) ]; then (cd vendor && git pull origin $(FRAMEWORK_VENDOR_BRANCH)); fi;


update-generated-file:
	@echo "update generated info"
	@echo -e "package config\n\nconst LastCommitLog = \""`git log -1 --pretty=format:"%h, %ad, %an, %s"` "\"\nconst BuildDate = \"`date`\"" > config/generated.go
	@echo -e "\nconst Version  = \"$(APP_VERSION)\"" >> config/generated.go


restore-generated-file:
	@echo "restore generated info"
	@echo -e "package config\n\nconst LastCommitLog = \"N/A\"\nconst BuildDate = \"N/A\"" > config/generated.go
	@echo -e "\nconst Version = \"0.0.1-SNAPSHOT\"" >> config/generated.go


update-vfs:
	@if [ -d $(APP_STATIC_FOLDER) ]; then  echo "generate static files";$(GO) get github.com/infinitbyte/framework/cmd/vfs;(cd $(APP_STATIC_FOLDER) && vfs -ignore="static.go|.DS_Store" -o static.go -pkg static . ) fi

update-template-ui:
	@if [ -d $(APP_UI_FOLDER) ]; then  (echo "generate main UI pages";$(GO) get github.com/infinitbyte/ego/cmd/ego;cd $(APP_UI_FOLDER)/ && ego) fi
	@if [ -d $(APP_PLUGIN_FOLDER) ]; then  (echo "generate plugin UI pages";$(GO) get github.com/infinitbyte/ego/cmd/ego;cd $(APP_PLUGIN_FOLDER)/ && ego) fi

#config: init update-vfs update-template-ui
config: init update-vfs update-template-ui update-generated-file
	@echo "update configs"
	@# $(GO) env
	@mkdir -p bin
	@cp $(APP_CONFIG) bin/$(APP_CONFIG)


dist: cross-build package

dist-major-platform: all package

dist-all-platform: all-platform package-all-platform

package:
	@echo "Packaging"
	cd bin && tar cfz ../bin/darwin64.tar.gz darwin64  $(APP_CONFIG)
	cd bin && tar cfz ../bin/linux64.tar.gz linux64  $(APP_CONFIG)
	cd bin && tar cfz ../bin/windows64.tar.gz windows64  $(APP_CONFIG)

package-all-platform: package-darwin-platform package-linux-platform package-windows-platform
	@echo "Packaging all"
	cd bin && tar cfz ../bin/freebsd64.tar.gz     $(APP_NAME)-freebsd64  $(APP_CONFIG)
	cd bin && tar cfz ../bin/freebsd32.tar.gz     $(APP_NAME)-freebsd32  $(APP_CONFIG)
	cd bin && tar cfz ../bin/netbsd64.tar.gz      $(APP_NAME)-netbsd64  $(APP_CONFIG)
	cd bin && tar cfz ../bin/netbsd32.tar.gz      $(APP_NAME)-netbsd32  $(APP_CONFIG)
	cd bin && tar cfz ../bin/openbsd64.tar.gz     $(APP_NAME)-openbsd64  $(APP_CONFIG)
	cd bin && tar cfz ../bin/openbsd32.tar.gz     $(APP_NAME)-openbsd32  $(APP_CONFIG)


package-darwin-platform:
	@echo "Packaging Darwin"
	cd bin && tar cfz ../bin/darwin64.tar.gz      $(APP_NAME)-darwin64 $(APP_CONFIG)
	cd bin && tar cfz ../bin/darwin32.tar.gz      $(APP_NAME)-darwin32 $(APP_CONFIG)

package-linux-platform:
	@echo "Packaging Linux"
	cd bin && tar cfz ../bin/linux64.tar.gz     $(APP_NAME)-linux64 $(APP_CONFIG)
	cd bin && tar cfz ../bin/linux32.tar.gz     $(APP_NAME)-linux32 $(APP_CONFIG)
	cd bin && tar cfz ../bin/armv5.tar.gz       $(APP_NAME)-armv5   $(APP_CONFIG)

package-windows-platform:
	@echo "Packaging Windows"
	cd bin && tar cfz ../bin/windows64.tar.gz   $(APP_NAME)-windows64.exe $(APP_CONFIG)
	cd bin && tar cfz ../bin/windows32.tar.gz   $(APP_NAME)-windows32.exe $(APP_CONFIG)

test:
	go get -u github.com/kardianos/govendor
	go get github.com/stretchr/testify/assert
	govendor test +local
	#$(GO) test -timeout 60s ./... --ignore ./vendor
	#GORACE="halt_on_error=1" go test ./... -race -timeout 120s  --ignore ./vendor
	#go test -bench=. -benchmem

check:
	$(GO)  get github.com/golang/lint/golint
	$(GO)  get honnef.co/go/tools/cmd/megacheck
	test -z $(gofmt -s -l $GO_FILES)    # Fail if a .go file hasn't been formatted with gofmt
	$(GO) test -v -race $(PKGS)            # Run all the tests with the race detector enabled
	$(GO) vet $(PKGS)                      # go vet is the official Go static analyzer
	@echo "go tool vet"
	go tool vet main.go
	go tool vet core
	go tool vet modules
	megacheck $(PKGS)                      # "go vet on steroids" + linter
	golint -set_exit_status $(PKGS)    # one last linter

errcheck:
	go get github.com/kisielk/errcheck
	errcheck -blank $(PKGS)

cover:
	go get github.com/mattn/goveralls
	go test -v -cover -race -coverprofile=data/coverage.out
	goveralls -coverprofile=data/coverage.out -service=travis-ci -repotoken=$COVERALLS_TOKEN

cyclo:
	go get -u github.com/fzipp/gocyclo
	gocyclo -top 10 -over 12 $$(ls -d */ | grep -v vendor)
