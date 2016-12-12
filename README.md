
## rsyslog repeater. 

Listens to incoming syslog messages and replays these messages to the systemd journald.


#### Usage 
```sh
./rsyslog-journald-repeater -h
Usage of ./rsyslog-journald-repeater:
  -debug
        debug flag
  -host string
        hostname to listen on (default "0.0.0.0")
  -port int
        port to listen on (default 5514 on darwin, default 514 on Linux, Unix, etc.)
```


#### Testing
```bash
$ logger -d -p err  -n localhost  -P 514 -i foo -t centos  "MESSAGE3 FROM LOGGER OVER PORT 514" 
$ journalctl -n 1
-- Logs begin at Mon 2016-12-12 21:39:05 UTC, end at Mon 2016-12-12 21:51:12 UTC. --
Dec 12 21:51:12 fb765a1a956a foo[31]: MESSAGE3 FROM LOGGER OVER PORT 514
```
