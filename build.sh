#!/bin/bash

set -e

target_dir="_output"

go build mknote || {
  code=$?
  echo "failed build, exit code: $code"
  exit $code
}

version=`./mknote --version 2>&1` || {
  code=$?
  echo "failed get version, exit code: $code"
  exit $code
}

if [ -d $target_dir ]; then
  rm -rf "$target_dir/$version.tar"
else
  mkdir $target_dir
fi

tar -cf "$target_dir/$version.tar" mknote static/
