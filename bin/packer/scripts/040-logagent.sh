#!/bin/bash
set -e

echo Install log agent -------------------------------
mkdir /tmp/logagent
cd /tmp/logagent
sudo mv /home/centos/etc/systemd/system/rsyslog-journald-repeater.service /etc/systemd/system/rsyslog-journald-repeater.service
sudo chmod 664 /etc/systemd/system/rsyslog-journald-repeater.service

sudo systemctl enable rsyslog-journald-repeater.service
sudo rm -rf /tmp/logagent
echo DONE installing log agent -------------------------------
