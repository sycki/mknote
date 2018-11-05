#!/bin/bash -e

usage(){
  cat <<EOF
Usage:
    $0 [bin | tar | install [path]]
Example:
    $0 install /usr/local/mknote
EOF
}

# build binary
bin(){
  [ x$version == x ] && {
    usage
    exit 1
  }

  mkdir -p "_output"
  GOOS=$os go build -ldflags "-X main.version=$version" -o _output/mknote ./cmd/mknote || {
    code=$?
    echo "failed to build, exit code: $code"
    exit $code
  }
}

# make tar ball
tarball(){
  rm -rf "_output/$name"
  cp -r build _output/$name
  cd _output
  mv mknote $name/bin/
  rm -rf "$name.tar.gz"
  tar -zcf "$name.tar.gz" $name
}

install(){
  local dir=$1
  [[ x$dir == x ]] && dir=/usr/local/mknote
  echo "Prepare install mknote to $dir"

  [ -e $dir ] && {
    echo "The dir is exists already: $dir"
    exit 1
  }

  cp -r _output/$name $dir && {
    echo "Successful install mknote to $dir"
    echo 'Start mknote command:'
    echo -e "\t$dir/start.sh"
  }

  return 0
}

cd `dirname $0`
version=$(<version)
os=linux
name=mknote-$version-$os
cmd=$1

case $cmd in
  -h|--help)
    usage
    ;;
  bin)
    bin
    ;;
  tar)
    bin && tarball
    ;;
  install)
    bin && tarball && install $2
    ;;
  *)
    bin
    ;;
esac
