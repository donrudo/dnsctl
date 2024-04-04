DIR_CONF=${HOME}/.config/dnsctl
DIR_PLUGIN=${HOME}/.config/dnsctl/plugins
DIR_PKG_PLUGIN = ./build/plugins
DIR_SRC_PLUGIN = ./plugins

LDFLAGS=-ldflags "-X main.Version=`date  +%Y%d%m%H`"
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

plugin: build-plugin

build-plugin: $(OBJ_PLUGIN)

%: $(DIR_SRC_PLUGIN)/%.go
	go build -buildmode=plugin $(LDFLAGS) -o $(DIR_BUILD_PLUGIN)/$@.so $<

run: all
	$(DIR_BUILD)/linux/$(APP)

install:
	mkdir -p $(DIR_PLUGIN)
	cp $(DIR_PKG)/plugins/* $(DIR_PLUGIN)/
	go install $(MAIN)

uninstall:
	rm -rf $(DIR_PLUGIN)
	rm ${GOPATH}/bin/dnsctl

