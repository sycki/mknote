#!/bin/bash

cd `dirname $0`

target_dir=$1

[ ! -z $target_dir ] || [ "${target_dir:0:1}" != "/" ] {
  echo "请指定mknote安装目录，以/开头，如：/usr/local/mknote"
  exit 11
}

go build mknote || {
  echo "编译失败！exit code $?"
  exit 13
}

[ -d $target_dir ] || mkdir $target_dir

/bin/cp -rf mknote conf articles static  $target_dir/

