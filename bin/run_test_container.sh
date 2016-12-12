#!/usr/bin/env bash
docker pull advantageous/golang-rsyslog-journald-repeater:latest
docker run  -it --name runner2  \
-p 5514:514/udp \
-v `pwd`:/gopath/src/github.com/advantageous/rsyslog-journald-repeater \
advantageous/golang-rsyslog-journald-repeater
docker rm runner2
