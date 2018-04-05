SHELL=/bin/bash

# Default PROXY version
PROXY_VERSION := 0.1.0_SNAPSHOT

# Get release version from environment
ifneq "$(VERSION)" ""
   PROXY_VERSION := $(VERSION)
endif

# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

# Go environment
CURDIR := $(shell pwd)
OLDGOPATH:= $(GOPATH)
NEWGOPATH:= $(CURDIR):$(CURDIR)/vendor:$(GOPATH)

GO        := GO15VENDOREXPERIMENT="1" go
GOBUILD  := GOPATH=$(NEWGOPATH) CGO_ENABLED=1  $(GO) build -ldflags -s
GOTEST   := GOPATH=$(NEWGOPATH) CGO_ENABLED=1  $(GO) test -ldflags -s

ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"

.PHONY: all build update test clean

default: build

build: config
	@#echo $(GOPATH)
	@echo $(NEWGOPATH)
	$(GOBUILD) -o bin/proxy
	@$(MAKE) restore-generated-file

cross-build: clean config update-ui
	$(GO) test
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/proxy-windows64.exe
	GOOS=darwin  GOARCH=amd64 $(GOBUILD) -o bin/proxy-darwin64
	GOOS=linux  GOARCH=amd64 $(GOBUILD) -o bin/proxy-linux64
	@$(MAKE) restore-generated-file

update-generated-file:
	@echo "update generated info"
	@echo -e "package config\n\nconst LastCommitLog = \""`git log -1 --pretty=format:"%h, %ad, %an, %s"` "\"\nconst BuildDate = \"`date`\"" > config/generated.go
	@echo -e "\nconst Version  = \"$(PROXY_VERSION)\"" >> config/generated.go


restore-generated-file:
	@echo "restore generated info"
	@echo -e "package config\n\nconst LastCommitLog = \"N/A\"\nconst BuildDate = \"N/A\"" > config/generated.go
	@echo -e "\nconst Version = \"0.0.1-SNAPSHOT\"" >> config/generated.go


format:
	gofmt -l -s -w .

clean_data:
	rm -rif data
	rm -rif log

clean: clean_data
	rm -rif bin
	mkdir bin

init-version:
	@echo building PROXY $(PROXY_VERSION)


update-ui:
	@echo "generate static files"
	@$(GO) get github.com/infinitbyte/esc
	@(cd static&& esc -ignore="static.go|build_static.sh|.DS_Store" -o static.go -pkg static ../static )

update-template-ui:
	@echo "generate UI pages"
	@$(GO) get github.com/infinitbyte/ego/cmd/ego
	@cd ui/ && ego

config: init-version update-ui update-template-ui update-generated-file
	@echo "update configs"
	@# $(GO) env
	@mkdir -p bin
	@cp proxy.yml bin/proxy.yml
