#!/bin/sh

test_description="test the ipget command"

. ./lib/sharness/sharness.sh

test_expect_success "retrieve a known popular single file" "
    ipget QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif &&
    echo 'c5ea0d6cacf1e54635685803ec4edbe0d4fe8465' > expected &&
    sha1sum cat.gif | cut -d ' ' -f 1 > actual &&
    diff expected actual
"

test_expect_success "retrieve a known popular file with -o" "
    ipget -o meow.gif QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif &&
    echo 'c5ea0d6cacf1e54635685803ec4edbe0d4fe8465' > expected &&
    sha1sum meow.gif | cut -d ' ' -f 1 > actual &&
    diff expected actual
"

# TODO(noffle): recursive retrieval (directory)

# TODO(noffle): start a daemon, add + fetch a unique hash from it

test_done
