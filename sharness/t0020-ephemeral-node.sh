#!/usr/bin/env bash

test_description="test the ipget command by spawning a shell"

. lib/test-lib.sh

# start the local ipfs node
test_init_ipfs
test_launch_ipfs_daemon

PEER_ID="$(ipfs id -f'<id>')"

test_expect_success "create a test directory" '
 	mkdir test_dir &&
	echo "hello ipget" > test_dir/test_file &&
	ipfs add -Qr test_dir > hash
'

test_expect_success "retrieve a single file" '
	echo "$SWARM_MADDR"
    ipget --peers="$SWARM_MADDR/p2p/$PEER_ID" --node=spawn "/ipfs/$(<hash)/test_file" &&
    test_cmp test_dir/test_file test_file
'

test_expect_success "retrieve a single file with -o" '
    ipget --peers="$SWARM_MADDR/p2p/$PEER_ID" --node=spawn -o out_file "/ipfs/$(<hash)/test_file" &&
    test_cmp test_dir/test_file out_file
'

test_expect_success "retrieve a directory" '
    ipget --peers="$SWARM_MADDR/p2p/$PEER_ID" --node=spawn "/ipfs/$(<hash)" &&
    diff -ru "$(<hash)" test_dir
'

# kill the local ipfs node
test_kill_ipfs_daemon

test_done
