DIR_CONF=${HOME}/.config/dnsctl
DIR_PLUGIN=${HOME}/.config/dnsctl/plugins
DIR_PKG_PLUGIN = ./build/plugins
DIR_SRC_PLUGIN = ./plugins
DIR_PKG = ./pkg
DIR_API = ./api

LDFLAGS=-ldflags "-X main.Version=`date  +%Y%d%m%H`"
SRC_PKG	   := $(wildcard $(DIR_PKG)/*_test.go)
SRC_API	   := $(wildcard $(DIR_API)/*.go)
SRC_PLUGIN := $(wildcard $(DIR_SRC_PLUGIN)/*.go)
OBJ_PLUGIN := $(SRC_PLUGIN:$(DIR_SRC_PLUGIN)/%.go=%)
DIR_BUILD = ./build
DIR_BUILD_PLUGIN = ./build/plugins

GOPATH := $(shell go env GOPATH)
APP=dnsctl
MAIN=cmd/dnsctl/main.go

.PHONY: build-all run plugin
all: clean build-all plugin

clean:
	rm -rf $(DIR_BUILD)

build-all: build-linux
build-linux:
	GOOS=linux go build $(LDFLAGS) -o $(DIR_BUILD)/linux/$(APP) $(MAIN)

plugin: $(OBJ_PLUGIN)

%: $(DIR_SRC_PLUGIN)/%.go
	go build -buildmode=plugin $(LDFLAGS) -o $(DIR_BUILD_PLUGIN)/$@.so $<

run: all
	$(DIR_BUILD)/linux/$(APP)

test:
	cd $(DIR_PKG);	go test
#test: $(SRC_PKG:$(DIR_PKG)/%.go)
#	go test $(SRC_PKG)

install:
	mkdir -p $(DIR_PLUGIN)
	cp $(DIR_PKG)/plugins/* $(DIR_PLUGIN)/
	go install $(MAIN)

uninstall:
	rm -rf $(DIR_PLUGIN)
	rm ${GOPATH}/bin/dnsctl

