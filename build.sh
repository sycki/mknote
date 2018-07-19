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
  GOOS=linux go build -ldflags "-X main.version=$version" -o _output/mknote ./cmd/mknote || {
    code=$?
    echo "failed to build, exit code: $code"
    exit $code
  }
}

# make tar ball
tarball(){
  rm -rf "_output/mknote-$version"
  cp -r build _output/mknote-$version
  cd _output
  cp mknote mknote-$version/bin/
  rm -rf "mknote-$version.tar"
  tar -cf "mknote-$version.tar" mknote
}

install(){
  local dir=$1
  [[ x$dir == x ]] && dir=/usr/local/mknote
  echo 'Prepare install mknote to $dir'

  [ -e $dir ] && {
    echo 'The dir is exists already! at $dir'
    exit 1
  }

  cp -r _output/mknote-$version $dir && {
    echo 'Successful install mknote to $dir'
  }

  return 0
}

cd `dirname $0`
version=$(<version)
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

