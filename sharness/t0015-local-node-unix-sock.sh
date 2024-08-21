#!/usr/bin/env bash

test_description="test the ipget command using a local node"


. lib/test-lib.sh

# start the local ipfs node
test_init_ipfs

# We can't use the trash dir as the full name must be longer less than 108 bytes
# long (because that's the max unix domain socket path length).
SOCKDIR="$(mktemp -d "${TMPDIR:-/tmp}ipget-sharness.XXXXXX")"
SOCKPATH="${SOCKDIR}/sock"
SOCKADDR="/unix${SOCKPATH}"

test_expect_success "configure with unix socket" '
    ipfs config Addresses.API "$SOCKADDR"
'

test_launch_ipfs_daemon_unixsocket "$SOCKPATH"

test_expect_success "create a test file" '
	echo "hello ipget" > test_file &&
	ipfs --api="$SOCKADDR" add -q test_file > hash &&
	cat hash
'
test_expect_success "retrieve a single file" '
    ipget --node=local "$(<hash)" &&
    test_cmp test_file "$(<hash)"
'

test_expect_success "retrieve a single file with -o" '
    ipget -o data.txt --node=local "$(<hash)" &&
    test_cmp test_file "data.txt"
'

test_expect_success "create a test directory" '
    mkdir test_dir &&
    cp test_file test_dir/data.txt &&
    ipfs add -rQ test_dir > dir_hash
'

test_expect_success "retrieve a directory" '
    ipget --node=local -o got_dir "$(<dir_hash)" &&
	diff -ru test_dir got_dir
'

# kill the local ipfs node
test_kill_ipfs_daemon

test_done
