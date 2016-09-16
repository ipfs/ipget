#!/bin/sh

test_description="test the ipget command by spawning a shell"

. ./lib/sharness/sharness.sh

test_expect_success "retrieve a known popular single file" "
    ipget --node=spawn QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif &&
    echo 'c5ea0d6cacf1e54635685803ec4edbe0d4fe8465' > expected &&
    shasum cat.gif | cut -d ' ' -f 1 > actual &&
    diff expected actual
"

test_expect_success "retrieve a known popular file with -o" "
    ipget -o meow.gif --node=spawn QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif &&
    echo 'c5ea0d6cacf1e54635685803ec4edbe0d4fe8465' > expected &&
    shasum meow.gif | cut -d ' ' -f 1 > actual &&
    diff expected actual
"

test_expect_success "retrieve a directory" "
    ipget --node=spawn QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF &&
    ls QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF > /dev/null &&
    ls QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif > /dev/null
"

test_done
