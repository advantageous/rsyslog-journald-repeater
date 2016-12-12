#!/usr/bin/env bash


rm rsyslog-journald-repeater_linux

set -e

/usr/lib/systemd/systemd-journald &

cd /gopath/src/github.com/advantageous/rsyslog-journald-repeater/

echo "Running go clean"
go clean
echo "Running go get"
go get
echo "Running go build"
go build
echo "Renaming output to _linux"
mv rsyslog-journald-repeater rsyslog-journald-repeater_linux
cp  rsyslog-journald-repeater_linux /usr/lib/rsyslog-journald-repeater

pkill -9 systemd