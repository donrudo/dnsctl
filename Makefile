DIR_CONF=${HOME}/.config/dnsctl
DIR_PLUGIN=${HOME}/.config/dnsctl/plugins
DIR_PKG_PLUGIN = ./pkg/plugins
DIR_SRC_PLUGIN = ./plugins

SRC_PLUGIN := $(wildcard $(DIR_SRC_PLUGIN)/*.go)
OBJ_PLUGIN := $(SRC_PLUGIN:$(DIR_SRC_PLUGIN)/%.go=%)
DIR_PKG=build
DIR_BUILD = ./build
DIR_BUILD_PLUGIN = ./build/plugins

GOPATH := $(shell go env GOPATH)
APP=dnsctl
MAIN=cmd/dnsctl/main.go

.PHONY: build-all run plugin
all: clean build-all plugin

clean:
	rm -rf $(DIR_PKG)
	rm -rf $(DIR_BUILD)

build-all: build-linux
build-linux:
	mkdir -p $(DIR_BUILD)/linux
	GOOS=linux go build -o $(DIR_BUILD)/linux/$(APP) $(MAIN)

plugin: build-plugin
	rm -rf  $(DIR_BUILD_PLUGIN)
	mv $(DIR_PKG_PLUGIN) $(DIR_BUILD_PLUGIN)

build-plugin: $(OBJ_PLUGIN)
	echo compiling

%: $(DIR_SRC_PLUGIN)/%.go
	go build -buildmode=plugin -o $(DIR_PKG_PLUGIN)/$@.so $<

install:
	mkdir -p $(DIR_PLUGIN)
	cp $(DIR_PKG)/plugins/* $(DIR_PLUGIN)/
	go install $(MAIN)

uninstall:
	rm -rf $(DIR_PLUGIN)
	rm ${GOPATH}/bin/dnsctl

