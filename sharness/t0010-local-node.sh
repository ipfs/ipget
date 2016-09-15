#!/bin/sh

test_description="test the ipget command using a local node"


. lib/test-lib.sh

# start the local ipfs node
test_init_ipfs
test_launch_ipfs_daemon

test_expect_success "create a test file" "
    echo 'hello ipget' | ipfs add -q > hash
    cat hash
"
file=$(cat hash)

test_expect_success "retrieve a single file" "
    ipget --node=local $file &&
    echo '541472026d4a69afc7e668926d3d6893a3024f8e' > expected &&
    shasum $file | cut -d ' ' -f 1 > actual &&
    diff expected actual
"

test_expect_success "retrieve a single file with -o" "
    ipget -o data.txt --node=local $file &&
    shasum data.txt | cut -d ' ' -f 1 > actual &&
    diff expected actual
"

test_expect_success "create a test directory" "
    mkdir test_dir &&
    cp $file test_dir/data.txt &&
    ipfs add -rq test_dir | tail -n 1 > dir_hash
"

dir=$(cat dir_hash)
test_expect_success "retrieve a directory" "
    ipget --node=local -o got_dir $dir &&
    ls got_dir > /dev/null &&
    ls got_dir/data.txt > /dev/null
"

# kill the local ipfs node
test_kill_ipfs_daemon

test_done
