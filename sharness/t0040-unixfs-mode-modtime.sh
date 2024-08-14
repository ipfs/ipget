#!/usr/bin/env bash

test_description="test the ipget command setting unix mode and modification time"


. lib/test-lib.sh

# start the local ipfs node
test_init_ipfs
test_launch_ipfs_daemon

test_expect_success "create test file with mode and mtime" '
	echo "hello ipget" > test_file &&
	ipfs add -q --mode=0664 --mtime=660000000 test_file > hash
	cat hash
'
test_expect_success "retrieve file with mode and mtime" '
    ipget -o data.txt --node=local "/ipfs/$(<hash)" &&
    test_cmp test_file "data.txt" &&
    if test $(uname -s) = "Darwin"; then
        stat -f "%m %p" data.txt > out &&
        echo "660000000 100664" > expect &&
        test_cmp expect out
    elif test $(uname -s) = "Linux"; then
        stat --format="%Z %a" data.txt > out &&
        echo "660000000 664" > expect2 && 
        test_cmp expect out
    fi
'

test_expect_success "create a test directory" '
    mkdir test_dir2 &&
    cp test_file test_dir2/data.txt &&
    ipfs add --mode=0775 --mtime=660000000 -rQ test_dir2 > dir_hash
'

test_expect_success "retrieve a directory" '
    ipget --node=local -o got_dir "/ipfs/$(<dir_hash)" &&
    if test $(uname -s) = "Darwin"; then
        stat -f "%m %p" got_dir > out2 &&
        echo "660000000 40775" > expect2 &&
        test_cmp expect2 out2
    elif test $(uname -s) = "Linux"; then
        stat --format="%Z %a" got_dir > out2 &&
        echo "660000000 40775" > expect2 &&
        test_cmp expect2 out2
    fi

    
'

# kill the local ipfs node
test_kill_ipfs_daemon

test_done
