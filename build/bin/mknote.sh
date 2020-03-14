#!/bin/bash

cmd::help() {
    cat <<EOF
    Usage:
        $0 <Command>
    Command:
        build
        start
        stop
        backup
        help
EOF

    exit
}

cmd::start() {
    if [ x`ls conf/*.pem 2>/dev/null | wc -l` = x2 ]; then
        exec bin/mknote --tls &>> log/mknote.log &
    else
        exec bin/mknote &>> log/mknote.log &
    fi
    sleep 2
    cmd::status
}

cmd::stop() {
    cmd::status &> /dev/null || return $?
    local pid=`ps -ef | grep mknote | grep -v grep | grep -v $sname | head -n 1 | awk '{print $2}'`
    kill $pid
    sleep 1
    cmd::status
}

cmd::status() {
    local count=`ps -ef | grep mknote | grep -v grep | grep -v $sname | wc -l`
    if (($count >= 1)); then
        echo "RUNNING"
        return 0
    else
        echo "STOPPED"
        return 1
    fi
}

cmd::backup() {
    local dir_name=`basename $(pwd)`
    local dt=`date +%Y%m%d`
    cd ..
    /bin/rm -f "mknote-*.bak.tar.gz"
    tar -zcf mknote-${dt}.bak.tar.gz $dir_name --exclude f/tmp
    cd $dir_name
}


sname=$(basename $0)
cd $(dirname $0)/..

case $1 in
    start)
        cmd::start;;
    stop)
        cmd::stop;;
    status)
        cmd::status;;
    restart)
        cmd::stop ; cmd::start;;
    backup)
        cmd::backup;;
    *)
        cmd::help;;
esac
