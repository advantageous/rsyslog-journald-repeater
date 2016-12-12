#!/usr/bin/env bash
source ~/.bash_profile
docker pull advantageous/golang-rsyslog-journald-repeater:latest
docker run -it --name build \
-v `pwd`:/gopath/src/github.com/advantageous/rsyslog-journald-repeater \
advantageous/golang-rsyslog-journald-repeater \
/bin/sh -c "/gopath/src/github.com/advantageous/rsyslog-journald-repeater/bin/build_linux.sh"
docker rm build

