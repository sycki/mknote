#!/bin/bash

set -e
target_dir="_output"
version=$1
cd `dirname $0`

usage(){
    cat <<-EOF
	Usage:
	    $0 <version>
EOF
}

build(){
    [ x$version == x ] && {
        usage
        exit 1
    }

    GOOS=linux go build -ldflags "-X main.version=$version" -o $target_dir/mknote ./cmd/mknote || {
        code=$?
        echo "failed to build, exit code: $code"
        exit $code
    }

    if [ -d "$target_dir" ]; then
      rm -rf "$target_dir/mknote-$version.tar"
    else
      mkdir $target_dir
    fi

    cp -r static $target_dir/
    cd $target_dir
    tar -cf "mknote-$version.tar" mknote static
    rm -rf static
}

build
