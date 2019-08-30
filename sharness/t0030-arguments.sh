#!/bin/sh

test_description="test the ipget argument parser"

. ./lib/test-lib.sh

test_expect_success "retrieve a known popular single file with a (HTTP) gateway URL" "
    ipget http://ipfs.io/ipfs/QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif &&
    echo 'c5ea0d6cacf1e54635685803ec4edbe0d4fe8465' > expected &&
    shasum cat.gif | cut -d ' ' -f 1 > actual &&
    diff expected actual
"

test_expect_success "retrieve a known popular single file with browser protocol URI" "
    ipget ipfs://QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif &&
    echo 'c5ea0d6cacf1e54635685803ec4edbe0d4fe8465' > expected &&
    shasum cat.gif | cut -d ' ' -f 1 > actual &&
    diff expected actual
"

test_expect_failure "don't allow non-(HTTP)gateway URLS" "
    ipget ftp://ipfs.io/ipfs/QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif
"

test_done
