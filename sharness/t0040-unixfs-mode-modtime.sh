#!/usr/bin/env bash

test_description="test the ipget command setting unix mode and modification time"


. lib/test-lib.sh

# start the local ipfs node
test_init_ipfs
test_launch_ipfs_daemon

test_expect_success "create test file with mode and mtime" '
	echo "hello ipget" > test_file &&
	ipfs add -q --mode=0666 --mtime=660000000 test_file > hash
	cat hash
'
test_expect_success "retrieve file with mode and mtime" '
    ipget -o data.txt --node=local "/ipfs/$(<hash)" &&
    test_cmp test_file "data.txt" &&
    case $(uname -s) in
    Linux)
        stat --format="%Y %a" data.txt > out &&
        echo "660000000 666" > expect && 
        test_cmp expect out
        ;;
    Darwin|FreeBSD)
        stat -f "%m %p" data.txt > out &&
        echo "660000000 100666" > expect &&
        test_cmp expect out
        ;;
    *)
        echo "unsupported system: $(uname)"
    esac
'

test_expect_success "create a test directory" '
    mkdir test_dir &&
    cp test_file test_dir/data.txt &&
    ipfs add --mode=0777 --mtime=660000000 -rQ test_dir > dir_hash
'

test_expect_success "retrieve a directory with mode and mtime" '
    ipget --node=local -o got_dir "/ipfs/$(<dir_hash)" &&
    case $(uname -s) in
    Linux)
        stat --format="%Y %a" got_dir > out2 &&
        echo "660000000 777" > expect2 &&
        test_cmp expect2 out2 &&
        stat --format="%Y %a" got_dir/data.txt > out3 &&
        echo "660000000 777" > expect3 &&
        test_cmp expect3 out3
        ;;
    Darwin|FreeBSD)
        stat -f "%m %p" got_dir > out2 &&
        echo "660000000 40777" > expect2 &&
        test_cmp expect2 out2 &&
        stat -f "%m %p" got_dir/data.txt > out3 &&
        echo "660000000 100777" > expect3 &&
        test_cmp expect3 out3
        ;;
    *)
        echo "unsupported system: $(uname)"
    esac
'

test_expect_success "create a test directory with symlink" '
    case $(uname -s) in
    Linux|FreeBSD|Darwin)
        mkdir test_dir2 &&
        cp test_file test_dir2/data.txt &&
        ln -s test_file test_dir2/test_file_link &&
        ipfs add --mtime=660000000 -rQ test_dir2 > dir2_hash
        ;;
    *)
        echo "unsupported system: $(uname)"
    esac
'

test_expect_success "retrieve a directory with symlink with mode and mtime" '
    case $(uname -s) in
    Linux)
        ipget --node=local -o got_dir2 "/ipfs/$(<dir2_hash)" &&
        readlink got_dir2/test_file_link > link_target &&
        echo "test_file" > expect_target &&
        test_cmp expect_target link_target &&
        stat --format="%Y" got_dir2/test_file_link > out4 &&
        echo "660000000" > expect4 &&
        test_cmp expect4 out4
        ;;
    Darwin|FreeBSD)
        ipget --node=local -o got_dir2 "/ipfs/$(<dir2_hash)" &&
        readlink got_dir2/test_file_link > link_target &&
        echo "test_file" > expect_target &&
        test_cmp expect_target link_target
        ;;
    *)
        echo "unsupported system: $(uname)"
    esac
'

# kill the local ipfs node
test_kill_ipfs_daemon

test_done
