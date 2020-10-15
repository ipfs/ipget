# Minimum version numbers for software required to build IPFS
IPFS_MIN_GO_VERSION = 1.14

# use things in our bin before any other system binaries
export PATH := bin:$(PATH)

go_check:
	@bin/check_go_version $(IPFS_MIN_GO_VERSION)

deps: go_check
	go mod download

install: deps
	go install

build: deps
	go build

clean:
	rm -rf ./ipget

uninstall:
	go clean github.com/ipfs/ipget

PHONY += help go_check deps install build clean

##############################################################
# A semi-helpful help message

help:
	@echo 'DEPENDENCY TARGETS:'
	@echo ''
	@echo '  deps            - Download dependencies'
	@echo ''
	@echo 'BUILD TARGETS:'
	@echo ''
	@echo '  help         - print this help message'
	@echo '  build        - Build binary'
	@echo '  install      - Build binary and install into $$GOPATH/bin'
	@echo ''
	@echo 'CLEANING TARGETS:'
	@echo ''
	@echo '  clean        - Remove binary from build directory'
	@echo ''
	@echo 'TESTING TARGETS:'
	@echo ''
	@echo '  COMING SOON(TM)'
	@echo ''

PHONY += help

.PHONY: $(PHONY)
