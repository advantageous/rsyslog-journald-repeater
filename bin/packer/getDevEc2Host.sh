#!/usr/bin/env bash

aws ec2 describe-instances --filters  "Name=tag:Name,Values=i.int.dev.rysyslog-journald-repeater" | jq --raw-output .Reservations[].Instances[].PublicDnsName
