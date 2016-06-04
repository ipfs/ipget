# Minimum version numbers for software required to build IPFS
IPFS_MIN_GO_VERSION = 1.5.2
IPFS_MIN_GX_VERSION = 0.6
IPFS_MIN_GX_GO_VERSION = 1.1

ifeq ($(TEST_NO_FUSE),1)
  go_test=go test -tags nofuse
else
  go_test=go test
endif


dist_root=/ipfs/QmXZQzBAFuoELw3NtjQZHkWSdA332PyQUj6pQjuhEukvg8
gx_bin=bin/gx-v0.7.0
gx-go_bin=bin/gx-go-v1.2.0

# use things in our bin before any other system binaries
export PATH := bin:$(PATH)
export IPFS_API ?= v04x.ipfs.io

all: help

go_check:
	@bin/check_go_version $(IPFS_MIN_GO_VERSION)

bin/gx-v%:
	@echo "installing gx $(@:bin/gx-%=%)"
	@bin/dist_get ${dist_root} gx $@ $(@:bin/gx-%=%)
	rm -f bin/gx
	ln -s $(@:bin/%=%) bin/gx

bin/gx-go-v%:
	@echo "installing gx-go $(@:bin/gx-go-%=%)"
	@bin/dist_get ${dist_root} gx-go $@ $(@:bin/gx-go-%=%)
	rm -f bin/gx-go
	ln -s $(@:bin/%=%) bin/gx-go

gx_check: ${gx_bin} ${gx-go_bin}

path_check:
	@bin/check_go_path $(realpath $(shell pwd)) $(realpath $(GOPATH)/src/github.com/ipfs/ipget)

deps: go_check gx_check path_check
	${gx_bin} --verbose install --global

install: deps
	go install

build: deps
	go build

clean:
	rm -rf ./ipget

uninstall:
	go clean github.com/ipfs/ipget

PHONY += all help gx_check
PHONY += go_check deps install build nofuse clean uninstall

##############################################################
# A semi-helpful help message

help:
	@echo 'DEPENDENCY TARGETS:'
	@echo ''
	@echo '  gx_check        - Installs or upgrades gx and gx-go'
	@echo '  deps            - Download dependencies using gx'
	@echo ''
	@echo 'BUILD TARGETS:'
	@echo ''
	@echo '  all          - print this help message'
	@echo '  build        - Build binary'
	@echo '  install      - Build binary and install into $$GOPATH/bin'
	@echo ''
	@echo 'CLEANING TARGETS:'
	@echo ''
	@echo '  clean        - Remove binary from build directory'
	@echo '  uninstall    - Remove binary from $$GOPATH/bin'
	@echo ''
	@echo 'TESTING TARGETS:'
	@echo ''
	@echo '  COMING SOON(TM)'
	@echo ''

PHONY += help

.PHONY: $(PHONY)
