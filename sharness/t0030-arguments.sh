#!/usr/bin/env bash

test_description="test the ipget argument parser"

. lib/test-lib.sh

# start the local ipfs node
test_init_ipfs
test_launch_ipfs_daemon

test_expect_success "create a test file" '
	echo "hello ipget" > test_file &&
	ipfs add -q test_file > hash
	cat hash
'
test_expect_success "retrieve a file with a known-gateway URL" '
    ipget -o from_gateway --node=local "https://ipfs.io/ipfs/$(<hash)" &&
    test_cmp test_file from_gateway
'

test_expect_success "retrieve a file with protocol URI" '
    ipget -o from_gateway --node=local "ipfs://$(<hash)" &&
    test_cmp test_file from_gateway
'

test_expect_success "don't allow non-(HTTP)gateway URLS" '
    test_expect_code 1 ipget --node=local "ftp://ipfs.io/ipfs/$(<hash)"
'

# kill the local ipfs node
test_kill_ipfs_daemon

test_done
