# Run tests
#
# Copyright (c) 2014 Christian Couder
# MIT Licensed; see the LICENSE file in this repository.
#

# NOTE: run with TEST_VERBOSE=1 for verbose sharness tests.

T = $(sort $(wildcard t[0-9][0-9][0-9][0-9]-*.sh))
BINS = bin/ipget bin/ipfs bin/go-sleep bin/pollEndpoint
SHARNESS = lib/sharness/sharness.sh

find_go_files = $(shell find $(1) -name "*.go")

# User might want to override those on the command line
GOFLAGS =

all: aggregate

clean: clean-test-results
	@echo "*** $@ ***"
	-rm -rf $(BINS)

clean-test-results:
	@echo "*** $@ ***"
	-rm -rf test-results

$(T): clean-test-results deps
	@echo "*** $@ ***"
	./$@

aggregate: clean-test-results $(T)
	@echo "*** $@ ***"
	lib/test-aggregate-results.sh

deps: $(SHARNESS) $(BINS)

$(SHARNESS): FORCE
	@echo "*** checking $@ ***"
	lib/install-sharness.sh

bin/ipget:
	mkdir -p bin
	$(MAKE) -C .. build
	cp ../ipget bin/ipget

# Usee the linked go-ipfs.
bin/ipfs:
	mkdir -p bin
	go build -o bin/ipfs github.com/ipfs/go-ipfs/cmd/ipfs

# Use go-sleep from our test deps.
bin/go-sleep:
	mkdir -p bin
	cd dependencies && go build -o ../bin/go-sleep github.com/chriscool/go-sleep

# Use pollEndpoints from our test deps.
bin/pollEndpoint:
	mkdir -p bin
	cd dependencies && go build -o ../bin/pollEndpoint github.com/whyrusleeping/pollEndpoint

race:
	make GOFLAGS=-race all

.PHONY: all clean $(T) aggregate FORCE
